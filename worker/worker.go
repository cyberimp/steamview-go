package worker

import (
	"fmt"
	"net/http"
	"steamview-go/aligndb"
	"steamview-go/steam"
	"sync"
	"time"
)

type Message struct {
	Logo  string `json:"logo"`
	Hero  string `json:"hero"`
	Align string `json:"align"`
}

var (
	receivers    map[uint]chan Message
	forceRefresh = false
	panicFlag    = false
	counter      uint
	lock         sync.Mutex
	appID        uint64
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
	align := aligndb.GetAlign(appID)
	if appID > 0 {
		logo = fmt.Sprintf("/cache/logo_%d.png", appID)
		hero = fmt.Sprintf("/cache/hero_%d.jpg", appID)
	}
	return Message{Logo: logo, Hero: hero, Align: align}
}

func SetAlign(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" && appID > 0 {
		_ = r.ParseForm()
		align := r.Form.Get("submit")
		aligndb.SetAlign(appID, align)
		forceRefresh = true
	}
	http.ServeFile(w, r, "./assets/align.html")
}

func Panic() {
	panicFlag = true
}

func Serve() {
	receivers = map[uint]chan Message{}
	for range time.Tick(time.Second / 3) {
		newAppID := steam.GetAppId()
		if panicFlag {
			sendAll(Message{Hero: defaultHero, Align: "absolute-center", Logo: "/images/error.png"})
			break
		}
		if appID != newAppID || forceRefresh {
			appID = newAppID
			sendAll(genMessage())
			forceRefresh = false
		}
	}
}
