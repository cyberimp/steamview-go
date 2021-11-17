package main

import (
	"log"
	"net/http"
	"steamView/steam"
	"steamView/worker"
	"steamView/wsserver"
)

func main() {
	http.HandleFunc("/socket", wsserver.ServeWs)
	http.HandleFunc("/cache/", steam.ServeCache)
	http.Handle("/", http.FileServer(http.Dir("./assets")))
	go worker.Serve()
	log.Fatal(http.ListenAndServe(":3000", nil))
}
