package handler

import (
	"gogs.mikescher.com/BlackForestBytes/goext/ginext"
	"mikescher.com/musicply/logic"
	"mikescher.com/musicply/models"
)

type CoverHandler struct {
	app *logic.Application
}

func NewCoverHandler(app *logic.Application) CoverHandler {
	return CoverHandler{
		app: app,
	}
}

func (h CoverHandler) GetCover(pctx ginext.PreContext) ginext.HTTPResponse {
	type uri struct {
		Cover models.CoverHash `uri:"cvrhash"`
	}

	var u uri
	ctx, _, errResp := pctx.URI(&u).Start()
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	cvr, err := h.app.Database.GetCover(ctx, u.Cover)
	if err != nil {
		return ginext.Error(err)
	}

	return ginext.Data(200, cvr.MimeType, cvr.Data)
}
