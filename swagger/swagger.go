package swagger

import (
	"embed"
	_ "embed"
	"github.com/gin-gonic/gin"
	"gogs.mikescher.com/BlackForestBytes/goext/ginext"
	mply "mikescher.com/musicply"
	"net/http"
	"strings"
)

//go:embed *.html
//go:embed *.json
//go:embed *.yaml
//go:embed *.js
//go:embed *.css
//go:embed themes/*
var assets embed.FS

func getAsset(fn string) ([]byte, string, bool) {
	data, err := assets.ReadFile(fn)
	if err != nil {
		return nil, "", false
	}

	mime := mply.FilenameToMime(fn, "text/plain")

	return data, mime, true
}

func Handle(pctx ginext.PreContext) ginext.HTTPResponse {
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
		index, _, _ := getAsset("index.html")
		return ginext.Data(http.StatusOK, "text/html", index)
	}

	if data, mime, ok := getAsset(u.Filename); ok {
		return ginext.Data(http.StatusOK, mime, data)
	}

	return ginext.JSON(http.StatusNotFound, gin.H{"error": "AssetNotFound", "filename": u.Filename})
}
