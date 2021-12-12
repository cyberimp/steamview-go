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
	systray.SetTitle("SteamView")
	mOpen := systray.AddMenuItem("Open browser", "Open app in browser")
	mQuitOrig := systray.AddMenuItem("Quit", "Quit the whole app")

	_ = browser.OpenURL("http://127.0.0.1:3000")
	go func() {
		for {
			select {
			case <-mQuitOrig.ClickedCh:
				quitChan <- struct{}{}
				systray.Quit()
				return
			case <-mOpen.ClickedCh:
				_ = browser.OpenURL("http://127.0.0.1:3000")
			}
		}
	}()
}

func onExit() {

}

func Quit() {
	systray.Quit()
}
