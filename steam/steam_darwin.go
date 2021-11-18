package steam

import (
	"os"
	"path"
)

func init() {
	homeDir, _ := os.UserHomeDir()
	imgRoot = path.Join(homeDir, "Library", "Application Support", "Steam", "appcache", "librarycache")
}
