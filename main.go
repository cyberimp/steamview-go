//go:generate go-winres make
package main

import (
	"context"
	"embed"
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
	"syscall"
	"time"
)

//go:embed assets
var assets embed.FS

func main() {
	ctx, cancel := context.WithCancel(context.Background())

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
		if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("HTTP server ListenAndServe: %v", err)
		}
	}()

	signalChan := make(chan os.Signal, 1)

	signal.Notify(
		signalChan,
		syscall.SIGHUP,  // kill -SIGHUP XXXX
		syscall.SIGINT,  // kill -SIGINT XXXX or Ctrl+c
		syscall.SIGQUIT, // kill -SIGQUIT XXXX
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

func setupHandles() {
	var staticFS = fs.FS(assets)
	htmlContent, err := fs.Sub(staticFS, "assets")
	if err != nil {
		log.Fatal(err)
	}
	fs := http.FileServer(http.FS(htmlContent))
	http.HandleFunc("/socket", wsserver.ServeWs)
	http.HandleFunc("/cache/", steam.ServeCache)
	http.Handle("/", fs)
}
