package wsserver

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"steamView/worker"
)

type Client struct {
	soc *websocket.Conn
	ch  chan worker.Message
	h   uint
}

func (c Client) Loop() {
	defer func() {
		worker.Remove(c.h)
		_ = c.soc.Close()
	}()

	for i := range c.ch {
		err := c.soc.WriteJSON(i)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func ServeWs(w http.ResponseWriter, r *http.Request) {
	upgrade := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	ws, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		if _, ok := err.(websocket.HandshakeError); !ok {
			log.Println(err)
		}
		return
	}
	cl := Client{soc: ws}
	cl.ch, cl.h = worker.GetChan()
	go cl.Loop()
}
