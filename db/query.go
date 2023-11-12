package db

import (
	"context"
	"fmt"
	"gogs.mikescher.com/BlackForestBytes/goext/exerr"
	"gogs.mikescher.com/BlackForestBytes/goext/ginext"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	mply "mikescher.com/musicply"
	"mikescher.com/musicply/models"
	"sort"
)

func (db *Database) ListTracks(ctx context.Context) ([]models.Track, error) {
	db.lock.RLock()
	defer db.lock.RUnlock()

	playlists := langext.MapValueArr(db.playlists)
	langext.SortBy(playlists, func(v models.Playlist) int { return v.Sort })

	rAll := make([]models.Track, 0, len(db.tracks)*32)

	for _, playlist := range db.playlists {

		rPTracks := langext.MapValueArr(db.tracks[playlist.ID])
		sort.SliceStable(rPTracks, func(i1, i2 int) bool { return models.CompareTracks(playlist.TrackSort, rPTracks[i1], rPTracks[i2]) })

		rAll = append(rAll, rPTracks...)
	}

	return rAll, nil
}

func (db *Database) GetTrack(ctx context.Context, plid models.PlaylistID, trckid models.TrackID) (models.Track, error) {
	db.lock.RLock()
	defer db.lock.RUnlock()

	_, ok := db.playlists[plid]
	if !ok {
		return models.Track{}, exerr.New(mply.ErrEntityNotFound, fmt.Sprintf("playlist '%s' not found", plid)).Build()
	}

	trck, ok := db.tracks[plid][trckid]
	if !ok {
		return models.Track{}, exerr.New(mply.ErrEntityNotFound, fmt.Sprintf("track '%s' not found in playlist '%s'", trckid, plid)).Build()
	}

	return trck, nil
}

func (db *Database) GetTrackDirect(ctx context.Context, trckid models.TrackID) (models.Track, error) {
	db.lock.RLock()
	defer db.lock.RUnlock()

	for _, pltracks := range db.tracks {
		if trck, ok := pltracks[trckid]; ok {
			return trck, nil
		}
	}

	return models.Track{}, exerr.New(mply.ErrEntityNotFound, fmt.Sprintf("track '%s' not found", trckid)).Build()
}

func (db *Database) GetPlaylist(ctx context.Context, plid models.PlaylistID) (models.Playlist, error) {
	db.lock.RLock()
	defer db.lock.RUnlock()

	pl, ok := db.playlists[plid]
	if !ok {
		return models.Playlist{}, exerr.New(mply.ErrEntityNotFound, fmt.Sprintf("playlist '%s' not found", plid)).Build()
	}

	return pl, nil
}

func (db *Database) ListPlaylistTracks(ctx context.Context, plid models.PlaylistID) ([]models.Track, error) {
	db.lock.RLock()
	defer db.lock.RUnlock()

	playlist, ok := db.playlists[plid]
	if !ok {
		return nil, exerr.New(mply.ErrEntityNotFound, fmt.Sprintf("playlist '%s' not found", plid)).Build()
	}

	r := make([]models.Track, 0, 64)

	for _, track := range db.tracks[plid] {
		r = append(r, track)
	}

	sort.SliceStable(r, func(i1, i2 int) bool { return models.CompareTracks(playlist.TrackSort, r[i1], r[i2]) })

	return r, nil
}

func (db *Database) CountPlaylistTracks(ctx context.Context, plid models.PlaylistID) (int, error) {
	db.lock.RLock()
	defer db.lock.RUnlock()

	_, ok := db.playlists[plid]
	if !ok {
		return 0, exerr.New(mply.ErrEntityNotFound, fmt.Sprintf("playlist '%s' not found", plid)).Build()
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

	langext.SortBy(r, func(v models.Playlist) int { return v.Sort })

	return r, nil
}

func (db *Database) GetCover(ctx *ginext.AppContext, cover models.CoverHash) (models.CoverData, error) {
	db.lock.RLock()
	defer db.lock.RUnlock()

	pl, ok := db.covers[cover]
	if !ok {
		return models.CoverData{}, exerr.New(mply.ErrEntityNotFound, fmt.Sprintf("playlist '%s' not found", pl)).Build()
	}

	return pl, nil
}

func (db *Database) ListSources(ctx context.Context) []models.Source {
	db.lock.RLock()
	defer db.lock.RUnlock()

	r := make([]models.Source, 0, len(db.sources))

	r = append(r, db.sources...)

	langext.SortBy(r, func(v models.Source) string { return v.Path })

	return r
}

func (db *Database) GetSource(ctx context.Context, sid models.SourceID) (models.Source, error) {
	db.lock.RLock()
	defer db.lock.RUnlock()

	for _, v := range db.sources {
		if v.ID == sid {
			return v, nil
		}
	}

	return models.Source{}, exerr.New(mply.ErrEntityNotFound, fmt.Sprintf("source '%s' not found", sid)).Build()
}
