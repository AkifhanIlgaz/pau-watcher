package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	TelegramToken string `mapstructure:"TELEGRAM_TOKEN"`
	ChatId        int64  `mapstructure:"CHAT_ID"`
	Pau           string `mapstructure:"PAU"`
	Chain         string
}

func Load(path string) (*Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName("app")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	return &config, nil
}
