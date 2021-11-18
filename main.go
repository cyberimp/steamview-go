package main

import (
	"log"
	"net/http"
	"steamview-go/steam"
	"steamview-go/worker"
	"steamview-go/wsserver"
)

func main() {
	http.HandleFunc("/socket", wsserver.ServeWs)
	http.HandleFunc("/align", worker.SetAlign)
	http.HandleFunc("/cache/", steam.ServeCache)
	http.Handle("/", http.FileServer(http.Dir("./assets")))
	go worker.Serve()
	log.Fatal(http.ListenAndServe(":3000", nil))
}
