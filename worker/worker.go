package worker

import (
	"fmt"
	"steamview-go/appinfo"
	"steamview-go/steam"
	"sync"
	"time"
)

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
)

const (
	defaultHero = "/images/no_hero.png"
	defaultLogo = "/images/default.png"
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
	logo := defaultLogo
	hero := defaultHero
	align := "CenterCenter"
	width := "50"
	height := "50"
	name := "Steam"
	if appID > 0 {
		logo = fmt.Sprintf("/cache/logo_%d.png", appID)
		hero = fmt.Sprintf("/cache/hero_%d.jpg", appID)
		info := appinfo.GetAppInfo(appID)
		align = info.GetAlign()
		name = info.GetName()
		if align == "" {
			align = "hidden"
		} else {
			width = info.GetWidth()
			height = info.GetHeight()
		}
	}
	return Message{Logo: logo, Hero: hero, Align: align, Width: width, Height: height, Name: name}
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
			if trySend() {
				return
			}
		case <-blocker:
			if trySend() {
				return
			}
		}
	}
}

func trySend() bool {
	if panicFlag {
		sendAll(Message{Hero: defaultHero, Align: "CenterCenter", Logo: "/images/error.png"})
		return true
	}
	newAppID := steam.GetAppId()
	if appID != newAppID || forceRefresh {
		appID = newAppID
		sendAll(genMessage())
		forceRefresh = false
	}
	return false
}
