package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// WaitTimeout is the duration for which the server should gracefully wait for existing connections to finish.
// Measured in seconds - for example 15 (15 seconds) or 60 (1 minute).
var WaitTimeout time.Duration

func main() {

	port := os.Getenv("PORT")

	to, err := strconv.Atoi(os.Getenv("WAIT_TIMEOUT_SECONDS"))
	if err != nil {
		log.Println("Invalid WAIT_TIMEOUT_SECONDS, setting to 15 seconds")
		to = 15
	}
	WaitTimeout = time.Duration(to) * time.Second

	router := mux.NewRouter()
	router.HandleFunc("/timer", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Minute)
	})

	srv := &http.Server{
		Addr: "0.0.0.0:" + port,
		// Good practice to set timeouts to avoid Slowloris attacks.
		// These are not 12-Factored as they should not be changed.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router,
	}

	// Run the server in a goroutine so that it doesn't block.
	go func() {
		log.Println("Listening on http://localhost:" + port + " ...")
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// Accept graceful shutdowns when quit via SIGINT (Ctrl+C).
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) are not caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c
	log.Println("Shutdown request (Ctrl-C) caught")

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), WaitTimeout)
	defer cancel()
	// Don't block if no connections, otherwise wait until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if application should wait for other services
	// to finalize based on context cancellation.
	log.Println("Shutting down ...")
	os.Exit(0)
}
