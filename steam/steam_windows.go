package steam

import (
	"golang.org/x/sys/windows/registry"
	"log"
	"path"
	"steamview-go/appinfo"
)

func init() {
	k, err := registry.OpenKey(registry.CURRENT_USER, `SOFTWARE\Valve\Steam`, registry.QUERY_VALUE)
	if err != nil {
		log.Fatal(err, "! Do you have Steam installed?")
	}

	defer func(k registry.Key) {
		_ = k.Close()
	}(k)

	root, _, err := k.GetStringValue("SteamPath")
	if err != nil {
		log.Fatal(err, "! Do you have Steam installed?")
	}
	CacheRoot = path.Join(root, "appcache")
	imgRoot = path.Join(CacheRoot, "librarycache")
}

func GetAppId() uint32 {
	if appinfo.Reading {
		return 0
	}
	k, err := registry.OpenKey(registry.CURRENT_USER, `SOFTWARE\Valve\Steam`, registry.QUERY_VALUE)
	if err != nil {
		log.Fatal(err, "! Do you have Steam installed?")
	}

	defer func(k registry.Key) {
		_ = k.Close()
	}(k)

	result, _, err := k.GetIntegerValue("RunningAppID")
	if err != nil {
		log.Fatal(err, "! Do you have Steam installed?")
	}
	return uint32(result)
}
