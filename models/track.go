package models

import (
	"github.com/dhowden/tag"
	"gogs.mikescher.com/BlackForestBytes/goext/rfctime"
	"io/fs"
)

type Track struct {
	ID        TrackID        `json:"id"`
	SourceID  SourceID       `json:"sourceID"`
	FileMeta  TrackFileMeta  `json:"fileMeta"`
	AudioMeta TrackAudioMeta `json:"audioMeta"`
	Tags      TrackTags      `json:"tags"`
}

type TrackFileMeta struct {
	Path      string                   `json:"path"`
	Filename  string                   `json:"filename"`
	Extension string                   `json:"extension"`
	Size      int64                    `json:"size"`
	Filemode  fs.FileMode              `json:"filemode"`
	ModTime   rfctime.RFC3339NanoTime  `json:"modTime"`
	CTime     *rfctime.RFC3339NanoTime `json:"ctime"`
	ATime     *rfctime.RFC3339NanoTime `json:"atime"`
}

type TrackAudioMeta struct {
	Duration   float64 `json:"duration"`
	BitRate    float64 `json:"bitRate"`
	Channels   int     `json:"channels"`
	CodecShort string  `json:"codecShort"`
	CodecLong  string  `json:"codecLong"`
	Samplerate string  `json:"samplerate"`
}

type TrackTags struct {
	Format      *tag.Format     `json:"format"`
	FileType    *tag.FileType   `json:"fileType"`
	Title       *string         `json:"title"`
	Album       *string         `json:"album"`
	Artist      *string         `json:"artist"`
	AlbumArtist *string         `json:"albumArtist"`
	Composer    *string         `json:"composer"`
	Year        *int            `json:"year"`
	Genre       *string         `json:"genre"`
	TrackIndex  *int            `json:"trackIndex"`
	TrackTotal  *int            `json:"trackTotal"`
	DiscIndex   *int            `json:"discIndex"`
	DiscTotal   *int            `json:"discTotal"`
	Picture     *tag.Picture    `json:"-"`
	Lyrics      *string         `json:"lyrics"`
	Comment     *string         `json:"comment"`
	Raw         *map[string]any `json:"-"`
}
