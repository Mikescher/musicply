package jobs

import (
	"context"
	"gogs.mikescher.com/BlackForestBytes/goext/exerr"
	"mikescher.com/musicply/logic"
	"time"
)

type RefreshSourcesData struct {
	//
}

func NewRefreshSourcesJob(app *logic.Application) logic.Job {
	return NewJobRunner(app, "RefreshSourcesJob", 11*time.Hour, 15*time.Minute, runRefreshSources, &RefreshSourcesData{})
}

func runRefreshSources(ctx context.Context, app *logic.Application, lstr *JobListener, _ *RefreshSourcesData) (int, error) {
	sources := app.Database.ListSources(ctx)

	for _, src := range sources {
		err := app.Database.RefreshSource(src, func(v string) { /*noop*/ })
		if err != nil {
			return 0, exerr.Wrap(err, "failed to refresh source").Str("path", src.Path).Build()
		}
	}

	return len(sources), nil
}
