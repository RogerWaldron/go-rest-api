package webHTTP

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/RogerWaldron/go-rest-api/server/internal/comment"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/rs/zerolog/log"
)

const (
	ErrCommentToJSON = "Failed to encode the comment into JSON"
	ErrRequestBodyToJSON = "Failed to decode the request body into JSON"
	ErrCommentFailedValidation = "Failed to validate comment content"
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

type PostCommentRequest struct {
	Slug   string `json:"slug" validate:"required"`
	Author string `json:"author" validate:"required"`
	Body   string `json:"body" validate:"required"`
}

func commentFromPostCommentRequest(u PostCommentRequest) comment.Comment {
	return comment.Comment{
		Slug:   u.Slug,
		Author: u.Author,
		Body:   u.Body,
	}
}

func (h *Handler) PostComment(w http.ResponseWriter, r *http.Request) {
	var postCmt PostCommentRequest

	if err := json.NewDecoder(r.Body).Decode(&postCmt); err != nil {
		log.Error().Err(err).Msg(ErrRequestBodyToJSON)
		return
	}

	validate := validator.New()
	err := validate.Struct(postCmt)
	if err != nil {
		log.Info().Err(err).Msg(ErrCommentFailedValidation)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cmt, err := h.Service.PostComment(r.Context(), commentFromPostCommentRequest(postCmt))
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

type UpdateCommentRequest struct {
	Slug   string `json:"slug" validate:"required"`
	Author string `json:"author" validate:"required"`
	Body   string `json:"body" validate:"required"`
}

func commentFromUpdateCommentRequest(u UpdateCommentRequest) comment.Comment {
	return comment.Comment{
		Slug:   u.Slug,
		Author: u.Author,
		Body:   u.Body,
	}
}

func (h *Handler) UpdateComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var updateCmtRequest UpdateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&updateCmtRequest); err != nil {
		return
	}

	validate := validator.New()
	err := validate.Struct(updateCmtRequest)
	if err != nil {
		log.Error().Err(err).Msg(ErrCommentFailedValidation)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err:= json.NewDecoder(r.Body).Decode(&updateCmtRequest); err != nil {
		log.Error().Err(err).Msg(ErrRequestBodyToJSON)
		return 
	}
	cmt, err := h.Service.UpdateComment(r.Context(), id, commentFromUpdateCommentRequest(updateCmtRequest))
	if err != nil {
		log.Error().Err(err).Msg("")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(cmt); err != nil {
		panic(err)
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
