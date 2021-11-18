//go:build !windows && !darwin
// +build !windows,!darwin

package steam

import (
	"os"
	"path"
)

func init() {
	homeDir, _ := os.UserHomeDir()
	imgRoot = path.Join(homeDir, ".steam", "steam", "appcache", "librarycache")
}
