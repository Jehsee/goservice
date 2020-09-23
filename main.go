package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Jehsee/goservice/handlers"
)

func main() {
	l := log.New(os.Stdout, "product-api", log.LstdFlags)

	// create the handlers
	ph := handlers.NewProducts(l)
	gh := handlers.NewGoodbye(l)

	// create a new serve mux and register the handlers
	sm := http.NewServeMux()
	sm.Handle("/", ph)
	sm.Handle("/goodbye", gh)

	// create a new server
	s := &http.Server{
		Addr:         ":8000",
		Handler:      sm,
		ErrorLog:     l,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// start the server
	go func() {
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	// trap sigterm or interup and gracefully shutdown the server
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	l.Println("Recieved terminate, graceful shutdown", sig)

	// timeout Context
	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)
}
