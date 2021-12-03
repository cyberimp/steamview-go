//go:generate go-winres make
package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"steamview-go/steam"
	"steamview-go/trayicon"
	"steamview-go/worker"
	"steamview-go/wsserver"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	setupHandles()

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
	http.HandleFunc("/socket", wsserver.ServeWs)
	http.HandleFunc("/align", worker.SetAlign)
	http.HandleFunc("/cache/", steam.ServeCache)
	http.Handle("/", http.FileServer(http.Dir("./assets")))
}
