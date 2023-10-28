package models

type MediaType string

const (
	MediaTypeText  MediaType = "text"
	MediaTypeImage MediaType = "image"
	MediaTypeVideo MediaType = "video"
	MediaTypeAudio MediaType = "audio"
)
