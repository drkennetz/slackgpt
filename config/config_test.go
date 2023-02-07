package configs

import (
	"errors"
	"github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParseConfigFromPath(t *testing.T) {
	type args struct {
		cfg     string
		cfgType string
	}
	type expectedResult struct {
		cfgParts configParts
		e        error
	}
	tests := []struct {
		name string
		args args
		want expectedResult
	}{
		{
			"empty config type error",
			args{
				"config",
				"",
			},
			expectedResult{
				configParts{},
				errors.New("cfg type not passed, and no file extension. cannot infer config type"),
			},
		},
		{
			"yml conversion, valid config",
			args{
				"config.yml",
				"",
			},
			expectedResult{
				configParts{"not_empty", "config.yml", "yaml"},
				nil,
			},
		},
		{
			"invalid config type error",
			args{
				"config.xml",
				"xml",
			},
			expectedResult{
				configParts{},
				errors.New("invalid cfg type"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfgParts, err := ParseConfigFromPath(tt.args.cfg, tt.args.cfgType)
			assert.Equal(t, cfgParts.Type, tt.want.cfgParts.Type)
			assert.Equal(t, cfgParts.Name, tt.want.cfgParts.Name)
			if err != nil {
				require.Error(t, err)
			}
		})
	}
}

func TestLoadConfig(t *testing.T) {
	type args struct {
		cfgParts configParts
	}
	type expectedResult struct {
		cfg Config
		e   error
	}
	tests := []struct {
		name string
		args args
		want expectedResult
	}{
		{
			"bad file path",
			args{
				configParts{
					"./test_files",
					"nofile",
					"env",
				},
			},
			expectedResult{
				Config{},
				errors.New("Config File \"nofile\" Not Found "),
			},
		},
		{
			"no gpt key",
			args{
				configParts{
					"./test_files",
					"no_cgpt.json",
					"env",
				},
			},
			expectedResult{
				Config{},
				errors.New("missing chat-gpt API key"),
			},
		},
		{
			"no app token",
			args{
				configParts{
					"./test_files",
					"no_slack_app.json",
					"env",
				},
			},
			expectedResult{
				Config{
					ChatGPTKey: "test",
				},
				errors.New("missing slack app token"),
			},
		},
		{
			"no bot token",
			args{
				configParts{
					"./test_files",
					"no_slack_bot.json",
					"env",
				},
			},
			expectedResult{
				Config{
					ChatGPTKey:    "test",
					SlackAppToken: "test",
				},
				errors.New("missing slack bot token"),
			},
		},
		{
			"bad app token",
			args{
				configParts{
					"./test_files",
					"bad_slack_app.json",
					"env",
				},
			},
			expectedResult{
				Config{
					ChatGPTKey:    "test",
					SlackAppToken: "test",
					SlackBotToken: "test",
				},
				errors.New("slack app token should begin with xapp-"),
			},
		},
		{
			"bad bot token",
			args{
				configParts{
					"./test_files",
					"bad_slack_bot.json",
					"env",
				},
			},
			expectedResult{
				Config{
					ChatGPTKey:    "test",
					SlackAppToken: "xapp-1",
					SlackBotToken: "test",
				},
				errors.New("slack bot token should begin with xoxb-"),
			},
		},
		{
			"good",
			args{
				configParts{
					"./test_files",
					"good.json",
					"env",
				},
			},
			expectedResult{
				Config{
					ChatGPTKey:    "test",
					SlackAppToken: "xapp-1",
					SlackBotToken: "xoxb-1",
				},
				errors.New(""),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := LoadConfig(tt.args.cfgParts)
			assert.Equal(t, cfg.SlackBotToken, tt.want.cfg.SlackBotToken)
			assert.Equal(t, cfg.SlackAppToken, tt.want.cfg.SlackAppToken)
			assert.Equal(t, cfg.ChatGPTKey, tt.want.cfg.ChatGPTKey)
			if err != nil {
				require.ErrorContains(t, err, tt.want.e.Error())
			}
		})
	}

}
