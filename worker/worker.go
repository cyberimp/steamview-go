package worker

import (
	"fmt"
	"steamview-go/appinfo"
	"steamview-go/steam"
	"sync"
	"time"
)

//Message is sent through websocket to clients
type Message struct {
	Logo   string `json:"logo"`
	Hero   string `json:"hero"`
	Align  string `json:"align"`
	Width  string `json:"width"`
	Height string `json:"height"`
	Name   string `json:"name"`
}

var (
	receivers    map[uint]chan Message
	forceRefresh = false
	panicFlag    = false
	counter      uint
	lock         sync.Mutex
	appID        uint32
	blocker      chan int
	ticker       *time.Ticker
	oldReading   = true
)

const (
	//defaultHero is transparent background
	defaultHero = "/images/no_hero.png"
	//defaultLogo is Steam logo
	defaultLogo = "/images/default.png"
	//defaultAlign is centered logo position
	defaultAlign = "CenterCenter"
	//errorLogo is sent on server shutdown
	errorLogo = "/images/error.png"
)

func GetChan() (chan Message, uint) {
	lock.Lock()
	defer lock.Unlock()
	result := make(chan Message)
	receivers[counter] = result
	counter++
	blocker <- 1
	forceRefresh = true
	return result, counter - 1
}

func Remove(i uint) {
	lock.Lock()
	defer lock.Unlock()
	close(receivers[i])
	delete(receivers, i)
}

func sendAll(m Message) {
	lock.Lock()
	defer lock.Unlock()
	for _, messages := range receivers {
		messages <- m
	}
}

func genMessage() Message {
	result := Message{
		Logo:   defaultLogo,
		Hero:   defaultHero,
		Align:  defaultAlign,
		Width:  "50",
		Height: "50",
		Name:   "Steam",
	}

	if appinfo.Reading {
		result.Width = fmt.Sprintf("%f", appinfo.GetProgress())
		result.Name = "_VDF_READING"
		return result
	}

	if appID > 0 {
		result.Logo = fmt.Sprintf("/cache/logo_%d.png", appID)
		result.Hero = fmt.Sprintf("/cache/hero_%d.jpg", appID)
		info := appinfo.GetAppInfo(appID)
		result.Align = info.GetAlign()
		result.Name = info.GetName()
		if result.Align == "" {
			result.Align = "hidden"
		} else {
			result.Width = info.GetWidth()
			result.Height = info.GetHeight()
		}
	}

	return result
}

func Panic() {
	panicFlag = true
	blocker <- 1
}

func Serve() {
	receivers = map[uint]chan Message{}

	blocker = make(chan int)
	ticker = time.NewTicker(time.Second / 3)

	for {
		select {
		case <-ticker.C:
		case <-blocker:
		}
		if trySend() {
			return
		}
	}
}

func trySend() bool {
	if panicFlag {
		sendAll(Message{
			Hero:   defaultHero,
			Align:  defaultAlign,
			Logo:   errorLogo,
			Width:  "50",
			Height: "50",
			Name:   "Error",
		})
		return true
	}
	newAppID := steam.GetAppId()
	if appID != newAppID || forceRefresh || oldReading {
		appID = newAppID
		sendAll(genMessage())
		forceRefresh = false
		oldReading = appinfo.Reading
	}
	return false
}
