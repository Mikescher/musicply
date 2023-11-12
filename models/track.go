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
	Title      string         `json:"title"`
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

func CompareTracks(sortarr []SortKey, t1 Track, t2 Track) bool {

	for _, skey := range sortarr {

		switch skey {

		case SortFilename:
			if t1.FileMeta.Filename != t2.FileMeta.Filename {
				return t1.FileMeta.Filename < t2.FileMeta.Filename
			}

		case SortFilepath:
			if t1.Path != t2.Path {
				return t1.Path < t2.Path
			}

		case SortTitle:
			if t1.Tags.Title != nil && t2.Tags.Title != nil && t1.Tags.Title != t2.Tags.Title {
				return *t1.Tags.Title < *t2.Tags.Title
			}

		case SortArtist:
			if t1.Tags.Artist != nil && t2.Tags.Artist != nil && *t1.Tags.Artist != *t2.Tags.Artist {
				return *t1.Tags.Artist < *t2.Tags.Artist
			}

		case SortAlbum:
			if t1.Tags.Album != nil && t2.Tags.Album != nil && *t1.Tags.Album != *t2.Tags.Album {
				return *t1.Tags.Album < *t2.Tags.Album
			}

		case SortTrackIndex:
			if t1.Tags.TrackIndex != nil && t2.Tags.TrackIndex != nil {
				if t1.Tags.TrackTotal != nil && t2.Tags.TrackTotal != nil && *t1.Tags.TrackTotal == *t2.Tags.TrackTotal && *t1.Tags.TrackIndex != *t2.Tags.TrackIndex {
					return *t1.Tags.TrackIndex < *t2.Tags.TrackIndex
				} else if *t1.Tags.TrackIndex != *t2.Tags.TrackIndex && t1.Tags.TrackTotal == nil && t2.Tags.TrackTotal == nil {
					return *t1.Tags.TrackIndex < *t2.Tags.TrackIndex
				}
			}

		case SortYear:
			if t1.Tags.Year != nil && t2.Tags.Year != nil && t1.Tags.Year != t2.Tags.Year {
				return *t1.Tags.Year < *t2.Tags.Year
			}

		case SortFileMDate:
			if t1.FileMeta.ModTime != t2.FileMeta.ModTime {
				return t1.FileMeta.ModTime.Unix() < t2.FileMeta.ModTime.Unix()
			}

		case SortFileCDate:
			if t1.FileMeta.CTime != nil && t2.FileMeta.CTime != nil && t1.FileMeta.CTime != t2.FileMeta.CTime {
				return t1.FileMeta.CTime.Unix() < t2.FileMeta.CTime.Unix()
			}

		case SortFileADate:
			if t1.FileMeta.ATime != nil && t2.FileMeta.ATime != nil && t1.FileMeta.ATime != t2.FileMeta.ATime {
				return t1.FileMeta.ATime.Unix() < t2.FileMeta.ATime.Unix()
			}

		default:
			panic("unknown sort-key: " + skey)

		}

	}

	return false

}
