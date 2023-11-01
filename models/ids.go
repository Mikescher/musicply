package models

import (
	"gogs.mikescher.com/BlackForestBytes/goext/rext"
)

//go:generate go run ../_gen/id-generate.go -- ids_gen.go

type EntityID interface {
	String() string
	Valid() error
	Prefix() string
	Raw() string
	CheckString() string
	Regex() rext.Regex
}

type JobLogID string //@csid:type [JLG]

type JobExecutionID string //@csid:type [JEX]

type SourceID string //@csid:type [SRC]

type PlaylistID string //@csid:type [PLS]

type HierarchicalPlaylistID string //@csid:type [HPL]

type TrackID string //@csid:type [TRK]

type FooterLinkID string //@csid:type [FLK]

func (id PlaylistID) ToHierarchical() HierarchicalPlaylistID {
	return HierarchicalPlaylistID(generateIDFromSeed(prefixHierarchicalPlaylistID, id.Raw()))
}

func NewHierarchicalPlaylistIDFromSeed(seed string) HierarchicalPlaylistID {
	return HierarchicalPlaylistID(generateIDFromSeed(prefixHierarchicalPlaylistID, seed))
}

func NewPlaylistIDFromPath(path string) PlaylistID {
	return PlaylistID(generateIDFromSeed(prefixPlaylistID, path))
}

func NewTrackIDFromPath(path string) TrackID {
	return TrackID(generateIDFromSeed(prefixTrackID, path))
}
