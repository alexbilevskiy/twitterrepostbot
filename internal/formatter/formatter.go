package formatter

import (
	"fmt"
	"strings"
	"time"

	twitterscraper "github.com/imperatrona/twitter-scraper"
)

func FormatTweet(t *twitterscraper.Tweet) (string, []Media) {
	date := t.TimeParsed.Local().Format(time.DateTime)
	text := t.Text
	text += "\n" + fmt.Sprintf(`— %s (@%s) <a href="%s">%s</a>`, t.Name, t.Username, t.PermanentURL, date)
	medias := buildMedia(t)
	for _, m := range medias {
		// either twitter-scraper or twitter itself adds plaintext links to media in text
		if m.PreviewURL != "" {
			text = strings.Replace(text, m.PreviewURL, "", 1)
		} else {
			text = strings.Replace(text, m.URL, "", 1)
		}
	}

	return text, medias
}

func buildMedia(t *twitterscraper.Tweet) []Media {
	medias := make([]Media, 0)
	if len(t.Photos) > 0 {
		for _, photo := range t.Photos {
			medias = append(medias, Media{MediaType: MediaTypePhoto, URL: photo.URL})
		}
	}
	if len(t.Videos) > 0 {
		for _, video := range t.Videos {
			medias = append(medias, Media{MediaType: MediaTypeVideo, URL: video.URL, PreviewURL: video.Preview})
		}
	}
	if len(t.GIFs) > 0 {
		for _, gif := range t.GIFs {
			medias = append(medias, Media{MediaType: MediaTypeGIF, URL: gif.URL, PreviewURL: gif.Preview})
		}
	}

	return medias
}
