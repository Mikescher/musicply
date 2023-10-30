package mply

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gogs.mikescher.com/BlackForestBytes/goext/exerr"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"os"
)

func Init() {
	//nolint:exhaustruct
	exerr.Init(exerr.ErrorPackageConfigInit{
		ZeroLogErrTraces:       langext.PTrue,
		ZeroLogAllTraces:       langext.PTrue,
		RecursiveErrors:        langext.PTrue,
		ExtendedGinOutput:      &Conf.ReturnRawErrors,
		IncludeMetaInGinOutput: &Conf.ReturnRawErrors,
	})

	//nolint:exhaustruct
	cw := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "2006-01-02 15:04:05 Z07:00",
	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	multi := zerolog.MultiLevelWriter(cw)
	logger := zerolog.New(multi).With().
		Timestamp().
		Caller().
		Logger()

	log.Logger = logger

	if Conf.GinDebug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	zerolog.SetGlobalLevel(Conf.LogLevel)

	log.Debug().Msg("Initialized")
}

var instID xid.ID

func InstanceID() string {
	return instID.String()
}

func init() {
	instID = xid.New()
}
