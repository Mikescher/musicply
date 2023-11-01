package models

import (
	"strings"
)

type Playlist struct {
	ID       PlaylistID `json:"id"`
	Sort     int        `json:"sort"`
	SourceID SourceID   `json:"sourceID"`
	Path     string     `json:"path"`
	Name     string     `json:"name"`

	Cover *CoverHash `json:"cover"`
}

func (p Playlist) NameParts() []string {
	splt := strings.Split(p.Name, "/")
	if len(splt) == 0 {
		return []string{""}
	}
	return splt
}

type HierarchicalPlaylist struct {
	ID       *PlaylistID            `json:"id"`
	HierID   HierarchicalPlaylistID `json:"hierarchicalID"`
	Name     string                 `json:"name"`
	Children []HierarchicalPlaylist `json:"children"`

	TrackCount int        `json:"trackCount"`
	Cover      *CoverHash `json:"cover"`
}
