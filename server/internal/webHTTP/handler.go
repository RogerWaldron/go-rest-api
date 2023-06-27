package webHTTP

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)


type Handler struct {
	Router 	*mux.Router
	Service CommentService
	Server 	*http.Server
}

func NewHandler(service CommentService) *Handler {
	h := &Handler {
		Service: service,
	}
	h.Router = mux.NewRouter()
	h.mapRoutes()
	h.Router.Use(JSONMiddleware)
	h.Router.Use(LoggingMiddleware)
	h.Router.Use(TimeOutMiddleware)

	h.Server = &http.Server{
		// host.docker.internal didn't work
		Addr: "0.0.0.0:8080",
		WriteTimeout: time.Second * 60,
		ReadTimeout: time.Second * 60,
		IdleTimeout: time.Second * 60,
		Handler: h.Router,
	}

	return h
}

func (h *Handler) mapRoutes() {
	h.Router.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, "Hello World")
	})

	h.Router.HandleFunc("/api/v1/comment", h.PostComment).Methods("POST")
	h.Router.HandleFunc("/api/v1/comment", h.GetComments).Methods("GET")
	h.Router.HandleFunc("/api/v1/comment/{id}", h.GetCommentByID).Methods("GET")
	h.Router.HandleFunc("/api/v1/comment/{id}", h.UpdateComment).Methods("PUT")
	h.Router.HandleFunc("/api/v1/comment/{id}", h.DeleteComment).Methods("DELETE")
}

func (h *Handler) Serve() error {
	log.Info().Msg("Starting HTTP server")
	
	go func() {
		err := h.Server.ListenAndServe()
		if err != nil {
			log.Error().Err(err).Msg("")
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<- c

	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 15)
	defer cancel()
	h.Server.Shutdown(ctx)
	log.Info().Msg("shutting down gracefully HTTP server")

	return nil
}