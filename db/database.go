package db

import (
	"mikescher.com/musicply/models"
	"sync"
)

type Database struct {
	sources []models.Source

	tracks map[models.SourceID][]models.Track
	lock   sync.RWMutex
}

func NewDatabase() *Database {
	return &Database{
		tracks: make(map[models.SourceID][]models.Track),
	}
}
