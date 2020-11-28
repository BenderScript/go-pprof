package main

import (
	"context"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	// These channels are not really needed in this entire example but make things nicer
	termPprof := make(chan bool)
	pprofDone := make(chan bool)
	done := make(chan bool)

	InitPprof(termPprof, pprofDone)

	newRouter := mux.NewRouter()

	newRouter.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte("Go Away!"))
	})

	server := &http.Server{
		Handler:      newRouter,
		Addr:         "0.0.0.0:15120",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("==== Starting Server at: %s =====\n", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on %s: %v", server.Addr, err)
		}
	}()

	go func() {
		// A interrupt signal
		<-quit
		log.Println("Server is shutting down...")
		// We signal pprof server to stop
		termPprof <- true
		// prof server tells us it is done
		<-pprofDone

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("Could not gracefully shutdown the server: %v", err)
		}
		// We are done
		close(done)
	}()

	<-done
	log.Println("Server stopped")

}
