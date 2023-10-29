package webassets

import (
	"embed"
	_ "embed"
	"fmt"
	"github.com/rs/zerolog/log"
	"gogs.mikescher.com/BlackForestBytes/goext/exerr"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"gogs.mikescher.com/BlackForestBytes/goext/rext"
	"io"
	mply "mikescher.com/musicply"
	"mikescher.com/musicply/models"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

//go:embed *.html
//go:embed *.js
//go:embed *.css
var _assets embed.FS

type templateCacheEntry struct {
	MDate time.Time
	Value ITemplate
}

type fileCacheEntry struct {
	MDate time.Time
	Value []byte
}

type Footerlink struct {
	ID       models.FooterLinkID
	IconPath string
	IconData []byte
	Name     string
	Link     string
}

type Assets struct {
	templateCache map[string]templateCacheEntry
	fileCache     map[string]fileCacheEntry
	footerlinks   []Footerlink
	lock          sync.RWMutex
}

func NewAssets() *Assets {
	return &Assets{
		templateCache: make(map[string]templateCacheEntry, 128),
		fileCache:     make(map[string]fileCacheEntry, 128),
		footerlinks:   make([]Footerlink, 0),
		lock:          sync.RWMutex{},
	}
}

type ITemplate interface {
	Execute(wr io.Writer, data any) error
}

func (a *Assets) ListAssets() []string {
	result := make([]string, 0)

	entries, err := _assets.ReadDir(".")
	if err != nil {
		panic(err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			panic("TODO implement recursion")
		}

		if entry.Name() == "index.html" {
			continue
		}

		result = append(result, "/"+entry.Name())
	}

	return result
}

func (a *Assets) Read(fp string) ([]byte, error) {
	if mply.Conf.LiveReload == nil {

		// no live-reload: use embedded data

		bin, err := _assets.ReadFile(fp)
		if err != nil {
			return nil, err
		}
		return bin, nil

	} else {

		liveFP := filepath.Join(*mply.Conf.LiveReload, fp)

		stat, err := os.Stat(liveFP)
		if err != nil {
			return nil, err
		}

		a.lock.RLock()
		v, ok := a.fileCache[fp]
		a.lock.RUnlock()

		if !ok {

			// initial load

			bin, err := os.ReadFile(liveFP)
			if err != nil {
				return nil, err
			}

			a.lock.Lock()
			a.fileCache[fp] = fileCacheEntry{MDate: stat.ModTime(), Value: bin}
			a.lock.Unlock()

			return bin, nil

		} else if v.MDate != stat.ModTime() {

			// live reload

			log.Info().Msg(fmt.Sprintf("[>>] Live reload file '%s' from filesystem (file changed)", fp))

			bin, err := os.ReadFile(liveFP)
			if err != nil {
				return nil, err
			}

			a.lock.Lock()
			a.fileCache[fp] = fileCacheEntry{MDate: stat.ModTime(), Value: bin}
			a.lock.Unlock()

			return bin, nil

		} else {
			// return from cache
			return v.Value, nil
		}
	}
}

func (a *Assets) Template(fp string, builder func([]byte) (ITemplate, error)) (ITemplate, error) {
	if mply.Conf.LiveReload == nil {

		// no live-reload: use embedded data, and permanently cache template

		a.lock.RLock()
		v, ok := a.templateCache[fp]
		a.lock.RUnlock()
		if ok {
			return v.Value, nil
		}

		bin, err := _assets.ReadFile(fp)
		if err != nil {
			return nil, err
		}
		t, err := builder(bin)
		if err != nil {
			panic(err)
		}

		a.lock.Lock()
		a.templateCache[fp] = templateCacheEntry{MDate: time.Now(), Value: t}
		a.lock.Unlock()

		return t, nil

	} else {

		a.lock.RLock()
		v, ok := a.templateCache[fp]
		a.lock.RUnlock()

		liveFP := filepath.Join(*mply.Conf.LiveReload, fp)

		stat, err := os.Stat(liveFP)
		if err != nil {
			return nil, err
		}

		if !ok {

			// initial load

			bin, err := os.ReadFile(liveFP)
			if err != nil {
				return nil, err
			}
			t, err := builder(bin)
			if err != nil {
				return nil, err
			}

			a.lock.Lock()
			a.templateCache[fp] = templateCacheEntry{MDate: stat.ModTime(), Value: t}
			a.lock.Unlock()

			return t, nil

		} else if v.MDate != stat.ModTime() {

			// live reload

			log.Info().Msg(fmt.Sprintf("[>>] Live reload file '%s' from filesystem (file changed)", fp))

			bin, err := os.ReadFile(liveFP)
			if err != nil {
				return nil, err
			}
			t, err := builder(bin)
			if err != nil {
				return nil, err
			}

			a.lock.Lock()
			a.templateCache[fp] = templateCacheEntry{MDate: stat.ModTime(), Value: t}
			a.lock.Unlock()

			return t, nil

		} else {
			// return from cache
			return v.Value, nil
		}
	}
}

func (a *Assets) LoadDynamicAssets() {

	regexFooterLink := rext.W(regexp.MustCompile("^FOOTERLINK(_[A-Z0-9]+)?$"))

	envs := os.Environ()
	langext.Sort(envs)

	for _, env := range envs {
		idx := strings.Index(env, "=")
		key := env[:idx]
		val := env[idx+1:]
		if regexFooterLink.IsMatch(key) {
			split := strings.Split(val, ";")
			if len(split) != 3 {
				exerr.New(mply.ErrEnviron, "failed to parse environment variable: "+key).Str("key", key).Str("val", val).Fatal()
			}

			data, err := os.ReadFile(split[0])
			if err != nil {
				exerr.Wrap(err, "failed to read icon-file from environment variable: "+key).Str("key", key).Str("val", val).Fatal()
			}

			a.footerlinks = append(a.footerlinks, Footerlink{
				ID:       models.NewFooterLinkID(),
				IconPath: split[0],
				IconData: data,
				Name:     split[1],
				Link:     split[2],
			})
		}
	}
}

func (a *Assets) GetFooterLink(id models.FooterLinkID) *Footerlink {
	for _, v := range a.footerlinks {
		if v.ID == id {
			return &v
		}
	}
	return nil
}

func (a *Assets) ListFooterLinks() []Footerlink {
	return a.footerlinks
}

func (a *Assets) NoCover() []byte {
	bin, err := _assets.ReadFile("no_cover.png")
	if err != nil {
		panic(err)
	}
	return bin
}
