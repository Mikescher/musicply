package models

import (
	"github.com/dhowden/tag"
	"gogs.mikescher.com/BlackForestBytes/goext/rfctime"
	"io/fs"
	"strings"
)

type Track struct {
	ID         TrackID        `json:"id"`
	SourceID   SourceID       `json:"sourceID"`
	PlaylistID PlaylistID     `json:"playlistID"`
	Path       string         `json:"path"`
	FileMeta   TrackFileMeta  `json:"fileMeta"`
	AudioMeta  TrackAudioMeta `json:"audioMeta"`
	Tags       TrackTags      `json:"tags"`
	Cover      *CoverHash     `json:"cover"`
}

func (t Track) Mimetype() string {
	for _, v := range strings.Split(t.AudioMeta.CodecShort, ",") {
		if strings.EqualFold(v, "mp3") {
			return "audio/mpeg"
		}
		if strings.EqualFold(v, "flac") {
			return "audio/flac"
		}
		if strings.EqualFold(v, "m4a") {
			return "audio/mp4"
		}
	}
	if strings.EqualFold(t.FileMeta.Extension, "asf") {
		return "audio/x-ms-wma"
	}

	return "application/octet-stream"
}

func (t Track) IsFilterMatch(v string) bool {
	for _, v := range strings.Split(v, " ") {

		v = strings.ToLower(v)

		if t.Tags.Title != nil && strings.Contains(strings.ToLower(*t.Tags.Title), v) {
			continue
		}

		if t.Tags.Album != nil && strings.Contains(strings.ToLower(*t.Tags.Album), v) {
			continue
		}

		if t.Tags.Artist != nil && strings.Contains(strings.ToLower(*t.Tags.Artist), v) {
			continue
		}

		if t.Tags.AlbumArtist != nil && strings.Contains(strings.ToLower(*t.Tags.AlbumArtist), v) {
			continue
		}

		if strings.Contains(strings.ToLower(t.FileMeta.Filename), v) {
			continue
		}

		return false
	}

	return true
}

type TrackFileMeta struct {
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

func CompareTracks(t1 Track, t2 Track) bool {
	if t1.Tags.Artist != nil && t2.Tags.Artist != nil && *t1.Tags.Artist != *t2.Tags.Artist {
		return *t1.Tags.Artist < *t2.Tags.Artist
	}

	if t1.Tags.Album != nil && t2.Tags.Album != nil && *t1.Tags.Album != *t2.Tags.Album {
		return *t1.Tags.Album < *t2.Tags.Album
	}

	if t1.Tags.TrackIndex != nil && t2.Tags.TrackIndex != nil {

		if t1.Tags.TrackTotal != nil && t2.Tags.TrackTotal != nil && *t1.Tags.TrackTotal == *t2.Tags.TrackTotal && *t1.Tags.TrackIndex != *t2.Tags.TrackIndex {
			return *t1.Tags.TrackIndex < *t2.Tags.TrackIndex
		} else if *t1.Tags.TrackIndex != *t2.Tags.TrackIndex && t1.Tags.TrackTotal == nil && t2.Tags.TrackTotal == nil {
			return *t1.Tags.TrackIndex < *t2.Tags.TrackIndex
		}

	}

	return t1.FileMeta.Filename < t2.FileMeta.Filename
}
