package db

import (
	"fmt"
	"github.com/dhowden/tag"
	"gogs.mikescher.com/BlackForestBytes/goext/dataext"
	"gogs.mikescher.com/BlackForestBytes/goext/exerr"
	"gogs.mikescher.com/BlackForestBytes/goext/fsext"
	json "gogs.mikescher.com/BlackForestBytes/goext/gojson"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"io/fs"
	mply "mikescher.com/musicply"
	"mikescher.com/musicply/models"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"syscall"
	"time"
)

type Database struct {
	sources []models.Source

	tracks map[models.SourceID][]models.Track
	lock   sync.Mutex
}

func NewDatabase() *Database {
	return &Database{
		tracks: make(map[models.SourceID][]models.Track),
	}
}

func (db *Database) LoadSourcesFromEnv(envkey string) {

	sourceList := make([]dataext.Tuple[string, string], 0)

	regex := regexp.MustCompile(regexp.QuoteMeta(envkey) + `(_[0-9]+)?`)

	for _, v := range os.Environ() {

		idx := strings.Index(v, "=")
		k := v[0:idx]
		v := v[idx+1:]

		if regex.MatchString(k) {
			sourceList = append(sourceList, dataext.Tuple[string, string]{V1: k, V2: v})
		}
	}

	fmt.Printf("\n")
	fmt.Printf("================ INITIALIZE SOURCES ================\n")
	fmt.Printf("\n")

	for _, v := range sourceList {
		err := db.loadSourceFromSingleEnv(v.V1, v.V2)
		if err != nil {
			exerr.Wrap(err, "failed to load source from env").Str("key", v.V1).Str("value", v.V2).Fatal()
		}
	}

	fmt.Printf("\n")
	fmt.Printf("================ ================== ================\n")
	fmt.Printf("\n")
}

func (db *Database) loadSourceFromSingleEnv(key string, value string) error {

	type sourceSpec struct {
		Name      string `json:"name"`
		Path      string `json:"path"`
		Recursive bool   `json:"recursive"`
	}

	if strings.HasPrefix(value, "{") {

		ss := sourceSpec{}

		err := json.Unmarshal([]byte(value), &ss)
		if err != nil {
			return exerr.Wrap(err, fmt.Sprintf("cannot parse source [%s] as a source-spec", key)).Str("value", value).Build()
		}

		src := models.Source{
			ID:        models.NewSourceID(),
			Name:      ss.Name,
			Path:      ss.Path,
			Recursive: ss.Recursive,
		}

		db.sources = append(db.sources, src)

	} else {

		exist, err := fsext.DirectoryExists(value)
		if err != nil {
			return exerr.Wrap(err, fmt.Sprintf("failed to access directory '%s' (from [%s])", value, key)).Str("path", value).Build()
		}
		if !exist {
			return exerr.New(mply.ErrSourceNotFound, fmt.Sprintf("directory '%s' does not exist (from [%s])", value, key)).Str("path", value).Build()
		}

		basename := path.Base(value)

		name := basename

		for i := 1; langext.ArrAny(db.sources, func(src models.Source) bool { return src.Name == name }); i++ {
			name = fmt.Sprintf("%s (%d)", basename, i)
		}

		src := models.Source{
			ID:        models.NewSourceID(),
			Name:      name,
			Path:      value,
			Recursive: false,
		}

		db.sources = append(db.sources, src)

		fmt.Printf("[%s] Initialize '%s' from source: '%s'\n", key, src.Name, src.Path)

	}

	return nil
}

func (db *Database) RefreshAllInitial() {

	fmt.Printf("\n")
	fmt.Printf("================ INITIALIZE SOURCES ================\n")
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
	for _, f := range files {

		var ctime *time.Time = nil
		var atime *time.Time = nil

		if v, ok := f.V2.Sys().(*syscall.Stat_t); ok {
			ctime = langext.Ptr(time.Unix(v.Ctim.Sec, v.Ctim.Nsec))
			atime = langext.Ptr(time.Unix(v.Atim.Sec, v.Atim.Nsec))
		}

		fptr, err := os.OpenFile(f.V1, os.O_RDONLY, 0755)
		if err != nil {
			return err
		}

		md, err := tag.ReadFrom(fptr)
		_ = fptr.Close()
		if err != nil {
			return err
		}

		mdTrackIndex, mdTrackTotal := md.Track()
		mdDiscIndex, mdDiscTotal := md.Disc()

		tracks = append(tracks, models.Track{
			ID: models.NewTrackID(),
			FileMeta: models.TrackFileMeta{
				Path:      f.V1,
				Filename:  filepath.Base(f.V1),
				Extension: strings.TrimLeft(filepath.Ext(strings.ToLower(filepath.Base(f.V1))), "."),
				Size:      f.V2.Size(),
				Filemode:  f.V2.Mode(),
				ModTime:   f.V2.ModTime(),
				CTime:     ctime,
				ATime:     atime,
			},
			Tags: models.TrackTags{
				Format:      md.Format(),
				FileType:    md.FileType(),
				Title:       md.Title(),
				Album:       md.Album(),
				Artist:      md.Artist(),
				AlbumArtist: md.AlbumArtist(),
				Composer:    md.Composer(),
				Year:        md.Year(),
				Genre:       md.Genre(),
				TrackIndex:  mdTrackIndex,
				TrackTotal:  mdTrackTotal,
				DiscIndex:   mdDiscIndex,
				DiscTotal:   mdDiscTotal,
				Picture:     md.Picture(),
				Lyrics:      md.Lyrics(),
				Comment:     md.Comment(),
				Raw:         md.Raw(),
			},
		})
	}

	db.lock.Lock()
	defer db.lock.Unlock()

	fmt.Printf("Found %d tracks in source '%s'\n", len(tracks), src.Name)

	db.tracks[src.ID] = tracks

	return nil
}
