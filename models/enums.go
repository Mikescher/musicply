package models

//go:generate go run ../_gen/enum-generate.go -- enums_gen.go

type DeDupKey string //@enum:type

const (
	DeDupKeyTitle      DeDupKey = "title"
	DeDupKeyArtist     DeDupKey = "artist"
	DeDupKeyAlbum      DeDupKey = "album"
	DeDupKeyYear       DeDupKey = "year"
	DeDupKeyTrackIndex DeDupKey = "track_index"
	DeDupKeyTrackTotal DeDupKey = "track_total"
	DeDupKeyFilename   DeDupKey = "filename"
)

type DeDupSelector string //@enum:type

const (
	DeDupSelectorAny     DeDupSelector = "any"
	DeDupSelectorNewest  DeDupSelector = "newest"
	DeDupSelectorOldest  DeDupSelector = "oldest"
	DeDupSelectorBiggest DeDupSelector = "biggest"
)
