package main

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"gogs.mikescher.com/BlackForestBytes/goext/ginext"
	mply "mikescher.com/musicply"
	"mikescher.com/musicply/api"
	"mikescher.com/musicply/db"
	"mikescher.com/musicply/jobs"
	"mikescher.com/musicply/logic"
	"mikescher.com/musicply/webassets"
)

func main() {
	mply.Init()

	log.Info().Msg(fmt.Sprintf("Starting with config-namespace <%s>", mply.Conf.Namespace))

	appdb := db.NewDatabase()

	assets := webassets.NewAssets()
	assets.LoadDynamicAssets()

	app := logic.NewApp(appdb, assets)

	appdb.LoadSourcesFromEnv("SOURCES")

	appdb.RefreshAllInitial()

	ginengine := ginext.NewEngine(mply.Conf.Cors, mply.Conf.GinDebug, true, mply.Conf.RequestTimeout)

	router := api.NewRouter(app)

	appjobs := []logic.Job{
		jobs.NewRefreshSourcesJob(app),
	}

	app.Init(mply.Conf, ginengine, appjobs)

	router.Init(ginengine)

	app.Run()
}
