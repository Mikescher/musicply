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

func (db *Database) GetTrack(ctx *ginext.AppContext, plid models.PlaylistID, trckid models.TrackID) (models.Track, error) {
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

func (db *Database) GetPlaylist(ctx *ginext.AppContext, plid models.PlaylistID) (models.Playlist, error) {
	db.lock.RLock()
	defer db.lock.RUnlock()

	pl, ok := db.playlists[plid]
	if !ok {
		return models.Playlist{}, exerr.New(mply.ErrEntityNotFound, fmt.Sprintf("playlist '%s' not found", pl)).Build()
	}

	return pl, nil
}

func (db *Database) ListPlaylistTracks(ctx *ginext.AppContext, plid models.PlaylistID) ([]models.Track, error) {
	db.lock.RLock()
	defer db.lock.RUnlock()

	pl, ok := db.playlists[plid]
	if !ok {
		return nil, exerr.New(mply.ErrEntityNotFound, fmt.Sprintf("playlist '%s' not found", pl)).Build()
	}

	r := make([]models.Track, 0, 64)

	for _, track := range db.tracks[plid] {
		r = append(r, track)
	}

	return r, nil
}

func (db *Database) ListPlaylists(ctx *ginext.AppContext) ([]models.Playlist, error) {
	db.lock.RLock()
	defer db.lock.RUnlock()

	r := make([]models.Playlist, 0, len(db.playlists))

	for _, e := range db.playlists {
		r = append(r, e)
	}

	return r, nil
}
