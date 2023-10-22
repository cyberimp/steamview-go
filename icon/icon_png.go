//go:build !windows
// +build !windows

package icon

import _ "embed"

//go:embed icons/favicon-32x32.png
var icon byte[]
