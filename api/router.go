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

	commonHandler   handler.CommonHandler
	trackHandler    handler.TrackHandler
	playlistHandler handler.PlaylistHandler
	coverHandler    handler.CoverHandler
	websiteHandler  handler.WebsiteHandler
	databaseHandler handler.DatabaseHandler
}

func NewRouter(app *logic.Application) *Router {
	return &Router{
		app: app,

		commonHandler:   handler.NewCommonHandler(app),
		trackHandler:    handler.NewTrackHandler(app),
		playlistHandler: handler.NewPlaylistHandler(app),
		coverHandler:    handler.NewCoverHandler(app),
		websiteHandler:  handler.NewWebsiteHandler(app),
		databaseHandler: handler.NewDatabaseHandler(app),
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

	// ================ Website ================

	website := e.Routes().Group("")
	{
		website.GET("/").Handle(r.websiteHandler.ServeIndexHTML)
		website.GET("/index.html").Handle(r.websiteHandler.ServeIndexHTML)
		website.GET("/scripts/:cs").Handle(r.websiteHandler.ServeScriptJS)
		website.GET("/:fp1").Handle(r.websiteHandler.ServeAssets)
		website.GET("/:fp1/:fp2").Handle(r.websiteHandler.ServeAssets)
		website.GET("/:fp1/:fp2/:fp3").Handle(r.websiteHandler.ServeAssets)

		website.GET("/footerlinks/:id/icon").Handle(r.websiteHandler.GetFooterLinkIcon)
	}

	// ================ API ================

	api.GET("/playlists").Handle(r.playlistHandler.ListPlaylists)
	api.GET("/playlists/hierarchical").Handle(r.playlistHandler.ListHierarchicalPlaylists)
	api.GET("/playlists/:plid").Handle(r.playlistHandler.GetPlaylist)
	api.GET("/playlists/:plid/tracks").Handle(r.playlistHandler.ListPlaylistTracks)
	api.GET("/playlists/:plid/tracks/:trackid").Handle(r.trackHandler.GetTrack)
	api.GET("/playlists/:plid/tracks/:trackid/stream").Handle(r.trackHandler.StreamTrack)

	api.GET("/tracks").Handle(r.trackHandler.ListTracks)
	api.GET("/tracks/:trackid").Handle(r.trackHandler.GetTrackDirect)

	api.GET("/covers/:cvrhash").Handle(r.coverHandler.GetCover)

	api.POST("/refresh").Handle(r.databaseHandler.RefreshSources)

	// ================  ================

	if r.app.Config.Custom404 {
		e.NoRoute(r.commonHandler.NoRoute)
	}

}
