package handler

import (
	"gogs.mikescher.com/BlackForestBytes/goext/ginext"
	"mikescher.com/musicply/api/responses"
	"mikescher.com/musicply/logic"
	"mikescher.com/musicply/models"
)

type TrackHandler struct {
	app *logic.Application
}

func (h TrackHandler) ListPlaylists(pctx ginext.PreContext) ginext.HTTPResponse {
	type query struct {
	}

	var q query
	ctx, _, errResp := pctx.Query(&q).Start()
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	return ginext.NotImplemented()
}

func (h TrackHandler) GetPlaylist(pctx ginext.PreContext) ginext.HTTPResponse {
	type uri struct {
	}

	var u uri
	ctx, _, errResp := pctx.URI(&u).Start()
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	return ginext.NotImplemented()
}

func (h TrackHandler) ListPlaylistTracks(pctx ginext.PreContext) ginext.HTTPResponse {
	type uri struct {
	}
	type query struct {
	}

	var u uri
	var q query
	ctx, _, errResp := pctx.URI(&u).Query(&q).Start()
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	return ginext.NotImplemented()
}

func (h TrackHandler) StreamTrack(pctx ginext.PreContext) ginext.HTTPResponse {
	type uri struct {
		PlaylistId models.PlaylistID `uri:"plid"`
		TrackId    models.TrackID    `uri:"trackid"`
	}

	var u uri
	ctx, _, errResp := pctx.URI(&u).Start()
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	track, err := h.app.Database.GetTrack(u.PlaylistId, u.TrackId)
	if err != nil {
		return ginext.Error(err)
	}

	return responses.Stream(track.Mimetype(), track.Path)
}

func (h TrackHandler) ListTracks(pctx ginext.PreContext) ginext.HTTPResponse {
	type query struct {
		Search string `form:"search"` //TODO
	}
	type response struct {
		Tracks []models.Track `json:"tracks"`
	}

	var q query
	ctx, _, errResp := pctx.Query(&q).Start()
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	tracks, err := h.app.Database.ListTracks(ctx)
	if err != nil {
		return ginext.Error(err)
	}

	return ginext.JSON(200, response{Tracks: tracks})
}

func (h TrackHandler) GetTrack(pctx ginext.PreContext) ginext.HTTPResponse {
	type uri struct {
	}

	var u uri
	ctx, _, errResp := pctx.URI(&u).Start()
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	return ginext.NotImplemented()
}

func NewTrackHandler(app *logic.Application) TrackHandler {
	return TrackHandler{
		app: app,
	}
}
