package db

import (
	"mikescher.com/musicply/models"
	"sync"
)

type Database struct {
	sources []models.Source

	playlists map[models.PlaylistID]models.Playlist
	tracks    map[models.PlaylistID]map[models.TrackID]models.Track
	lock      sync.RWMutex
}

func NewDatabase() *Database {
	return &Database{
		sources:   make([]models.Source, 0),
		tracks:    make(map[models.PlaylistID]map[models.TrackID]models.Track),
		playlists: make(map[models.PlaylistID]models.Playlist),
		lock:      sync.RWMutex{},
	}
}
