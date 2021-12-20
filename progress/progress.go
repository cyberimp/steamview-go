package progress

import (
	"github.com/ncruces/zenity"
	"log"
	"os"
)

var (
	full   int64
	dialog zenity.ProgressDialog
)

func SetLength(file *os.File) {
	info, err := file.Stat()
	if err != nil {
		panic("cannot read size of file:" + err.Error())
	}
	full = info.Size()
}

func Display() {
	var err error
	dialog, err = zenity.Progress(zenity.NoCancel(),
		zenity.Icon(zenity.NoIcon),
		zenity.Title("reading appinfo.vdf..."),
		zenity.MaxValue(100))

	if err != nil {
		log.Print("You don't have zenity installed:", err)
	}

	_ = dialog.Value(0)
}

func SetValue(value int64) {
	_ = dialog.Value(int(float32(value) / float32(full) * 100))

}

func Close() {
	_ = dialog.Close()
}
