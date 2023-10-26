package models

import (
	"github.com/dhowden/tag"
	"io/fs"
	"time"
)

type Track struct {
	ID        TrackID
	FileMeta  TrackFileMeta
	AudioMeta TrackAudioMeta
	Tags      TrackTags
}

type TrackFileMeta struct {
	Path      string
	Filename  string
	Extension string
	Size      int64
	Filemode  fs.FileMode
	ModTime   time.Time
	CTime     *time.Time
	ATime     *time.Time
}

type TrackAudioMeta struct {
	Duration   float64
	BitRate    float64
	Channels   int
	CodecShort string
	CodecLong  string
	Samplerate string
}

type TrackTags struct {
	Format      tag.Format
	FileType    tag.FileType
	Title       string
	Album       string
	Artist      string
	AlbumArtist string
	Composer    string
	Year        int
	Genre       string
	TrackIndex  int
	TrackTotal  int
	DiscIndex   int
	DiscTotal   int
	Picture     *tag.Picture
	Lyrics      string
	Comment     string
	Raw         map[string]any
}
