package main

import (
	"fmt"
	"log"
	"time"

	"github.com/sbstjn/hanu"
	"github.com/wlbr/mitch/skills"
	"github.com/spf13/viper"
)

var (
	// Version is the bot version
	Version string = "Unknown build."
	Githash string = "Unknown git commit hash."
	Buildstamp string = "Unknown build timestamp."

	// SlackToken will be set by ENV or config file in init()
	SlackToken = ""
)


func init() {

	fmt.Printf("Version: %s (%s) of %s \n", Version, Githash, Buildstamp)

	fmt.Println("Configuring...")
	viper.AddConfigPath("$HOME")
	viper.AddConfigPath(".")
	viper.SetConfigName(".mitch")
	viper.AutomaticEnv()
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		log.Println(err)
	}

	SlackToken = viper.GetString("MITCH_SLACK_TOKEN")
	if SlackToken == "" {
		log.Fatal("No Slack auth token found.")
	}
}

func main() {

	bot, err := hanu.New(SlackToken)

	if err != nil {
		log.Fatal(err)
	}

	skills.Version = fmt.Sprintf("%s (%s) of %s", Version, Githash, Buildstamp)
	skills.Start = time.Now()
	cmdList := skills.List()
	for i := 0; i < len(cmdList); i++ {
		bot.Register(cmdList[i])
	}

	fmt.Println("Listening...")

	bot.Listen()
}
