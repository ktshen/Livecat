package twitchstreams_crawler

// Twitch API /helix/first=100
type TwitchJson struct {
	Data       []Data
	Pagination Pagination
}
type Data struct {
	ID            string
	User_ID       string
	User_Name     string
	Game_ID       string
	Community_IDs []string
	Type          string
	Title         string
	Viewer_Count  int
	Started_At    string
	Language      string
	Thumbnail_URL string
	Tag_IDs       []string
}
type Community_ids struct {
	id []string
}
type Pagination struct {
	Cursor string
}

// Twitch API summary
type Summary struct {
	Channels int
	viewers  int
	_links   []string
}

// Twitch API Get game name
type GameNameJson struct {
	DATA       []Data2
	Pagination Pagination2
}
type Data2 struct {
	Id          string
	Name        string
	Box_art_url string
}
type Pagination2 struct {
	Cursor string
}
type GameNameFileJson struct {
	DATA []Data2
}
