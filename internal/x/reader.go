package x

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/alexbilevskiy/twitterbot/internal/config"
	twitterscraper "github.com/imperatrona/twitter-scraper"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

const dateLimit = 4 * 24 * time.Hour

var errInvalidAuth = errors.New("invalid auth")

type XReader struct {
	cfg *config.AppCfg
	t   *twitterscraper.Scraper
}

func NewXReader(cfg *config.AppCfg) *XReader {
	return &XReader{cfg: cfg}
}

func (r *XReader) Login() error {
	scraper := twitterscraper.New()

	_, err := os.Stat(r.cfg.CookiesCache)
	if os.IsNotExist(err) {
		log.Printf("starting auth")
		err := scraper.Login(r.cfg.Login, r.cfg.Password)
		if err != nil {
			return fmt.Errorf("%w: %w", errInvalidAuth, err)
		}
		cookies := scraper.GetCookies()
		bytes, _ := json.Marshal(cookies)
		_ = os.WriteFile(r.cfg.CookiesCache, bytes, 0644)
	} else {
		log.Printf("using cached cookies")
		bytes, _ := os.ReadFile(r.cfg.CookiesCache)
		var cookies []*http.Cookie
		_ = json.Unmarshal(bytes, &cookies)
		scraper.SetCookies(cookies)
	}

	if !scraper.IsLoggedIn() {
		return fmt.Errorf("%w: not logged in", errInvalidAuth)
	}
	r.t = scraper
	return nil
}

func (r *XReader) ReadHome() error {
	if r.t == nil {
		return errors.New("auth first")
	}
	var cursor string
	page := 0
	var first *twitterscraper.Tweet
	var last *twitterscraper.Tweet
	_, err := os.Stat(r.cfg.LastTweetFile)
	if err == nil {
		log.Printf("using last tweet file")
		bytes, _ := os.ReadFile(r.cfg.LastTweetFile)
		_ = json.Unmarshal(bytes, &last)
	} else {
		log.Printf("not using last tweet file")
	}
	for {
		log.Printf("Page: %d, cursor: %s", page, cursor)

		tweets, newCursor, err := r.t.FetchHomeTweets(20, cursor)
		if err != nil {
			log.Printf("fetch home tweets: %v", err)
			continue
		}
		cookies := r.t.GetCookies()
		bytes, _ := json.Marshal(cookies)
		_ = os.WriteFile(r.cfg.CookiesCache, bytes, 0644)

		cursor = newCursor
		page++
		finished := false
		log.Printf("Fetched %d tweets", len(tweets))
		for _, tweet := range tweets {
			if first == nil {
				first = tweet
			}
			url := fmt.Sprintf("https://x.com/%s/status/%s", tweet.Username, tweet.ID)
			date := tweet.TimeParsed.Local().Format(time.DateTime)

			if last != nil {
				lastId, _ := strconv.ParseInt(last.ID, 10, 64)
				curId, _ := strconv.ParseInt(tweet.ID, 10, 64)
				if lastId >= curId {
					log.Printf("[%s] Received already processed tweet: %s", date, url)
					finished = true
					break
				}
			}
			if time.Since(tweet.TimeParsed) > dateLimit {
				log.Printf("[%s] Received too old tweet: %s", date, url)
				finished = true
				break
			}

			if tweet.IsRetweet {
				//log.Printf("[%s] Skip retweet: %s", date, url)
				continue
			}
			log.Printf("[%s] Received new tweet: %s", date, url)
		}
		if finished {
			break
		}
	}
	if first != nil {
		log.Printf("storing last tweet: %s", first.TimeParsed.Local().Format(time.DateTime))
		bytes, _ := json.Marshal(first)
		_ = os.WriteFile(r.cfg.LastTweetFile, bytes, 0644)
	}

	return nil
}
