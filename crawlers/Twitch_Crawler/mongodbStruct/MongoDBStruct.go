package mongodbStruct

import "time"

type MongoDB struct {
	Title            string
	Description      string
	Platform         string
	VideoID          string
	Host             string
	Status           string
	Thumbnails       string
	Published        string
	Tags             string
	GeneralTag       string
	Timestamp        time.Time
	Language         string
	ViewCount        int
	Viewers          int
	VideoURL         string
	VideoEmbedded    string
	ChatRoomEmbedded string
	Channel          string
	CreatedAt        time.Time
}

type Web struct {
	Host string
}
