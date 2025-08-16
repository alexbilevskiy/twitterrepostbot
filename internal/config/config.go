package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type AppCfg struct {
	LastTweetFile string `envconfig:"last_tweet_file" default:".cache/last.json"`
	CookiesCache  string `envconfig:"cookies_cache" default:".cache/cookies.json"`
	Login         string `envconfig:"x_login" default:""`
	Password      string `envconfig:"x_password" default:""`
	BotToken      string `envconfig:"tg_token" default:""`
	ChatId        int64  `envconfig:"tg_chat_id" default:""`
}

func ReadConfig() *AppCfg {
	cfg := &AppCfg{}
	_ = godotenv.Load(".env")
	_ = envconfig.Process("", cfg)
	return cfg
}
