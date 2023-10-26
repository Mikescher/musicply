package api

import (
	"fmt"
	"gogs.mikescher.com/BlackForestBytes/goext/ginext"
	mply "mikescher.com/musicply"
	"mikescher.com/musicply/api/handler"
	"mikescher.com/musicply/logic"
	"mikescher.com/musicply/swagger"
)

type Router struct {
	app *logic.Application

	commonHandler handler.CommonHandler
	trackHandler  handler.TrackHandler
}

func NewRouter(app *logic.Application) *Router {
	return &Router{
		app: app,

		commonHandler: handler.NewCommonHandler(app),
		trackHandler:  handler.NewTrackHandler(app),
	}
}

// Init swaggerdocs
//
//	@title		MusicPly API
//	@version	1.0
//	@host		localhost
//
//	@tag.name	MusicPly
//	@tag.name	Common
//
//	@BasePath	/api/v1/
func (r *Router) Init(e *ginext.GinWrapper) {

	api := e.Routes().Group("/api").Group(fmt.Sprintf("/v%d", mply.APILevel))

	// ================ General ================

	api.Any("/ping").Handle(r.commonHandler.Ping)
	api.GET("/health").Handle(r.commonHandler.Health)
	api.POST("/sleep/:secs").Handle(r.commonHandler.Sleep)

	// ================ Swagger ================

	docs := e.Routes().Group("/documentation")
	{
		docs.GET("/swagger").Handle(ginext.RedirectTemporary("/documentation/swagger/"))
		docs.GET("/swagger/*sub").Handle(swagger.Handle)
	}

	// ================ API ================

	api.GET("/playlists").Handle(r.trackHandler.ListPlaylists)
	api.GET("/playlists/:plid").Handle(r.trackHandler.GetPlaylist)
	api.GET("/playlists/:plid/tracks").Handle(r.trackHandler.ListPlaylistTracks)
	api.GET("/playlists/:plid/tracks/:trackid").Handle(r.trackHandler.GetPlaylistTrack)
	api.GET("/playlists/:plid/tracks/:trackid/stream").Handle(r.trackHandler.StreamPlaylistTrack)
	api.GET("/playlists/:plid/tracks/:trackid/cover").Handle(r.trackHandler.StreamPlaylistTrack)

	api.GET("/tracks").Handle(r.trackHandler.ListTracks)

	api.GET("/tracks/:trackid/stream").Handle(r.trackHandler.StreamTrack)

	// ================  ================

	if r.app.Config.Custom404 {
		e.NoRoute(r.commonHandler.NoRoute)
	}

}
