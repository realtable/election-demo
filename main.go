package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/realtable/template/handler"
	"github.com/realtable/template/model"
	"github.com/realtable/template/util"
)

func main() {
	// set up slogger
	util.SetDefaultSlogger()

	// load environment variables
	err := godotenv.Load()
	if err != nil {
		slog.Error("could not load .env file", "error", err)
		os.Exit(1)
	}

	// create app
	r := chi.NewRouter()

	// initialise trace provider
	tp, err := util.InitTracerProvider()
	if err != nil {
		slog.Error("could not initialize tracer provider", "error", err)
		os.Exit(1)
	}

	// add middleware
	r.Use(util.TelemetryMiddleware)
	r.Use(middleware.RequestLogger(&util.RequestLogFormatter{}))
	r.Use(middleware.Recoverer)
	// TODO openapi middleware

	// init db
	model.InitDB()

	// declare routes
	r.Handle("/*", http.FileServer(http.Dir("static")))
	r.Route("/api", func(r chi.Router) {
		r.Get("/voters", handler.GetVoters)
		r.Post("/voters", handler.AddVoter)
		r.Delete("/voters", handler.ClearVoters)
		r.Get("/votes", handler.GetVotes)
		r.Post("/votes", handler.AddVote)
		r.Delete("/votes", handler.ClearVotes)
	})

	// set up graceful shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	// start app
	go func() {
		err := http.ListenAndServeTLS(os.Getenv("NETWORK_ADDRESS"), os.Getenv("TLS_CERT_FILE"), os.Getenv("TLS_KEY_FILE"), r)
		if err != nil {
			slog.Error("app server returned an error", "error", err)
		}
		shutdown <- syscall.SIGINT
	}()

	// graceful shutdown
	<-shutdown
	slog.Info("shutting down")
	if err := tp.Shutdown(context.Background()); err != nil {
		slog.Error("error shutting down tracer provider", "error", err)
	}
	util.WaitForBackgroundTasks(1 * time.Second)
	slog.Info("shutdown complete")
}
