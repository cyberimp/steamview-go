package steam

import (
	"os"
	"path"
)

func init() {
	homeDir, _ := os.UserHomeDir()
	CacheRoot = path.Join(homeDir, ".steam", "steam", "appcache")
	imgRoot = path.Join(CacheRoot, "librarycache")
}
