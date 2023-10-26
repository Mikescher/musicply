package db

import (
	"fmt"
	"gogs.mikescher.com/BlackForestBytes/goext/exerr"
	"gogs.mikescher.com/BlackForestBytes/goext/ginext"
	mply "mikescher.com/musicply"
	"mikescher.com/musicply/models"
)

func (db *Database) ListTracks(ctx *ginext.AppContext) ([]models.Track, error) {
	db.lock.RLock()
	defer db.lock.RUnlock()

	r := make([]models.Track, 0, len(db.tracks)*32)

	for _, playlist := range db.playlists {
		for _, track := range db.tracks[playlist.ID] {
			r = append(r, track)
		}
	}

	return r, nil
}

func (db *Database) GetTrack(plid models.PlaylistID, trckid models.TrackID) (models.Track, error) {
	db.lock.RLock()
	defer db.lock.RUnlock()

	pl, ok := db.playlists[plid]
	if !ok {
		return models.Track{}, exerr.New(mply.ErrEntityNotFound, fmt.Sprintf("playlist '%s' not found", pl)).Build()
	}

	trck, ok := db.tracks[plid][trckid]
	if !ok {
		return models.Track{}, exerr.New(mply.ErrEntityNotFound, fmt.Sprintf("track '%s' not found in playlist '%s'", trckid, pl)).Build()
	}

	return trck, nil
}
