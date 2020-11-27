package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"time"
)

func main() {

	pprofRouter := mux.NewRouter()
	AddPprofRoutes(pprofRouter)

	pprofServer := &http.Server{
		Handler:      pprofRouter,
		Addr:         "127.0.0.1:15121",
		WriteTimeout: 100 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	// done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		log.Println("=============================================")
		log.Printf("==== Starting PProf at: %s =====\n", pprofServer.Addr)
		if err := pprofServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on %s: %v", pprofServer.Addr, err)
		}
	}()

	<-quit
}

// AddPprofRoutes
func AddPprofRoutes(r *mux.Router) {
	debugProf := r.PathPrefix("/debug/pprof").Subrouter()
	debugProf.HandleFunc("/", pprof.Index)
	debugProf.HandleFunc("/cmdline", pprof.Cmdline)
	debugProf.HandleFunc("/symbol", pprof.Symbol)
	debugProf.HandleFunc("/trace", pprof.Trace)
	debugProf.HandleFunc("/profile", pprof.Profile)

	// Manually add support for paths not easily linked as above
	// Hooking this up is actually very convoluted and only a few answers on how to do it
	// https://stackoverflow.com/questions/19591065/profiling-go-web-application-built-with-gorillas-mux-with-net-http-pprof
	// Alternatively one can go to localhost:6060
	debugProf.Handle("/goroutine", pprof.Handler("goroutine"))
	debugProf.Handle("/heap", pprof.Handler("heap"))
	debugProf.Handle("/threadcreate", pprof.Handler("threadcreate"))
	debugProf.Handle("/block", pprof.Handler("block"))
}
