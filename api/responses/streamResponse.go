package responses

import (
	"github.com/gin-gonic/gin"
	"gogs.mikescher.com/BlackForestBytes/goext/ginext"
)

type headerval struct {
	Key string
	Val string
}

type streamResponse struct {
	mimetype string
	filepath string
	headers  []headerval
}

func (j streamResponse) Write(g *gin.Context) {
	g.Header("Content-Type", j.mimetype) // if we don't set it here gin does weird file-sniffing later...
	for _, v := range j.headers {
		g.Header(v.Key, v.Val)
	}
	g.File(j.filepath)
}

func (j streamResponse) WithHeader(k string, v string) ginext.HTTPResponse {
	j.headers = append(j.headers, headerval{k, v})
	return j
}

func Stream(mimetype string, filepath string) ginext.HTTPResponse {
	return &streamResponse{mimetype: mimetype, filepath: filepath, headers: nil}
}
