package models

import (
	"gogs.mikescher.com/BlackForestBytes/goext/rext"
)

//go:generate go run ../_gen/id-generate.go -- ids_gen.go

type EntityID interface {
	String() string
	Valid() error
	Prefix() string
	Raw() string
	CheckString() string
	Regex() rext.Regex
}

type JobLogID string //@csid:type [JLG]

type JobExecutionID string //@csid:type [JEX]

type SourceID string //@csid:type [SRC]

type PlaylistID string //@csid:type [PLS]

type TrackID string //@csid:type [TRK]
