//go:generate go-winres make
package main

import (
	"context"
	"embed"
	"errors"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path"
	"steamview-go/appinfo"
	"steamview-go/steam"
	"steamview-go/trayicon"
	"steamview-go/worker"
	"steamview-go/wsserver"
	"strings"
	"syscall"
	"time"
)

//go:embed assets
var assets embed.FS

//go:embed md5sums.txt
var sums string

var sumsMap map[string]string

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	parseSums()
	setupHandles()

	go appinfo.ParseAsync(path.Join(steam.CacheRoot, "appinfo.vdf"))

	go worker.Serve()

	httpServer := &http.Server{
		Addr:        ":3000",
		Handler:     nil,
		BaseContext: func(_ net.Listener) context.Context { return ctx },
	}

	httpServer.RegisterOnShutdown(func() {
		worker.Panic()
	})

	go func() {
		if err := httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server ListenAndServe: %v", err)
		}
	}()

	signalChan := make(chan os.Signal, 1)

	signal.Notify(
		signalChan,
		syscall.SIGHUP,  // kill -SIGHUP processID
		syscall.SIGINT,  // kill -SIGINT processID or Ctrl+c
		syscall.SIGQUIT, // kill -SIGQUIT processID
	)

	quitChan := make(chan struct{})
	trayicon.Run(quitChan)
	defer trayicon.Quit()

	select {
	case <-signalChan:
	case <-quitChan:
	}

	gracefulCtx, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()
	err := httpServer.Shutdown(gracefulCtx)
	if err != nil {
		log.Fatalln("shutdown error: ", err)
	}

	cancel()
}

func parseSums() {
	sumsSlice := strings.Fields(sums)
	sumsMap = map[string]string{}
	for i := 0; i < len(sumsSlice); i += 2 {
		if sumsSlice[i+1] == "index.html" {
			sumsMap["/"] = sumsSlice[i]
		}
		sumsMap["/"+strings.ReplaceAll(sumsSlice[i+1], "\\", "/")] = sumsSlice[i]
	}
}

func ServeFSCached(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=2592000")
		md5Sum := sumsMap[r.URL.Path]
		w.Header().Set("ETag", md5Sum)
		next.ServeHTTP(w, r)
	})
}

func setupHandles() {
	var staticFS = fs.FS(assets)
	htmlContent, err := fs.Sub(staticFS, "assets")
	if err != nil {
		log.Fatal(err)
	}
	ServeFs := http.FileServer(http.FS(htmlContent))
	http.HandleFunc("/socket", wsserver.ServeWs)
	http.HandleFunc("/cache/", steam.ServeCache)
	http.Handle("/", ServeFSCached(ServeFs))
}
