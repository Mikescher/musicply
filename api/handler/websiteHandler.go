package handler

import (
	"bytes"
	"context"
	"github.com/gin-gonic/gin"
	"gogs.mikescher.com/BlackForestBytes/goext/ginext"
	"html/template"
	"mikescher.com/musicply/logic"
	"mikescher.com/musicply/models"
	"net/http"
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
		"RemoteIP": g.RemoteIP(),
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
		Filename string `uri:"sub"`
	}

	var u uri

	ctx, _, errResp := pctx.URI(&u).Start()
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	u.Filename = strings.TrimLeft(u.Filename, "/")

	if u.Filename == "" {
		u.Filename = "index.html"
	}

	data, err := h.app.Assets.Read(u.Filename)
	if err != nil {
		return ginext.JSON(http.StatusNotFound, gin.H{"error": "AssetNotFound", "filename": u.Filename})
	}

	mime := "text/plain"

	lowerFN := strings.ToLower(u.Filename)
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
