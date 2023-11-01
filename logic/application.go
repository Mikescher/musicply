package logic

import (
	"context"
	"github.com/rs/zerolog/log"
	"gogs.mikescher.com/BlackForestBytes/goext/ginext"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"gogs.mikescher.com/BlackForestBytes/goext/syncext"
	mply "mikescher.com/musicply"
	"mikescher.com/musicply/db"
	"mikescher.com/musicply/models"
	"mikescher.com/musicply/webassets"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Application struct {
	Config    mply.Config
	Gin       *ginext.GinWrapper
	Database  *db.Database
	Assets    *webassets.Assets
	Jobs      []Job
	stopChan  chan bool
	Port      string
	IsRunning *syncext.AtomicBool
}

func NewApp(db *db.Database, ass *webassets.Assets) *Application {
	//nolint:exhaustruct
	return &Application{
		Database:  db,
		Assets:    ass,
		stopChan:  make(chan bool),
		IsRunning: syncext.NewAtomicBool(false),
	}
}

func (app *Application) Init(cfg mply.Config, g *ginext.GinWrapper, jobs []Job) {
	app.Config = cfg
	app.Gin = g
	app.Jobs = jobs
}

func (app *Application) Stop() {
	syncext.WriteNonBlocking(app.stopChan, true)
}

func (app *Application) Run() {

	addr := net.JoinHostPort(app.Config.ServerIP, app.Config.ServerPort)

	errChan, httpserver := app.Gin.ListenAndServeHTTP(addr, func(port string) {
		app.Port = port
		app.IsRunning.Set(true)
	})

	sigstop := make(chan os.Signal, 1)
	signal.Notify(sigstop, os.Interrupt, syscall.SIGTERM)

	for _, job := range app.Jobs {
		err := job.Start()
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to start job")
		}
	}

	select {
	case <-sigstop:
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		log.Info().Msg("Stopping HTTP-Server")

		err := httpserver.Shutdown(ctx)

		if err != nil {
			log.Info().Err(err).Msg("Error while stopping the http-server")
		} else {
			log.Info().Msg("Stopped HTTP-Server")
		}

	case err := <-errChan:
		log.Error().Err(err).Msg("HTTP-Server failed")

	case <-app.stopChan:
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		log.Info().Msg("Manually stopping HTTP-Server")

		err := httpserver.Shutdown(ctx)

		if err != nil {
			log.Info().Err(err).Msg("Error while stopping the http-server")
		} else {
			log.Info().Msg("Manually stopped HTTP-Server")
		}
	}

	for _, job := range app.Jobs {
		job.Stop()
	}

	app.IsRunning.Set(false)
}

func (app *Application) ListHierarchicalPlaylists(ctx context.Context) (models.HierarchicalPlaylist, error) {

	playlists, err := app.Database.ListPlaylists(ctx)
	if err != nil {
		return models.HierarchicalPlaylist{}, err
	}

	root := models.HierarchicalPlaylist{Name: "", Children: make([]models.HierarchicalPlaylist, 0)} //nolint:exhaustruct

	getOrCreate := func(path []string) *models.HierarchicalPlaylist {

		c := &root

		currPathStr := ""

		for ppi, pp := range path {

			lastPathPart := ppi == len(path)-1

			currPathStr += "/" + pp

			found := false
			for ci := 0; ci < len(c.Children); ci++ {

				if c.Children[ci].Name == pp {

					if lastPathPart {
						return &c.Children[ci]
					}

					c = &(c.Children[ci])
					found = true
					break
				}

			}
			if !found {
				c.Children = append(c.Children, models.HierarchicalPlaylist{
					ID:         nil,
					HierID:     models.NewHierarchicalPlaylistIDFromSeed("@HIERARCHICAL:" + currPathStr),
					Name:       pp,
					Children:   make([]models.HierarchicalPlaylist, 0),
					Cover:      nil,
					TrackCount: 0,
				})
				if lastPathPart {
					r := &(c.Children[len(c.Children)-1])

					return r
				}
			}
		}

		return c
	}

	for _, plst := range playlists {

		parts := plst.NameParts()

		hp := getOrCreate(parts[:len(parts)-1])

		hplst := models.HierarchicalPlaylist{
			ID:         langext.Ptr(plst.ID),
			HierID:     plst.ID.ToHierarchical(),
			Name:       parts[len(parts)-1],
			Children:   nil,
			Cover:      plst.Cover,
			TrackCount: 0,
		}

		hp.Children = append(hp.Children, hplst)
	}

	var process func(hplst *models.HierarchicalPlaylist) error
	process = func(hplst *models.HierarchicalPlaylist) error {

		tc, err := app.recursiveTrackCount(ctx, *hplst)
		if err != nil {
			return err
		}

		hplst.TrackCount = tc

		for i := range hplst.Children {
			err = process(&hplst.Children[i])
			if err != nil {
				return err
			}
		}

		return nil
	}

	err = process(&root)
	if err != nil {
		return models.HierarchicalPlaylist{}, err
	}

	return root, nil
}

func (app *Application) recursiveTrackCount(ctx context.Context, plst models.HierarchicalPlaylist) (int, error) {
	count := 0
	if plst.ID != nil {
		c, err := app.Database.CountPlaylistTracks(ctx, *plst.ID)
		if err != nil {
			return 0, err
		}
		count += c
	}
	for _, child := range plst.Children {
		c, err := app.recursiveTrackCount(ctx, child)
		if err != nil {
			return 0, err
		}
		count += c
	}

	return count, nil
}
