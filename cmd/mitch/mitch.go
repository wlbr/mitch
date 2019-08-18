package main

import (
	"fmt"
	"log"
	"time"

	"github.com/nlopes/slack"
	"github.com/spf13/viper"
	"github.com/wlbr/mitch/bot"
	"github.com/wlbr/mitch/skills"
)

var (
	//Version is a linker injected variable for a git revision info used as version info
	Version = "Unknown build"
	/*Buildstamp is a linker injected variable for a buildtime timestamp used in version info */
	Buildstamp = "unknown build timestamp."

	// SlackToken will be set by ENV or config file in init()
	SlackToken = ""

	// ArchiveFile will be set by ENV or config file in init()
	ArchiveFile = ""

	config bot.Config
	mitch  bot.Bot
)

func init() {
	mitch = bot.Bot{}
	config = bot.Config{}
	mitch.Config = &config
	mitch.Config.Upstart = time.Now()

	btime, err := time.Parse("2006-01-02_15:04:05_MST", Buildstamp)
	if err != nil {
		config.BuildTimeStamp = time.Now()
	} else {
		config.BuildTimeStamp = btime
	}

	config.GitVersion = Version

	fmt.Printf("Version: %s of %s \n", config.GitVersion, config.BuildTimeStamp)

	fmt.Print("Configuring...")
	viper.AddConfigPath("$HOME")
	viper.AddConfigPath(".")
	viper.SetConfigName(".mitch")
	viper.AutomaticEnv()
	viper.SetConfigType("yaml")

	err = viper.ReadInConfig()
	if err != nil {
		log.Println(err)
	}

	config.SlackToken = viper.GetString("MITCH_SLACK_TOKEN")
	if config.SlackToken == "" {
		log.Fatal("No Slack auth token found.")
	}
	config.OpenWeatherMapToken = viper.GetString("openweathermap_api_token")
	config.ArchiveFile = viper.GetString("archive")

	fmt.Println(" done.")
}

func main() {
	api := slack.New(mitch.Config.SlackToken)
	//logger := log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)
	//slack.SetLogger(logger)
	//api.SetDebug(false)

	fmt.Print("Logging in...")
	mitch.Client = api
	mitch.Rtm = api.NewRTM()

	mitch.RegisterMessageHandler(skills.NewArchiver())
	mitch.RegisterSkillHandler(skills.NewHelp())
	mitch.RegisterSkillHandler(skills.NewStockInformer())
	mitch.RegisterSkillHandler(skills.NewHello())
	mitch.RegisterSkillHandler(skills.NewEchoSkill())
	mitch.RegisterSkillHandler(skills.NewUptimeInfo())
	mitch.RegisterSkillHandler(skills.NewVersionInfo())
	mitch.RegisterSkillHandler(skills.NewTimeIn())
	if config.OpenWeatherMapToken != "" {
		mitch.RegisterSkillHandler(skills.NewWeatherIn())
	}

	go mitch.Rtm.ManageConnection()

	mitch.Rtm.SendMessage(mitch.Rtm.NewOutgoingMessage("Hello world", "#blubb"))

	mitch.MainLoop()
}
