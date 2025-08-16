package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

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

	xReader := x.NewXReader(cfg)
	err := xReader.Login()
	if err != nil {
		log.Printf("failed to login: %v", err)
		os.Exit(1)
	}

	timer := time.NewTicker(60 * time.Second)
	defer timer.Stop()
	for {
		select {
		case <-timer.C:
			log.Printf("starting by timer")
			err := xReader.ReadHome()
			if err != nil {
				log.Printf("failed to read home: %v", err)
				os.Exit(1)
			}
		case <-ctx.Done():
			log.Printf("exiting...")
			os.Exit(0)
		}
	}
}
