package models

import (
	"fmt"
	json "gogs.mikescher.com/BlackForestBytes/goext/gojson"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
)

type Source struct {
	ID            SourceID
	SortIndex     int
	Name          string
	Path          string
	Recursive     bool
	Deduplication *DedupliationConfig
}

type DedupliationConfig struct {
	Keys     []DeDupKey
	Selector DeDupSelector
}

func (c DedupliationConfig) genkey(track Track) string {
	r := make([]any, 0)
	for _, key := range c.Keys {
		switch key {
		case DeDupKeyTitle:
			r = append(r, track.Tags.Title)
		case DeDupKeyArtist:
			r = append(r, track.Tags.Artist)
		case DeDupKeyAlbum:
			r = append(r, track.Tags.Album)
		case DeDupKeyYear:
			r = append(r, track.Tags.Year)
		case DeDupKeyTrackIndex:
			r = append(r, track.Tags.TrackIndex)
		case DeDupKeyTrackTotal:
			r = append(r, track.Tags.TrackTotal)
		case DeDupKeyFilename:
			r = append(r, track.FileMeta.Filename)
		default:
			panic("unknown dedup-key: " + key)
		}
	}
	b, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}

	return string(b)
}

func (c DedupliationConfig) getPrimary(tracks []Track) int {
	if c.Selector == DeDupSelectorAny {
		return 0
	}

	if c.Selector == DeDupSelectorOldest {
		idx := 0
		val := langext.Coalesce(tracks[0].FileMeta.CTime, tracks[0].FileMeta.ModTime).Unix()

		for i, track := range tracks {
			loopval := langext.Coalesce(track.FileMeta.CTime, track.FileMeta.ModTime).Unix()
			if loopval < val {
				val = loopval
				idx = i
			}
		}

		return idx
	}

	if c.Selector == DeDupSelectorNewest {
		idx := 0
		val := langext.Coalesce(tracks[0].FileMeta.CTime, tracks[0].FileMeta.ModTime).Unix()

		for i, track := range tracks {
			loopval := langext.Coalesce(track.FileMeta.CTime, track.FileMeta.ModTime).Unix()
			if loopval > val {
				val = loopval
				idx = i
			}
		}

		return idx
	}

	if c.Selector == DeDupSelectorBiggest {
		idx := 0
		val := tracks[0].FileMeta.Size

		for i, track := range tracks {
			loopval := track.FileMeta.Size
			if loopval > val {
				val = loopval
				idx = i
			}
		}

		return idx
	}

	panic("unknown dedup-selector: " + c.Selector)
}

func (c DedupliationConfig) Apply(alltracks []Track) []Track {

	duplicates := make(map[TrackID]bool, len(alltracks))

	d := make(map[string][]Track)

	for _, track := range alltracks {
		d[c.genkey(track)] = append(d[c.genkey(track)], track)
	}

	for _, duptracks := range d {
		if len(duptracks) > 1 {
			prim := c.getPrimary(duptracks)
			for i := 0; i < len(duptracks); i++ {
				if prim != i {
					duplicates[duptracks[i].ID] = true
					fmt.Printf("[DEDUPLICATE] Skip track %s in favor of %s: %s - %s\n", duptracks[i].ID, duptracks[prim].ID, langext.Coalesce(duptracks[i].Tags.Artist, "(NULL)"), duptracks[i].Title)
				}
			}
		}
	}

	result := make([]Track, 0, len(alltracks))
	for _, t := range alltracks {
		if _, ok := duplicates[t.ID]; !ok {
			result = append(result, t)
		}
	}

	return result
}
