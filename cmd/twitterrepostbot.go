package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"

	"github.com/alexbilevskiy/twitterbot/internal/config"
	"github.com/alexbilevskiy/twitterbot/internal/x"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()
	go func() {
		<-ctx.Done()
		<-time.After(5 * time.Second)
		log.Printf("service has not been stopped within the specified timeout; killed by force")
		os.Exit(1)
	}()

	cfg := config.ReadConfig()

	tgBot, err := gotgbot.NewBot(cfg.BotToken, nil)
	if err != nil {
		log.Printf("Error creating tg bot: %v", err)
		os.Exit(1)
	}

	xReader := x.NewXReader(cfg)
	err = xReader.Login()
	if err != nil {
		log.Printf("failed to login: %v", err)
		os.Exit(1)
	}

	timer := time.NewTicker(300 * time.Second)
	defer timer.Stop()
	for {
		select {
		case <-timer.C:
			log.Printf("starting by timer")
			tweets, err := xReader.ReadHome()
			if err != nil {
				log.Printf("failed to read home: %v", err)
				os.Exit(1)
			}
			for _, tweet := range tweets {
				_, err := tgBot.SendMessage(cfg.ChatId, fmt.Sprintf("%s", tweet.URL), nil)
				if err != nil {
					log.Printf("failed to send tg message: %v", err)
					continue
				}
			}
		case <-ctx.Done():
			log.Printf("exiting...")
			os.Exit(0)
		}
	}
}
