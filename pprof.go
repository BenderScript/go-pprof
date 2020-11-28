package main

import (
	"context"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"net/http/pprof"
	"time"
)

func InitPprof(termPprof, pprofDone chan bool) {

	pprofRouter := mux.NewRouter()
	AddPprofRoutes(pprofRouter)

	pprofServer := &http.Server{
		Handler:      pprofRouter,
		Addr:         "0.0.0.0:15121",
		WriteTimeout: 100 * time.Second,
		ReadTimeout:  100 * time.Second,
	}

	go func() {
		log.Printf("==== Starting PProf at: %s =====\n", pprofServer.Addr)
		if err := pprofServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on %s: %v", pprofServer.Addr, err)
		}
	}()

	go func() {
		<-termPprof
		log.Println("Pprof Server is shutting down...")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		pprofServer.SetKeepAlivesEnabled(false)
		if err := pprofServer.Shutdown(ctx); err != nil {
			log.Fatalf("Could not gracefully shutdown the server: %v", err)
		}
		pprofDone <- true
		log.Println("Pprof Server stopped")
	}()
}

// AddPprofRoutes
func AddPprofRoutes(r *mux.Router) {
	debugProf := r.PathPrefix("/debug/pprof").Subrouter()
	debugProf.HandleFunc("/check", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte("Pprof!"))
	})
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
	debugProf.Handle("/vars", http.DefaultServeMux)
}
