package handler

import (
	"gogs.mikescher.com/BlackForestBytes/goext/ginext"
	"mikescher.com/musicply/logic"
	"mikescher.com/musicply/models"
)

type DatabaseHandler struct {
	app *logic.Application
}

func NewDatabaseHandler(app *logic.Application) DatabaseHandler {
	return DatabaseHandler{
		app: app,
	}
}

func (h DatabaseHandler) RefreshSources(pctx ginext.PreContext) ginext.HTTPResponse {
	type response struct {
		SourcesCount int `json:"sourcesCount"`
	}

	ctx, _, errResp := pctx.Start()
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	sources := h.app.Database.ListSources(ctx)

	for _, src := range sources {
		err := h.app.Database.RefreshSource(src, func(v string) { /*noop*/ })
		if err != nil {
			return ginext.Error(err)
		}
	}

	return ginext.JSON(200, response{SourcesCount: len(sources)})
}

func (h DatabaseHandler) RefreshSingleSource(pctx ginext.PreContext) ginext.HTTPResponse {
	type uri struct {
		SourceID models.SourceID `uri:"sourceid"`
	}
	type response struct {
		Source models.Source `json:"src"`
	}

	var u uri
	ctx, _, errResp := pctx.URI(&u).Start()
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	src, err := h.app.Database.GetSource(ctx, u.SourceID)
	if err != nil {
		return ginext.Error(err)
	}

	err = h.app.Database.RefreshSource(src, func(v string) { /*noop*/ })
	if err != nil {
		return ginext.Error(err)
	}

	return ginext.JSON(200, response{Source: src})
}
