package db

import (
	"errors"
	"fmt"
	"github.com/dhowden/tag"
	"github.com/rs/zerolog/log"
	"github.com/titanous/json5"
	"github.com/vansante/go-ffprobe"
	"gogs.mikescher.com/BlackForestBytes/goext/cryptext"
	"gogs.mikescher.com/BlackForestBytes/goext/dataext"
	"gogs.mikescher.com/BlackForestBytes/goext/exerr"
	"gogs.mikescher.com/BlackForestBytes/goext/fsext"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"gogs.mikescher.com/BlackForestBytes/goext/rfctime"
	"io/fs"
	mply "mikescher.com/musicply"
	"mikescher.com/musicply/models"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func (db *Database) LoadSourcesFromEnv(envkey string) {

	fmt.Printf("\n")
	fmt.Printf("================ INITIALIZE SOURCES ================\n")
	fmt.Printf("\n")

	envval, ok := os.LookupEnv(envkey)
	if !ok {
		exerr.New(mply.ErrConfig, "No config specified (via "+envkey+" environment variable)").Fatal()
	}

	if isConfigFilepath(envval) {
		log.Info().Msg("Load config from file '" + envval + "'")
		v, err := os.ReadFile(envval)
		if err != nil {
			exerr.Wrap(err, fmt.Sprintf("failed to load config from file '%s'", envval)).Str("envval", envval).Fatal()
		}
		envval = string(v)
	}

	type SourceSpec struct {
		Name      *string `json:"name"`
		Path      *string `json:"path"`
		Recursive *bool   `json:"recursive"`
	}

	spec := make([]SourceSpec, 0)

	err := json5.Unmarshal([]byte(envval), &spec)
	if err != nil {
		exerr.Wrap(err, "failed to unmarshal config (from  "+envkey+" environment variable)").Fatal()
	}

	for _, srcspec := range spec {
		if srcspec.Path == nil {
			exerr.New(mply.ErrConfig, "missing path in source config").Fatal()
		}

		exist, err := fsext.DirectoryExists(*srcspec.Path)
		if err != nil {
			exerr.Wrap(err, fmt.Sprintf("failed to access directory '%s'", *srcspec.Path)).Str("path", *srcspec.Path).Fatal()
		}
		if !exist {
			exerr.New(mply.ErrSourceNotFound, fmt.Sprintf("directory '%s' does not exist", *srcspec.Path)).Str("path", *srcspec.Path).Fatal()
		}

		var name = ""

		if srcspec.Name == nil {
			basename := path.Base(*srcspec.Path)

			name = basename

			for i := 1; langext.ArrAny(db.sources, func(src models.Source) bool { return src.Name == name }); i++ {
				name = fmt.Sprintf("%s (%d)", basename, i)
			}
		} else {
			name = *srcspec.Name
		}

		if langext.ArrAny(db.sources, func(src models.Source) bool { return src.Name == name }) {
			exerr.New(mply.ErrSourceNotFound, fmt.Sprintf("Duplicate source name '%s'", name)).Str("name", name).Fatal()
		}

		src := models.Source{
			ID:        models.NewSourceID(),
			Name:      name,
			Path:      *srcspec.Path,
			Recursive: langext.Coalesce(srcspec.Recursive, false),
		}

		db.sources = append(db.sources, src)

		if src.Recursive {
			fmt.Printf("Initialize '%s' from source: '%s' (recursive)\n", src.Name, src.Path)
		} else {
			fmt.Printf("Initialize '%s' from source: '%s' (flat)\n", src.Name, src.Path)
		}

	}

	fmt.Printf("\n")
	fmt.Printf("================ ================== ================\n")
	fmt.Printf("\n")
}

func isConfigFilepath(v string) bool {
	v = strings.Trim(v, " \r\t\n")
	if strings.HasPrefix(v, "{") {
		return false
	}
	if strings.HasPrefix(v, "[") {
		return false
	}
	if strings.ContainsAny(v, "\r\n\t") {
		return false
	}
	if strings.HasPrefix(v, "/") {
		return true
	}

	tmp := make([]string, 0)
	if err := json5.Unmarshal([]byte(v), &tmp); err == nil {
		return false
	}

	if len(v) > 2048 {
		return false
	}

	return true
}

func (db *Database) RefreshAllInitial() {

	fmt.Printf("\n")
	fmt.Printf("================ ENUMERATE SOURCES ================\n")
	fmt.Printf("\n")

	for _, src := range db.sources {

		err := db.RefreshSource(src)
		if err != nil {
			exerr.Wrap(err, "").Fatal()
		}
	}

	fmt.Printf("\n")
	fmt.Printf("================ ================== ================\n")
	fmt.Printf("\n")

}

func (db *Database) RefreshSource(src models.Source) error {

	files := make([]dataext.Tuple[string, fs.FileInfo], 0)

	if src.Recursive {
		err := filepath.Walk(src.Path, func(path string, info fs.FileInfo, err error) error {
			if !info.IsDir() {
				files = append(files, dataext.Tuple[string, fs.FileInfo]{V1: path, V2: info})
			}
			return nil
		})
		if err != nil {
			return err
		}
	} else {
		entries, err := os.ReadDir(src.Path)
		if err != nil {
			return err
		}
		for _, v := range entries {
			if !v.IsDir() {
				info, err := v.Info()
				if err != nil {
					return err
				}
				files = append(files, dataext.Tuple[string, fs.FileInfo]{V1: filepath.Join(src.Path, v.Name()), V2: info})
			}
		}
	}

	potentialCovers := langext.ArrFilter(files, func(v dataext.Tuple[string, fs.FileInfo]) bool {
		return langext.InArray(filepath.Base(strings.ToLower(v.V1)), []string{
			"cover.png", "cover.jpeg", "cover.jpg", "cover.bmp", "cover.gif", "cover.webm",
			"folder.png", "folder.jpeg", "folder.jpg", "folder.bmp", "folder.gif", "folder.webm",
			"albumart.png", "albumart.jpeg", "albumart.jpg", "albumart.bmp", "albumart.gif", "albumart.webm",
			"albumartsmall.png", "albumartsmall.jpeg", "albumartsmall.jpg", "albumartsmall.bmp", "albumartsmall.gif", "albumartsmall.webm",
			"front.png", "front.jpeg", "front.jpg", "front.bmp", "front.gif", "front.webm",
		})
	})

	files = langext.ArrFilter(files, func(v dataext.Tuple[string, fs.FileInfo]) bool {
		return langext.InArray(filepath.Ext(strings.ToLower(v.V1)), []string{".mp3", ".flac", ".m4a", ".ogg", ".wav", ".wma"})
	})

	covers := make([]models.CoverData, 0)

	tracks := make([]models.Track, 0, len(files))
	for _, fobj := range files {

		fp := fobj.V1
		info := fobj.V2

		track, coverdata, err := db.analyzeAudioFile(src.ID, fp, info)
		if err != nil {
			log.Error().Msg(fmt.Sprintf("Failed to parse track from file '%s'", fp))
			exerr.Wrap(err, "").Print()
			continue
		}

		tracks = append(tracks, track)
		if coverdata != nil {
			covers = append(covers, *coverdata)
		}
	}

	var fileCover *models.CoverData = nil
	if len(potentialCovers) > 0 {
		mime := mply.FilenameToMime(potentialCovers[0].V1, "")
		if mime != "" {
			bin, err := os.ReadFile(potentialCovers[0].V1)
			if err != nil {
				log.Error().Msg(fmt.Sprintf("Failed to parse load cover file '%s'", potentialCovers[0].V1))
				exerr.Wrap(err, "").Fatal()
			}
			fileCover = &models.CoverData{
				Hash:     models.CoverHash(cryptext.BytesSha256(bin)),
				MimeType: mime,
				Data:     bin,
			}
			covers = append(covers, *fileCover)
		}

	}

	db.lock.Lock()
	defer db.lock.Unlock()

	playlists := make([]dataext.Tuple[models.Playlist, []models.Track], 0)
	if len(tracks) > 0 {

		existing := langext.ArrFirstOrNil(langext.MapValueArr(db.playlists), func(pl models.Playlist) bool { return pl.SourceID == src.ID })

		plid := models.NewPlaylistID()

		existingTrackmap := make(map[string]models.TrackID)
		for _, v := range db.tracks[plid] {
			existingTrackmap[v.Path] = v.ID
		}

		if existing != nil {
			plid = existing.ID
		}

		plst := models.Playlist{
			ID:       plid,
			SourceID: src.ID,
			Name:     src.Name,
			Path:     src.Path,
			Cover:    langext.ConditionalFn01(fileCover == nil, nil, func() *models.CoverHash { return langext.Ptr(fileCover.Hash) }),
		}

		pltracks := tracks
		for i := 0; i < len(pltracks); i++ {
			if v, ok := existingTrackmap[pltracks[i].Path]; ok {
				pltracks[i].ID = v
			}
			pltracks[i].PlaylistID = plst.ID
		}

		sort.SliceStable(pltracks, func(i1, i2 int) bool { return models.CompareTracks(pltracks[i1], pltracks[i2]) })

		if plst.Cover == nil {
			coverTrack := langext.ArrFirstOrNil(pltracks, func(v models.Track) bool { return v.Tags.Picture != nil })
			if coverTrack != nil {
				plst.Cover = coverTrack.Cover
			}
		}

		if plst.Cover != nil {
			for i := 0; i < len(pltracks); i++ {
				if pltracks[i].Cover == nil {
					pltracks[i].Cover = plst.Cover
				}
			}
		}

		playlists = append(playlists, dataext.Tuple[models.Playlist, []models.Track]{V1: plst, V2: pltracks})
	}

	for _, cvrdata := range covers {
		db.covers[cvrdata.Hash] = cvrdata
	}

	for _, v := range langext.ArrFilter(langext.MapValueArr(db.playlists), func(pl models.Playlist) bool { return pl.SourceID == src.ID }) {
		delete(db.playlists, v.ID)
		delete(db.tracks, v.ID)
	}

	fmt.Printf("Found %d tracks and %d playlists in source '%s'\n", len(tracks), len(playlists), src.Name)

	for _, v := range playlists {
		db.playlists[v.V1.ID] = v.V1
		db.tracks[v.V1.ID] = langext.ArrToMap(v.V2, func(v models.Track) models.TrackID { return v.ID })
	}

	db.recalcChecksum(false)

	return nil
}

func (db *Database) analyzeAudioFile(srcid models.SourceID, fp string, info fs.FileInfo) (models.Track, *models.CoverData, error) {

	fm, err := db.getFileMeta(fp, info)
	if err != nil {
		return models.Track{}, nil, exerr.Wrap(err, "").Build()
	}

	am, err := db.getAudioMeta(fp)
	if err != nil {
		return models.Track{}, nil, exerr.Wrap(err, "").Build()
	}

	tt, err := db.getTrackTags(fp)
	if err != nil {
		return models.Track{}, nil, exerr.Wrap(err, "").Build()
	}

	var cvr *models.CoverData = nil
	var chash *models.CoverHash = nil
	if tt.Picture != nil {
		chash = langext.Ptr(models.CoverHash(cryptext.BytesSha256(tt.Picture.Data)))
		cvr = &models.CoverData{
			Hash:     *chash,
			MimeType: tt.Picture.MIMEType,
			Data:     tt.Picture.Data,
		}
	}

	return models.Track{
		ID:         models.NewTrackID(),
		SourceID:   srcid,
		Path:       fp,
		FileMeta:   fm,
		AudioMeta:  am,
		Tags:       tt,
		Cover:      chash,
		PlaylistID: "", // will be set later
	}, cvr, nil
}

func (db *Database) getFileMeta(fp string, info fs.FileInfo) (models.TrackFileMeta, error) {
	if info == nil {
		return models.TrackFileMeta{}, exerr.New(mply.ErrInternal, "no file info").Build()
	}

	var ctime *time.Time = nil
	var atime *time.Time = nil

	if v, ok := info.Sys().(*syscall.Stat_t); ok {
		ctime = langext.Ptr(time.Unix(v.Ctim.Sec, v.Ctim.Nsec))
		atime = langext.Ptr(time.Unix(v.Atim.Sec, v.Atim.Nsec))
	}

	conv := func(t *time.Time) *rfctime.RFC3339NanoTime {
		if t == nil {
			return nil
		}
		return langext.Ptr(rfctime.NewRFC3339Nano(*t))
	}

	return models.TrackFileMeta{
		Filename:  filepath.Base(fp),
		Extension: strings.TrimLeft(filepath.Ext(strings.ToLower(filepath.Base(fp))), "."),
		Size:      info.Size(),
		Filemode:  info.Mode(),
		ModTime:   rfctime.NewRFC3339Nano(info.ModTime()),
		CTime:     conv(ctime),
		ATime:     conv(atime),
	}, nil
}

func (db *Database) getAudioMeta(fp string) (models.TrackAudioMeta, error) {

	pdata, err := ffprobe.GetProbeData(fp, 5*time.Second)
	if err != nil {
		return models.TrackAudioMeta{}, exerr.Wrap(err, "failed to get audiofile ffprobe data").Build()
	}

	if pdata.Format == nil {
		return models.TrackAudioMeta{}, exerr.Wrap(err, "failed to get audiofile ffprobe data (no format)").Build()
	}

	astream := pdata.GetFirstAudioStream()
	if astream == nil {
		return models.TrackAudioMeta{}, exerr.Wrap(err, "failed to get audiofile ffprobe data (no audio stream)").Build()
	}

	br, err := strconv.ParseFloat(pdata.Format.BitRate, 64)
	if err != nil {
		return models.TrackAudioMeta{}, exerr.Wrap(err, "failed to get bitrate from ffprobe data").Str("astream.BitRate", astream.BitRate).Build()
	}

	return models.TrackAudioMeta{
		Duration:   pdata.Format.Duration().Seconds(),
		BitRate:    br,
		Channels:   astream.Channels,
		CodecShort: astream.CodecName,
		CodecLong:  astream.CodecLongName,
		Samplerate: astream.SampleRate,
	}, nil
}

func (db *Database) getTrackTags(fp string) (models.TrackTags, error) {

	fptr, err := os.OpenFile(fp, os.O_RDONLY, 0755)
	if err != nil {
		return models.TrackTags{}, err
	}

	defer func() { _ = fptr.Close() }()

	md, err := tag.ReadFrom(fptr)
	if errors.Is(err, tag.ErrNoTagsFound) {
		return models.TrackTags{
			Format:      nil,
			FileType:    nil,
			Title:       nil,
			Album:       nil,
			Artist:      nil,
			AlbumArtist: nil,
			Composer:    nil,
			Year:        nil,
			Genre:       nil,
			TrackIndex:  nil,
			TrackTotal:  nil,
			DiscIndex:   nil,
			DiscTotal:   nil,
			Picture:     nil,
			Lyrics:      nil,
			Comment:     nil,
			Raw:         nil,
		}, nil
	}
	if err != nil {
		return models.TrackTags{}, exerr.Wrap(err, "failed to get audiofile tag-data").Build()
	}

	estrptr := func(x string) *string {
		if x == "" {
			return nil
		} else {
			return &x
		}
	}

	var mdTrackIndex *int = nil
	var mdTrackTotal *int = nil
	if a, b := md.Track(); a != 0 || b != 0 {
		mdTrackIndex = &a
		mdTrackTotal = &b
	}

	var mdDiscIndex *int = nil
	var mdDiscTotal *int = nil
	if a, b := md.Disc(); a != 0 || b != 0 {
		mdDiscIndex = &a
		mdDiscTotal = &b
	}

	return models.TrackTags{
		Format:      langext.Ptr(md.Format()),
		FileType:    langext.Ptr(md.FileType()),
		Title:       langext.Ptr(md.Title()),
		Album:       langext.Ptr(md.Album()),
		Artist:      langext.Ptr(md.Artist()),
		AlbumArtist: langext.Ptr(md.AlbumArtist()),
		Composer:    langext.Ptr(md.Composer()),
		Year:        langext.Ptr(md.Year()),
		Genre:       langext.Ptr(md.Genre()),
		TrackIndex:  mdTrackIndex,
		TrackTotal:  mdTrackTotal,
		DiscIndex:   mdDiscIndex,
		DiscTotal:   mdDiscTotal,
		Picture:     md.Picture(),
		Lyrics:      estrptr(md.Lyrics()),
		Comment:     estrptr(md.Comment()),
		Raw:         langext.Ptr(md.Raw()),
	}, nil
}

func (db *Database) recalcChecksum(lock bool) {
	if lock {
		db.lock.Lock()
		defer db.lock.Unlock()
	}

	plsts := langext.MapValueArr(db.playlists)
	trcks := make([]models.Track, 0)
	for _, v1 := range db.tracks {
		for _, v2 := range v1 {
			trcks = append(trcks, v2)
		}
	}

	langext.SortBy(plsts, func(v models.Playlist) models.PlaylistID { return v.ID })
	langext.SortBy(trcks, func(v models.Track) models.TrackID { return v.ID })

	str := fmt.Sprintf("%#+v\n%#+v", plsts, trcks)

	db.checksum = strings.ToUpper(cryptext.StrSha256(str))[:16]
}
