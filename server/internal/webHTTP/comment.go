package webHTTP

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/RogerWaldron/go-rest-api/server/internal/comment"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/rs/zerolog/log"
)

const (
	ErrCommentToJSON = "Failed to encode the comment into JSON"
	ErrRequestBodyToJSON = "Failed to decode the request body into JSON"
)

type CommentService interface {
	PostComment(ctx context.Context, cmt comment.Comment) (comment.Comment, error)
	GetComments(ctx context.Context, limit int, offset int) ([]comment.Comment, error)
	GetCommentByID(ctx context.Context, id string) (comment.Comment, error)
	UpdateComment(ctx context.Context, id string, content comment.Comment) (comment.Comment, error)
	DeleteComment(ctx context.Context, id string) error
	// Ping(ctx context.Context) error
}

type queryParms struct {
	Limit 	int `schema:"limit"`
	Offset 	int `schema:"offset"`
}

func (h *Handler) PostComment(w http.ResponseWriter, r *http.Request) {
	var cmt comment.Comment

	if err := json.NewDecoder(r.Body).Decode(&cmt); err != nil {
		log.Error().Err(err).Msg(ErrRequestBodyToJSON)
		return
	}

	cmt, err := h.Service.PostComment(r.Context(), cmt)
	if err != nil {
		log.Error().Err(err).Msg("")
		return
	}

	if err := json.NewEncoder(w).Encode(cmt); err != nil {
		log.Error().Err(err).Msg(ErrCommentToJSON)
		return
	}
}

func (h *Handler) GetComments(w http.ResponseWriter, r *http.Request) {
	var (
		decoder = schema.NewDecoder()
		params queryParms
	)
	err := decoder.Decode(&params, r.URL.Query())
	if err != nil {
		log.Error().Err(err).Msg("error decoding query params")
	}

	cmts, err := h.Service.GetComments(r.Context(), params.Limit, params.Offset)
	if err != nil {
		log.Error().Err(err).Msg("")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(cmts); err != nil {
		log.Error().Err(err).Msg(ErrCommentToJSON)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetCommentByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cmt, err := h.Service.GetCommentByID(r.Context(), id)
	if err != nil {
		log.Error().Err(err).Msg("")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(cmt); err != nil {
		log.Error().Err(err).Msg(ErrCommentToJSON)
	}

}

func (h *Handler) UpdateComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var cmt comment.Comment
	if err:= json.NewDecoder(r.Body).Decode(&cmt); err != nil {
		log.Error().Err(err).Msg(ErrRequestBodyToJSON)
		return 
	}
	cmt, err := h.Service.UpdateComment(r.Context(), id, cmt)
	if err != nil {
		log.Error().Err(err).Msg("")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := h.Service.DeleteComment(r.Context(), id)
	if err != nil {
		log.Error().Err(err).Msg("")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
