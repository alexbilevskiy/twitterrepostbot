package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"

	"github.com/alexbilevskiy/twitterbot/internal/config"
	"github.com/alexbilevskiy/twitterbot/internal/formatter"
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
				text, disablePreview := formatter.FormatTweet(tweet)
				_, err := tgBot.SendMessage(cfg.ChatId, text, &gotgbot.SendMessageOpts{ParseMode: "HTML", LinkPreviewOptions: &gotgbot.LinkPreviewOptions{IsDisabled: disablePreview}})
				if err != nil {
					log.Printf("failed to send tg message: %v", err)
					bytes, _ := json.Marshal(tweet)
					_ = os.WriteFile(fmt.Sprintf(".cache/%s.json", tweet.ID), bytes, 0644)
					continue
				}
				time.Sleep(5 * time.Second)
			}
		case <-ctx.Done():
			log.Printf("exiting...")
			os.Exit(0)
		}
	}
}
