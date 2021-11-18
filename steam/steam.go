package steam

import (
	"net/http"
	"path"
	"regexp"
)

var imgRoot string

func ServeCache(w http.ResponseWriter, r *http.Request) {

	d, f := path.Split(r.URL.Path)
	if d != "/cache/" {
		http.NotFound(w, r)
		return
	}

	//[1] - type, [2] - appid, [3] - extension
	split := regexp.MustCompile(`(hero|logo)_(\d+)\.(jpg|png)`)
	parse := split.FindStringSubmatch(f)
	if parse == nil {
		http.NotFound(w, r)
		return
	}

	var fullPath string

	if parse[1] == "hero" {
		fullPath = path.Join(imgRoot, parse[2]+"_library_hero.jpg")
	} else {
		fullPath = path.Join(imgRoot, parse[2]+"_logo.png")
	}

	http.ServeFile(w, r, fullPath)
}
