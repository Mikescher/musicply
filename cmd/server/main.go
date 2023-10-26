package main

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"gogs.mikescher.com/BlackForestBytes/goext/ginext"
	mply "mikescher.com/musicply"
	"mikescher.com/musicply/api"
	"mikescher.com/musicply/db"
	"mikescher.com/musicply/logic"
)

func main() {
	mply.Init()

	log.Info().Msg(fmt.Sprintf("Starting with config-namespace <%s>", mply.Conf.Namespace))

	appdb := db.NewDatabase()

	app := logic.NewApp(appdb)

	appdb.LoadSourcesFromEnv("SOURCE")

	appdb.RefreshAllInitial()

	ginengine := ginext.NewEngine(mply.Conf.Cors, mply.Conf.GinDebug, true, mply.Conf.RequestTimeout)

	router := api.NewRouter(app)

	appjobs := make([]logic.Job, 0)

	app.Init(mply.Conf, ginengine, appjobs)

	router.Init(ginengine)

	app.Run()
}
