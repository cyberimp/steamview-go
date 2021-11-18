//go:build !windows
// +build !windows

package steam

import (
	ps "github.com/shirou/gopsutil/process"
	"regexp"
	"strconv"
)

func GetAppId() uint64 {
	var result uint64

	split := regexp.MustCompile(`SteamLaunch AppId=(\d+) `)

	arr, _ := ps.Processes()
	for _, process := range arr {
		str, err := process.Cmdline()
		if err != nil {
			continue
		}
		parse := split.FindStringSubmatch(str)
		if parse != nil {
			result, _ = strconv.ParseUint(parse[1], 10, 64)
			return result
		}
	}
	return result
}
