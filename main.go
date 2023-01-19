package main

import (
	"context"
	"fmt"
	"github.com/PullRequestInc/go-gpt3"
	"github.com/spf13/viper"
	"go-slack-chat-gpt3/src/chatgpt"
	"log"
	"os"
	"strings"
)

func main() {
	log.SetOutput(os.Stdout)
	statement := strings.Join(os.Args[1:], "")
	viper.SetConfigFile(".env")
	viper.ReadInConfig()
	apiKey := viper.GetString("API_KEY")
	if apiKey == "" {
		panic("Missing API KEY")
	}
	ctx := context.Background()
	client := gpt3.NewClient(apiKey)
	resp, err := chatgpt.GetStringResponse(client, ctx, statement)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(resp)
}
