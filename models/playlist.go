package models

type Playlist struct {
	ID       PlaylistID `json:"id"`
	SourceID SourceID   `json:"sourceID"`
	Name     string     `json:"name"`
}
