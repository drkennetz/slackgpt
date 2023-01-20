package main

import (
	"context"
	"github.com/PullRequestInc/go-gpt3"
	"github.com/spf13/viper"
	gptslack "go-slack-chat-gpt3/src/slack"
	"log"
	"os"
)

func main() {
	log.SetOutput(os.Stdout)
	viper.SetConfigFile(".env")
	viper.ReadInConfig()
	cgptApiKey := viper.GetString("CGPT_API_KEY")
	if cgptApiKey == "" {
		log.Fatalln("Missing chat-gpt API KEY")
	}
	slackAppToken := viper.GetString("SLACK_APP_TOKEN")
	if slackAppToken == "" {
		log.Fatalln("Missing slack app token")
	}
	slackBotToken := viper.GetString("SLACK_BOT_TOKEN")
	if slackBotToken == "" {
		log.Fatalln("Missing slack bot token")
	}
	ctx := context.Background()
	client := gpt3.NewClient(cgptApiKey)
	gptslack.EventHandler(slackAppToken, slackBotToken, client, ctx)
}
