//go:build !windows
// +build !windows

package steam

import (
	ps "github.com/shirou/gopsutil/process"
	"regexp"
	"steamview-go/appinfo"
	"strconv"
)

func GetAppId() uint32 {
	var result uint64

	if appinfo.Reading {
		return 0
	}

	split := regexp.MustCompile(`/gameoverlayui.*-gameid (\d+)$`)

	arr, _ := ps.Processes()
	for _, process := range arr {
		str, err := process.Cmdline()
		if err != nil {
			continue
		}
		parse := split.FindStringSubmatch(str)
		if parse != nil {
			result, _ = strconv.ParseUint(parse[1], 10, 32)
			return uint32(result)
		}
	}
	return uint32(result)
}
