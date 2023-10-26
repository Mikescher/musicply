package models

import (
	"io/fs"
	"time"
)

type Track struct {
	ID       TrackID
	Path     string
	Size     int64
	Filemode fs.FileMode
	ModTime  time.Time
	CTime    *time.Time
	ATime    *time.Time
}
