package steam

import (
	"golang.org/x/sys/windows/registry"
	"log"
	"net/http"
	"path"
	"regexp"
	"runtime"
)

var ImgRoot string

func init() {
	os := runtime.GOOS
	if os == "windows" {
		k, err := registry.OpenKey(registry.CURRENT_USER, `SOFTWARE\Valve\Steam`, registry.QUERY_VALUE)
		if err != nil {
			log.Fatal(err, "! Do you have Steam installed?")
		}
		defer k.Close()

		root, _, err := k.GetStringValue("SteamPath")
		if err != nil {
			log.Fatal(err, "! Do you have Steam installed?")
		}
		ImgRoot = path.Join(root, "appcache", "librarycache")
	} else {
		ImgRoot = path.Join("~", ".steam", "appcache", "librarycache")
	}
}

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
		fullPath = path.Join(ImgRoot, parse[2]+"_library_hero.jpg")
	} else {
		fullPath = path.Join(ImgRoot, parse[2]+"_logo.png")
	}

	http.ServeFile(w, r, fullPath)
}

func GetAppId() uint64 {
	os := runtime.GOOS
	if os != "windows" {
		//TODO: linux/mac check for running app
		return 0
	}
	k, err := registry.OpenKey(registry.CURRENT_USER, `SOFTWARE\Valve\Steam`, registry.QUERY_VALUE)
	if err != nil {
		log.Fatal(err, "! Do you have Steam installed?")
	}
	defer k.Close()

	result, _, err := k.GetIntegerValue("RunningAppID")
	if err != nil {
		log.Fatal(err, "! Do you have Steam installed?")
	}
	return result
}
