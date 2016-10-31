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
	Version = "0.0.2"
	// SlackToken will be set by ENV or config file in init()
	SlackToken = ""
)

func init() {
	viper.AddConfigPath("$HOME")
	viper.AddConfigPath(".")
	viper.SetConfigName(".mitch")
	viper.AutomaticEnv()
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)

	}

	SlackToken = viper.GetString("MITCH_SLACK_TOKEN")
}

func main() {
	bot, err := hanu.New(SlackToken)

	if err != nil {
		log.Fatal(err)
	}

	skills.Version = Version
	skills.Start = time.Now()
	cmdList := skills.List()
	for i := 0; i < len(cmdList); i++ {
		bot.Register(cmdList[i])
	}

	fmt.Println("Starting up!")
	bot.Listen()
}
