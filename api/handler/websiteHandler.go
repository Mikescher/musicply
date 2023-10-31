package handler

import (
	"bytes"
	"context"
	"github.com/gin-gonic/gin"
	"gogs.mikescher.com/BlackForestBytes/goext/ginext"
	json "gogs.mikescher.com/BlackForestBytes/goext/gojson"
	template_html "html/template"
	mply "mikescher.com/musicply"
	"mikescher.com/musicply/logic"
	"mikescher.com/musicply/models"
	"mikescher.com/musicply/webassets"
	"net/http"
	"path/filepath"
	template_text "text/template"
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
		"DBChecksum":  h.app.Database.Checksum(),
	}

	bin := bytes.Buffer{}
	err = templ.Execute(&bin, data)
	if err != nil {
		return ginext.Error(err)
	}

	return ginext.Data(http.StatusOK, "text/html", bin.Bytes())
}

func (h WebsiteHandler) ServeScriptJS(pctx ginext.PreContext) ginext.HTTPResponse {
	type uri struct {
		Checksum string `uri:"cs"`
	}

	var u uri
	ctx, g, errResp := pctx.URI(&u).Start()
	if errResp != nil {
		return *errResp
	}
	defer ctx.Cancel()

	csDB := h.app.Database.Checksum() + ".js"

	if csDB != u.Checksum {
		return ginext.JSON(http.StatusNotFound, gin.H{"error": "AssetNotFound", "cs_db": csDB, "cs_uri": u.Checksum, "msg": "wrong checksum"})
	}

	templ, err := h.app.Assets.Template("script.js", h.buildScriptJSTemplate)
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
		"DBChecksum":  h.app.Database.Checksum(),
	}

	bin := bytes.Buffer{}
	err = templ.Execute(&bin, data)
	if err != nil {
		return ginext.Error(err)
	}

	return ginext.Data(http.StatusOK, "text/javascript", bin.Bytes())
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

	mime := mply.FilenameToMime(assetpath, "text/plain")

	return ginext.Data(http.StatusOK, mime, data)
}

func (h WebsiteHandler) buildIndexHTMLTemplate(content []byte) (webassets.ITemplate, error) {
	t := template_html.New("index.html")

	t.Funcs(template_html.FuncMap{
		"listPlaylists": func() models.HierarchicalPlaylist {
			v, err := h.app.ListHierarchicalPlaylists(context.Background())
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
		"safe": func(s string) template_html.HTML { return template_html.HTML(s) }, //nolint:gosec
		"json": func(obj any) string {
			v, err := json.Marshal(obj)
			if err != nil {
				panic(err)
			}
			return string(v)
		},
		"json_indent": func(obj any) string {
			v, err := json.MarshalIndent(obj, "", "  ")
			if err != nil {
				panic(err)
			}
			return string(v)
		},
	})

	_, err := t.Parse(string(content))
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (h WebsiteHandler) buildScriptJSTemplate(content []byte) (webassets.ITemplate, error) {
	t := template_text.New("script.js")

	t = t.Delims("/*{{", "}}*/")

	t.Funcs(template_text.FuncMap{
		"listPlaylists": func() models.HierarchicalPlaylist {
			v, err := h.app.ListHierarchicalPlaylists(context.Background())
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
		"safe": func(s string) template_html.HTML { return template_html.HTML(s) }, //nolint:gosec
		"json": func(obj any) string {
			v, err := json.Marshal(obj)
			if err != nil {
				panic(err)
			}
			return string(v)
		},
		"json_indent": func(obj any) string {
			v, err := json.MarshalIndent(obj, "", "  ")
			if err != nil {
				panic(err)
			}
			return string(v)
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

	mime := mply.FilenameToMime(fl.IconPath, "application/octet-stream")

	return ginext.Data(http.StatusOK, mime, fl.IconData)
}
