package models

import "strings"

type CoverData struct {
	Filepath string `json:"filepath"`
	MimeType string `json:"mimeType"`
	Data     []byte `json:"-"`
}

type CoverRef struct {
	Playlist PlaylistID `json:"playlistID"`
	Track    TrackID    `json:"trackID"`
}

type Playlist struct {
	ID       PlaylistID `json:"id"`
	SourceID SourceID   `json:"sourceID"`
	Path     string     `json:"path"`
	Name     string     `json:"name"`

	CoverData *CoverData `json:"coverData"`
	CoverRef  *CoverRef  `json:"coverRef"`
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
	Name     string                 `json:"name"`
	Children []HierarchicalPlaylist `json:"children"`

	TrackCount int        `json:"trackCount"`
	CoverData  *CoverData `json:"coverData"`
	CoverRef   *CoverRef  `json:"coverRef"`
}
