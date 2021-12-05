package trayicon

import (
	"github.com/getlantern/systray"
	"github.com/pkg/browser"
	"steamview-go/icon"
)

var quitChan chan struct{}

func Run(quit chan struct{}) {
	quitChan = quit
	go systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetIcon(icon.GetIcon())
	systray.SetTitle("Awesome App")
	mOpen := systray.AddMenuItem("Open &browser", "Open app in browser")
	mAlign := systray.AddMenuItem("Set &align", "Set align for current banner")
	mQuitOrig := systray.AddMenuItem("&Quit", "Quit the whole app")

	go func() {
		for {
			select {
			case <-mQuitOrig.ClickedCh:
				quitChan <- struct{}{}
				systray.Quit()
				return
			case <-mOpen.ClickedCh:
				_ = browser.OpenURL("http://127.0.0.1:3000")
			case <-mAlign.ClickedCh:
				_ = browser.OpenURL("http://127.0.0.1:3000/align")
			}
		}
	}()
}

func onExit() {

}

func Quit() {
	systray.Quit()
}
