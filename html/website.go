package html

import (
	"embed"
	_ "embed"
)

//go:embed *.html
//go:embed *.js
//go:embed *.css
var Assets embed.FS

func ListAssets() []string {
	result := make([]string, 0)

	entries, err := Assets.ReadDir(".")
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
