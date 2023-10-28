package handler

import (
	"bytes"
	"context"
	"github.com/gin-gonic/gin"
	"gogs.mikescher.com/BlackForestBytes/goext/ginext"
	"html/template"
	mply "mikescher.com/musicply"
	"mikescher.com/musicply/logic"
	"mikescher.com/musicply/models"
	"net/http"
	"path/filepath"
	"strings"
)

type WebsiteHandler struct {
	app *logic.Application
}

func NewWebsiteHandler(app *logic.Application) WebsiteHandler {
	return WebsiteHandler{
		app: app,
	}
}

func (h WebsiteHandler) ServeIndexHTML(pctx ginext.PreContext) ginext.HTTPResponse {
	ctx, g, errResp := pctx.Start()
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	templ, err := h.app.Assets.Template("index.html", h.buildIndexHTMLTemplate)
	if err != nil {
		return ginext.Error(err)
	}

	data := map[string]any{
		"RemoteIP":    g.RemoteIP(),
		"BranchName":  mply.BranchName,
		"CommitTime":  mply.CommitTime,
		"VCSType":     mply.VCSType,
		"CommitHash":  mply.CommitHash,
		"APILevel":    mply.APILevel,
		"FooterLinks": h.app.Assets.ListFooterLinks(),
	}

	bin := bytes.Buffer{}
	err = templ.Execute(&bin, data)
	if err != nil {
		return ginext.Error(err)
	}

	return ginext.Data(http.StatusOK, "text/html", bin.Bytes())
}

func (h WebsiteHandler) ServeAssets(pctx ginext.PreContext) ginext.HTTPResponse {
	type uri struct {
		FP1 *string `uri:"fp1"`
		FP2 *string `uri:"fp2"`
		FP3 *string `uri:"fp3"`
	}

	var u uri

	ctx, _, errResp := pctx.URI(&u).Start()
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	assetpath := ""
	if u.FP1 == nil && u.FP2 == nil && u.FP3 == nil {
		assetpath = filepath.Join()
	} else if u.FP2 == nil && u.FP3 == nil {
		assetpath = filepath.Join(*u.FP1)
	} else if u.FP3 == nil {
		assetpath = filepath.Join(*u.FP1, *u.FP2)
	} else {
		assetpath = filepath.Join(*u.FP1, *u.FP2, *u.FP3)
	}

	data, err := h.app.Assets.Read(assetpath)
	if err != nil {
		return ginext.JSON(http.StatusNotFound, gin.H{"error": "AssetNotFound", "assetpath": assetpath})
	}

	mime := "text/plain"

	lowerFN := strings.ToLower(assetpath)
	if strings.HasSuffix(lowerFN, ".html") || strings.HasSuffix(lowerFN, ".htm") {
		mime = "text/html"
	} else if strings.HasSuffix(lowerFN, ".css") {
		mime = "text/css"
	} else if strings.HasSuffix(lowerFN, ".js") {
		mime = "text/javascript"
	} else if strings.HasSuffix(lowerFN, ".json") {
		mime = "application/json"
	} else if strings.HasSuffix(lowerFN, ".jpeg") || strings.HasSuffix(lowerFN, ".jpg") {
		mime = "image/jpeg"
	} else if strings.HasSuffix(lowerFN, ".png") {
		mime = "image/png"
	} else if strings.HasSuffix(lowerFN, ".svg") {
		mime = "image/svg+xml"
	}

	return ginext.Data(http.StatusOK, mime, data)
}

func (h WebsiteHandler) buildIndexHTMLTemplate(content []byte) (*template.Template, error) {
	t := template.New("index.html")

	t.Funcs(template.FuncMap{
		"listPlaylists": func() []models.Playlist {
			v, err := h.app.Database.ListPlaylists(context.Background())
			if err != nil {
				panic(err)
			}
			return v
		},
		"listPlaylistTracks": func(plid models.PlaylistID) []models.Track {
			v, err := h.app.Database.ListPlaylistTracks(context.Background(), plid)
			if err != nil {
				panic(err)
			}
			return v
		},
		"listAllTracks": func() []models.Track {
			v, err := h.app.Database.ListTracks(context.Background())
			if err != nil {
				panic(err)
			}
			return v
		},
	})

	_, err := t.Parse(string(content))
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (h WebsiteHandler) GetFooterLinkIcon(pctx ginext.PreContext) ginext.HTTPResponse {
	type uri struct {
		ID models.FooterLinkID `uri:"id"`
	}

	var u uri

	ctx, _, errResp := pctx.URI(&u).Start()
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	fl := h.app.Assets.GetFooterLink(u.ID)

	if fl == nil {
		return ginext.JSON(http.StatusNotFound, gin.H{"error": "AssetNotFound", "id": u.ID})
	}

	mime := "application/octet-stream"

	lowerFN := strings.ToLower(fl.IconPath)
	if strings.HasSuffix(lowerFN, ".jpeg") || strings.HasSuffix(lowerFN, ".jpg") {
		mime = "image/jpeg"
	} else if strings.HasSuffix(lowerFN, ".png") {
		mime = "image/png"
	} else if strings.HasSuffix(lowerFN, ".gif") {
		mime = "image/gif"
	} else if strings.HasSuffix(lowerFN, ".svg") {
		mime = "image/svg+xml"
	}

	return ginext.Data(http.StatusOK, mime, fl.IconData)
}
