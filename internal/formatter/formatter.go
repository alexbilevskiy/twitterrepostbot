package formatter

import (
	"fmt"
	"strings"
	"time"

	twitterscraper "github.com/imperatrona/twitter-scraper"
)

func FormatTweet(t *twitterscraper.Tweet) (string, bool) {
	date := t.TimeParsed.Local().Format(time.DateTime)
	disablePreview := true
	text := t.Text
	text += "\n" + fmt.Sprintf(`— %s (@%s) <a href="%s">%s</a>`, t.Name, t.Username, t.PermanentURL, date)
	medias := buildMedia(t)
	if len(medias) == 1 {
		text = fmt.Sprintf(`<a href="%s">&#8204;</a>`, medias[0]) + text
		disablePreview = false
	} else if len(medias) > 1 {
		disablePreview = false
		var mediasText []string
		for i, url := range medias {
			mediasText = append(mediasText, fmt.Sprintf(`<a href="%s">m%d</a>`, url, i))
		}
		text = fmt.Sprintf(`<a href="%s">&#8204;</a>`, medias[0]) + strings.Join(mediasText, " ") + "\n" + text
	}

	return text, disablePreview
}

func buildMedia(t *twitterscraper.Tweet) []string {
	medias := make([]string, 0)
	if len(t.Photos) > 0 {
		for _, photo := range t.Photos {
			medias = append(medias, photo.URL)
		}
	}
	if len(t.Videos) > 0 {
		for _, video := range t.Videos {
			medias = append(medias, video.URL)
		}
	}
	if len(t.GIFs) > 0 {
		for _, gif := range t.GIFs {
			medias = append(medias, gif.URL)
		}
	}

	return medias
}
