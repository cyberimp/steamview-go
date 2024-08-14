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

// assets folder packed into binary
//
//go:embed assets
var assets embed.FS

// sumsRaw contains precalculated md5 sums in raw format
//
//go:embed md5sums.txt
var sumsRaw string

// sumsMap contains mapping filename -> md5
var sumsMap map[string]string

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	parseSums()
	setupHandles()

	//start async parser of Steam DB
	go appinfo.ParseAsync(path.Join(steam.CacheRoot, "appinfo.vdf"))

	//start worker watching changes in current working appID
	go worker.Serve()

	httpServer := &http.Server{
		Addr:        ":3000",
		Handler:     nil,
		BaseContext: func(_ net.Listener) context.Context { return ctx },
	}

	//stop worker on server shutdown, sending "Server is dead" to all clients
	httpServer.RegisterOnShutdown(func() {
		worker.Panic()
	})

	//start http server
	go func() {
		if err := httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server ListenAndServe: %v", err)
		}
	}()

	//setup channel for kill command
	signalChan := make(chan os.Signal, 1)
	signal.Notify(
		signalChan,
		syscall.SIGHUP,  // kill -SIGHUP processID
		syscall.SIGINT,  // kill -SIGINT processID or Ctrl+c
		syscall.SIGQUIT, // kill -SIGQUIT processID
	)

	//run tray icon in separate goroutine
	quitChan := make(chan struct{})
	trayicon.Run(quitChan)
	defer trayicon.Quit()

	//wait for shutdown signal from tray icon or terminal
	select {
	case <-signalChan:
	case <-quitChan:
	}

	//now shutdown server gracefully
	gracefulCtx, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()
	err := httpServer.Shutdown(gracefulCtx)
	if err != nil {
		log.Fatalln("shutdown error: ", err)
	}

	cancel()
}

// parseSums parses sumsRaw into sumsMap
func parseSums() {
	sumsSlice := strings.Fields(sumsRaw)
	sumsMap = map[string]string{}
	for i := 0; i < len(sumsSlice); i += 2 {
		//serving "/" is equal to serving embed.FS/index.html
		if sumsSlice[i+1] == "index.html" {
			sumsMap["/"] = sumsSlice[i]
		}
		//works both in linux and windows
		sumsMap["/"+strings.ReplaceAll(sumsSlice[i+1], "\\", "/")] = sumsSlice[i]
	}
}

// serveFSCached is middleware for sending embedded resources with their md5sum as e-tag
func serveFSCached(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=2592000")
		md5Sum := sumsMap[r.URL.Path]
		w.Header().Set("ETag", md5Sum)
		next.ServeHTTP(w, r)
	})
}

// setupHandles sets up server handles
func setupHandles() {
	var staticFS = fs.FS(assets)
	htmlContent, err := fs.Sub(staticFS, "assets")
	if err != nil {
		log.Fatal(err)
	}
	ServeFs := http.FileServer(http.FS(htmlContent))
	//handle for socketserver
	http.HandleFunc("/socket", wsserver.ServeWs)
	//handle for steam cache, app will get pictures from real steam cache
	http.HandleFunc("/cache/", steam.ServeCache)
	//any other files would come from assets folder
	http.Handle("/", serveFSCached(ServeFs))
}
