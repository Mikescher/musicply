package handler

import (
	"gogs.mikescher.com/BlackForestBytes/goext/ginext"
	"mikescher.com/musicply/logic"
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

func (h TrackHandler) GetPlaylistTrack(pctx ginext.PreContext) ginext.HTTPResponse {
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

func (h TrackHandler) StreamPlaylistTrack(pctx ginext.PreContext) ginext.HTTPResponse {
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

func (h TrackHandler) ListTracks(pctx ginext.PreContext) ginext.HTTPResponse {
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

func (h TrackHandler) StreamTrack(pctx ginext.PreContext) ginext.HTTPResponse {
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
