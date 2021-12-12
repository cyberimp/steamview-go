package steam

import (
	"os"
	"path"
)

func init() {
	homeDir, _ := os.UserHomeDir()
	CacheRoot = path.Join(homeDir, "Library", "Application Support", "Steam", "appcache")
	imgRoot = path.Join(CacheRoot, "librarycache")
}
