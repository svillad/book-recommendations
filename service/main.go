package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/book-recommendations/service/api"
	log "github.com/sirupsen/logrus"
)

func main() {
	if err := run(); err != nil {
		fmt.Println("error :", err)
		os.Exit(1)
	}
}

func run() error {
	router := api.Routes()

	srv := &http.Server{
		Addr:         ":5001",
		WriteTimeout: time.Second * time.Duration(10),
		ReadTimeout:  time.Second * time.Duration(10),
		IdleTimeout:  time.Second * time.Duration(60),
		Handler:      router,
	}

	log.Info("setting up api shutdown")

	serverErr := make(chan error, 1)

	go func() {
		log.Info("starting server")
		serverErr <- srv.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT or SIGKILL
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErr:
		return fmt.Errorf("server error: %w", err)
	case shutdownSignal := <-shutdown:
		log.WithField("shutdown_command", shutdownSignal).Info("starting shutdown")

		// Create a deadline to wait for.
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(60))
		defer cancel()

		// Attempt to gracefully shutdown API
		err := srv.Shutdown(ctx)

		if err != nil {
			log.WithError(err).Error("graceful shutdown failed")
			err = srv.Close()
		}

		// Log the status of this shutdown.
		switch {
		case shutdownSignal == syscall.SIGSTOP:
			return errors.New("integrity issue caused shutdown")
		case err != nil:
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}

		log.Info("server shutdown")
	}

	return nil
}
