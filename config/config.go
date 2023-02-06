// Package config handles all the configuration
package configs

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"golang.org/x/exp/slices"
	"path/filepath"
	"strings"
)

// Config stores the configurations required for the app
type Config struct {
	ChatGPTKey    string `mapstructure:"CGPT_API_KEY"`
	SlackAppToken string `mapstructure:"SLACK_APP_TOKEN"`
	SlackBotToken string `mapstructure:"SLACK_BOT_TOKEN"`
}

// configParts provide a convenience object for parsing input config
type configParts struct {
	AbsPath string
	Name    string
	Type    string
}

// ParseConfigFromPath extracts all relevant info from passed config
func ParseConfigFromPath(cfg, cfgType string) (configParts, error) {
	var cfgParts configParts
	var ext string
	var err error
	validTypes := []string{"yaml", "json", "hcl", "properties", "toml", "env", "ini"}
	abs, _ := filepath.Abs(filepath.Dir(cfg))
	if cfgType == "" {
		ext = strings.Replace(filepath.Ext(cfg), ".", "", -1)
	} else {
		ext = cfgType
	}
	if ext == "" {
		err = errors.New("cfg type not passed, and no file extension. cannot infer config type")
		return cfgParts, err
	}
	if ext == "yml" {
		ext = "yaml"
	}
	if !slices.Contains(validTypes, ext) {
		fmt.Println("cfg type must be in: json, toml, yaml, hcl, ini, env, properties")
		fmt.Println("this can be specified on the command line with -t <type>")
		err = errors.New("invalid cfg type")
		return cfgParts, err
	}
	//ext := strings.Replace(filepath.Ext(cfg), ".", "", -1)
	name := filepath.Base(cfg)
	cfgParts.AbsPath = abs
	cfgParts.Name = name
	cfgParts.Type = ext
	return cfgParts, nil
}

// LoadConfig reads configuration from config
func LoadConfig(cfgParts configParts) (config Config, err error) {
	viper.AddConfigPath(cfgParts.AbsPath)
	viper.SetConfigName(cfgParts.Name)
	viper.SetConfigType(cfgParts.Type)
	if err = viper.ReadInConfig(); err != nil {
		return
	}
	if err = viper.Unmarshal(&config); err != nil {
		return
	}
	if config.ChatGPTKey == "" {
		err = errors.New("missing chat-gpt API key")
		return
	}
	if config.SlackAppToken == "" {
		err = errors.New("missing slack app token")
		return
	}
	if config.SlackBotToken == "" {
		err = errors.New("missing slack bot token")
		return
	}
	if !strings.HasPrefix(config.SlackAppToken, "xapp-") {
		err = errors.New("slack app token should begin with xapp-")
		return
	}
	if !strings.HasPrefix(config.SlackBotToken, "xoxb-") {
		err = errors.New("slack bot token should begin with xoxb-")
		return
	}
	return
}
