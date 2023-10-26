package webassets

import (
	"embed"
	_ "embed"
	"fmt"
	"github.com/rs/zerolog/log"
	"html/template"
	mply "mikescher.com/musicply"
	"os"
	"path/filepath"
	"sync"
	"time"
)

//go:embed *.html
//go:embed *.js
//go:embed *.css
var _assets embed.FS

type templateCacheEntry struct {
	MDate time.Time
	Value *template.Template
}

type fileCacheEntry struct {
	MDate time.Time
	Value []byte
}

type Assets struct {
	templateCache map[string]templateCacheEntry
	fileCache     map[string]fileCacheEntry
	lock          sync.RWMutex
}

func NewAssets() *Assets {
	return &Assets{
		templateCache: make(map[string]templateCacheEntry, 128),
		fileCache:     make(map[string]fileCacheEntry, 128),
		lock:          sync.RWMutex{},
	}
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

		a.lock.RLock()
		v, ok := a.fileCache[fp]
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

func (a *Assets) Template(fp string, builder func([]byte) (*template.Template, error)) (*template.Template, error) {
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
