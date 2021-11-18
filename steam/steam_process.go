//go:build !windows
// +build !windows

package steam

import (
	"github.com/mitchellh/go-ps"
	"regexp"
	"strconv"
)

func GetAppId() uint64 {
	var result uint64

	split := regexp.MustCompile(`SteamLaunch AppId=(\d+) `)

	arr, _ := ps.Processes()
	for _, process := range arr {
		str := process.Executable()
		parse := split.FindStringSubmatch(str)
		if parse != nil {
			result, _ = strconv.ParseUint(parse[1], 10, 64)
			return result
		}
	}
	return result
}
