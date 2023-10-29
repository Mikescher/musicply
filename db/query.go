package db

import (
	"context"
	"fmt"
	"gogs.mikescher.com/BlackForestBytes/goext/exerr"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	mply "mikescher.com/musicply"
	"mikescher.com/musicply/models"
	"sort"
)

func (db *Database) ListTracks(ctx context.Context) ([]models.Track, error) {
	db.lock.RLock()
	defer db.lock.RUnlock()

	r := make([]models.Track, 0, len(db.tracks)*32)

	for _, playlist := range db.playlists {
		for _, track := range db.tracks[playlist.ID] {
			r = append(r, track)
		}
	}

	sort.SliceStable(r, func(i1, i2 int) bool { return models.CompareTracks(r[i1], r[i2]) })

	return r, nil
}

func (db *Database) GetTrack(ctx context.Context, plid models.PlaylistID, trckid models.TrackID) (models.Track, error) {
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

func (db *Database) GetPlaylist(ctx context.Context, plid models.PlaylistID) (models.Playlist, error) {
	db.lock.RLock()
	defer db.lock.RUnlock()

	pl, ok := db.playlists[plid]
	if !ok {
		return models.Playlist{}, exerr.New(mply.ErrEntityNotFound, fmt.Sprintf("playlist '%s' not found", pl)).Build()
	}

	return pl, nil
}

func (db *Database) ListPlaylistTracks(ctx context.Context, plid models.PlaylistID) ([]models.Track, error) {
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

	sort.SliceStable(r, func(i1, i2 int) bool { return models.CompareTracks(r[i1], r[i2]) })

	return r, nil
}

func (db *Database) CountPlaylistTracks(ctx context.Context, plid models.PlaylistID) (int, error) {
	db.lock.RLock()
	defer db.lock.RUnlock()

	pl, ok := db.playlists[plid]
	if !ok {
		return 0, exerr.New(mply.ErrEntityNotFound, fmt.Sprintf("playlist '%s' not found", pl)).Build()
	}

	return len(db.tracks[plid]), nil
}

func (db *Database) ListPlaylists(ctx context.Context) ([]models.Playlist, error) {
	db.lock.RLock()
	defer db.lock.RUnlock()

	r := make([]models.Playlist, 0, len(db.playlists))

	for _, e := range db.playlists {
		r = append(r, e)
	}

	langext.SortBy(r, func(v models.Playlist) string { return v.Name })

	return r, nil
}
