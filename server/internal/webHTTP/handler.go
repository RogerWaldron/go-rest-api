package webHTTP

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type CommentService interface {

}

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

	h.Server = &http.Server{
		Addr: "host.docker.internal",
		Handler: h.Router,
	}

	return h
}

func (h *Handler) mapRoutes() {
	h.Router.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, "Hello World")
	})
}

func (h *Handler) Serve() error {
	err := h.Server.ListenAndServe()
	if err != nil {
		return err
	}
	
	return nil
}