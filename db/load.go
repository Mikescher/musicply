package db

import (
	"errors"
	"fmt"
	"github.com/dhowden/tag"
	"github.com/rs/zerolog/log"
	"github.com/titanous/json5"
	"github.com/vansante/go-ffprobe"
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
		err := db.refreshSource(src)
		if err != nil {
			exerr.Wrap(err, "").Fatal()
		}
	}

	fmt.Printf("\n")
	fmt.Printf("================ ================== ================\n")
	fmt.Printf("\n")

}

func (db *Database) refreshSource(src models.Source) error {

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

	files = langext.ArrFilter(files, func(v dataext.Tuple[string, fs.FileInfo]) bool {
		return langext.InArray(filepath.Ext(strings.ToLower(v.V1)), []string{".mp3", ".flac", ".m4a", ".ogg", ".wav", ".wma"})
	})

	tracks := make([]models.Track, 0, len(files))
	for _, fobj := range files {

		fp := fobj.V1
		info := fobj.V2

		track, err := db.analyzeAudioFile(src.ID, fp, info)
		if err != nil {
			log.Error().Msg(fmt.Sprintf("Failed to parse track from file '%s'", fp))
			exerr.Wrap(err, "").Print()
			continue
		}

		tracks = append(tracks, track)
	}

	db.lock.Lock()
	defer db.lock.Unlock()

	fmt.Printf("Found %d tracks in source '%s'\n", len(tracks), src.Name)

	db.tracks[src.ID] = tracks

	return nil
}

func (db *Database) analyzeAudioFile(srcid models.SourceID, fp string, info fs.FileInfo) (models.Track, error) {

	fm, err := db.getFileMeta(fp, info)
	if err != nil {
		return models.Track{}, exerr.Wrap(err, "").Build()
	}

	am, err := db.getAudioMeta(fp)
	if err != nil {
		return models.Track{}, exerr.Wrap(err, "").Build()
	}

	tt, err := db.getTrackTags(fp)
	if err != nil {
		return models.Track{}, exerr.Wrap(err, "").Build()
	}

	return models.Track{
		ID:        models.NewTrackID(),
		SourceID:  srcid,
		FileMeta:  fm,
		AudioMeta: am,
		Tags:      tt,
	}, nil
}

func (db *Database) getFileMeta(fp string, info fs.FileInfo) (models.TrackFileMeta, error) {
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
		Path:      fp,
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
		return models.TrackTags{}, nil
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