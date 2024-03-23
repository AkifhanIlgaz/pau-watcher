package config

import (
	"errors"
	"fmt"

	"github.com/AkifhanIlgaz/pau-watcher/chain"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Config struct {
	TelegramToken string `mapstructure:"TELEGRAM_TOKEN"`
	ChatId        int64  `mapstructure:"CHAT_ID"`
	WatchAddress  string `mapstructure:"WATCH_ADDRESS"`
	Chain         chain.Chain
}

func Load(path string) (*Config, error) {
	var config Config

	ch := pflag.String("chain", "", "Chain")
	pflag.Parse()

	if len(*ch) == 0 {
		return nil, errors.New("chain is not provided")
	}

	if selectedChain, ok := chain.Chains[*ch]; !ok {
		return nil, fmt.Errorf("chain %s is not supported", *ch)
	} else {
		config.Chain = selectedChain
	}

	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName("app")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	return &config, nil
}
