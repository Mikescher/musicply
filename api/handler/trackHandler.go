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

func (h TrackHandler) ListHierarchicalPlaylists(pctx ginext.PreContext) ginext.HTTPResponse {
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

	if track.Tags.Picture != nil {
		return ginext.Data(200, track.Tags.Picture.MIMEType, track.Tags.Picture.Data)
	}

	plst, err := h.app.Database.GetPlaylist(ctx, u.PlaylistId)
	if err != nil {
		return ginext.Error(err)
	}

	if plst.CoverData != nil {
		return ginext.Data(200, plst.CoverData.MimeType, plst.CoverData.Data)
	}

	if plst.CoverRef != nil {
		reftrack, err := h.app.Database.GetTrack(ctx, plst.CoverRef.Playlist, plst.CoverRef.Track)
		if err != nil {
			return ginext.Error(err)
		}
		if reftrack.Tags.Picture != nil {
			return ginext.Data(200, reftrack.Tags.Picture.MIMEType, reftrack.Tags.Picture.Data)
		}
	}

	return ginext.Data(404, "image/png", h.app.Assets.NoCover())
}

func (h TrackHandler) GetPlaylistCover(pctx ginext.PreContext) ginext.HTTPResponse {
	type uri struct {
		PlaylistId models.PlaylistID `uri:"plid"`
	}

	var u uri
	ctx, _, errResp := pctx.URI(&u).Start()
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	plst, err := h.app.Database.GetPlaylist(ctx, u.PlaylistId)
	if err != nil {
		return ginext.Error(err)
	}

	if plst.CoverData != nil {
		return ginext.Data(200, plst.CoverData.MimeType, plst.CoverData.Data)
	}

	if plst.CoverRef != nil {
		reftrack, err := h.app.Database.GetTrack(ctx, plst.CoverRef.Playlist, plst.CoverRef.Track)
		if err != nil {
			return ginext.Error(err)
		}
		if reftrack.Tags.Picture != nil {
			return ginext.Data(200, reftrack.Tags.Picture.MIMEType, reftrack.Tags.Picture.Data)
		}
	}

	return ginext.Data(404, "image/png", h.app.Assets.NoCover())

}

func NewTrackHandler(app *logic.Application) TrackHandler {
	return TrackHandler{
		app: app,
	}
}
