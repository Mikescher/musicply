package handler

import (
	"gogs.mikescher.com/BlackForestBytes/goext/ginext"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"mikescher.com/musicply/logic"
	"mikescher.com/musicply/models"
)

type PlaylistHandler struct {
	app *logic.Application
}

func NewPlaylistHandler(app *logic.Application) PlaylistHandler {
	return PlaylistHandler{
		app: app,
	}
}

func (h PlaylistHandler) ListPlaylists(pctx ginext.PreContext) ginext.HTTPResponse {
	type query struct {
	}
	type response struct {
		Playlists []models.Playlist `json:"playlists"`
	}

	var q query
	ctx, _, errResp := pctx.Query(&q).Start()
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	playlists, err := h.app.Database.ListPlaylists(ctx)
	if err != nil {
		return ginext.Error(err)
	}

	return ginext.JSON(200, response{Playlists: playlists})
}

func (h PlaylistHandler) ListHierarchicalPlaylists(pctx ginext.PreContext) ginext.HTTPResponse {
	type query struct {
	}

	var q query
	ctx, _, errResp := pctx.Query(&q).Start()
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	root, err := h.app.ListHierarchicalPlaylists(ctx)
	if err != nil {
		return ginext.Error(err)
	}

	return ginext.JSON(200, root)
}

func (h PlaylistHandler) GetPlaylist(pctx ginext.PreContext) ginext.HTTPResponse {
	type uri struct {
		PlaylistId models.PlaylistID `uri:"plid"`
	}

	var u uri
	ctx, _, errResp := pctx.URI(&u).Start()
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	playlist, err := h.app.Database.GetPlaylist(ctx, u.PlaylistId)
	if err != nil {
		return ginext.Error(err)
	}

	return ginext.JSON(200, playlist)
}

func (h PlaylistHandler) ListPlaylistTracks(pctx ginext.PreContext) ginext.HTTPResponse {
	type uri struct {
		PlaylistId models.PlaylistID `uri:"plid"`
	}
	type query struct {
		Search *string `form:"search"`
	}
	type response struct {
		Tracks []models.Track `json:"tracks"`
	}

	var u uri
	var q query
	ctx, _, errResp := pctx.URI(&u).Query(&q).Start()
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	tracks, err := h.app.Database.ListPlaylistTracks(ctx, u.PlaylistId)
	if err != nil {
		return ginext.Error(err)
	}

	if q.Search != nil {
		tracks = langext.ArrFilter(tracks, func(v models.Track) bool { return v.IsFilterMatch(*q.Search) })
	}

	return ginext.JSON(200, response{Tracks: tracks})
}
