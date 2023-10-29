package models

import "strings"

type Playlist struct {
	ID       PlaylistID `json:"id"`
	SourceID SourceID   `json:"sourceID"`
	Name     string     `json:"name"`
}

func (p Playlist) NameParts() []string {
	splt := strings.Split(p.Name, "/")
	if len(splt) == 0 {
		return []string{""}
	}
	return splt
}

type HierarchicalPlaylist struct {
	ID          *PlaylistID            `json:"id"`
	Name        string                 `json:"name"`
	Children    []HierarchicalPlaylist `json:"children"`
	HasChildren bool                   `json:"hasChildren"`

	TrackCount int `json:"trackCount"`
	Cover      *struct {
		Playlist PlaylistID `json:"playlistID"`
		Track    PlaylistID `json:"trackID"`
	} `json:"cover"`
}
