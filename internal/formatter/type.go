package formatter

type Media struct {
	MediaType  MediaType
	URL        string
	PreviewURL string
}

type MediaType string

const MediaTypeVideo MediaType = "video"
const MediaTypeGIF MediaType = "GIF"
const MediaTypePhoto MediaType = "photo"
