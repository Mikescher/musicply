package mply

import (
	"gogs.mikescher.com/BlackForestBytes/goext/exerr"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
)

var (
	ErrInternal       = exerr.TypeInternal
	ErrPanic          = exerr.TypePanic
	ErrWrap           = exerr.TypeWrap
	ErrNotImplemented = exerr.TypeNotImplemented

	ErrBindFailURI      = exerr.TypeBindFailURI
	ErrBindFailQuery    = exerr.TypeBindFailQuery
	ErrBindFailJSON     = exerr.TypeBindFailJSON
	ErrBindFailFormData = exerr.TypeBindFailFormData

	ErrUnauthorized = exerr.TypeUnauthorized
	ErrAuthFailed   = exerr.TypeAuthFailed

	ErrEntityNotFound = exerr.NewType("NOT_FOUND", langext.Ptr(400))

	ErrJob            = exerr.NewType("JOB", langext.Ptr(500))
	ErrSourceNotFound = exerr.NewType("SOURCE_NOT_FOUND", langext.Ptr(400))
	ErrConfig         = exerr.NewType("CONFIG", langext.Ptr(400))

	ErrEnviron = exerr.NewType("ENVIRON", langext.Ptr(400))
)
