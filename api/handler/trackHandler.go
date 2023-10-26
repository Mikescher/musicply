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

func (h TrackHandler) GetPlaylist(pctx ginext.PreContext) ginext.HTTPResponse {
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

func (h TrackHandler) ListPlaylistTracks(pctx ginext.PreContext) ginext.HTTPResponse {
	type uri struct {
		PlaylistId models.PlaylistID `uri:"plid"`
	}
	type query struct {
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

	return ginext.JSON(200, response{Tracks: tracks})
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

	track, err := h.app.Database.GetTrack(ctx, u.PlaylistId, u.TrackId)
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
		PlaylistId models.PlaylistID `uri:"plid"`
		TrackId    models.TrackID    `uri:"trackid"`
	}

	var u uri
	ctx, _, errResp := pctx.URI(&u).Start()
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	track, err := h.app.Database.GetTrack(ctx, u.PlaylistId, u.TrackId)
	if err != nil {
		return ginext.Error(err)
	}

	return ginext.JSON(200, track)
}

func (h TrackHandler) GetTrackCover(pctx ginext.PreContext) ginext.HTTPResponse {
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

	track, err := h.app.Database.GetTrack(ctx, u.PlaylistId, u.TrackId)
	if err != nil {
		return ginext.Error(err)
	}

	if track.Tags.Picture == nil {
		return ginext.Status(404)
	}

	return ginext.Data(200, track.Tags.Picture.MIMEType, track.Tags.Picture.Data)
}

func NewTrackHandler(app *logic.Application) TrackHandler {
	return TrackHandler{
		app: app,
	}
}
