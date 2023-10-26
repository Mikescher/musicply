package db

import (
	"gogs.mikescher.com/BlackForestBytes/goext/ginext"
	"mikescher.com/musicply/models"
)

func (db *Database) ListTracks(ctx *ginext.AppContext) ([]models.Track, error) {
	db.lock.RLock()
	defer db.lock.RUnlock()

	r := make([]models.Track, 0, len(db.tracks)*32)

	for _, trackarr := range db.tracks {
		for _, track := range trackarr {
			r = append(r, track)
		}
	}

	return r, nil
}
