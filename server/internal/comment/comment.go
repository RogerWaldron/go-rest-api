package comment

import (
	"context"
	"errors"
	"fmt"
)

var (
	ErrFetchingCommentByID = errors.New("failed to fetch comment by ID")
	ErrNotImplemented = errors.New("not implemented")
)

type Store interface {
	CreateComment(context.Context, Comment) (Comment, error)
	GetComments(context.Context) ([]Comment, error)
	GetCommentByID(context.Context, string) (Comment, error)
	UpdateComment(context.Context, Comment) error
	DeleteComment(context.Context, string) error
}

type Comment struct {
	ID string
	Slug string
	Body string
	Author string
}

type Service struct{
	Store Store
}

func NewService(store Store) *Service {
	return &Service{
		Store: store,
	}
}

func (s *Service) CreateComment(ctx context.Context, content Comment) (Comment, error) {
	return Comment{}, ErrNotImplemented
}

func (S *Service) GetComments(ctx context.Context) ([]Comment, error) {
	return []Comment{}, ErrNotImplemented
}

func (s *Service) GetCommentByID(ctx context.Context, id string) (Comment, error) {
	comment, err := s.Store.GetCommentByID(ctx, id)
	if err != nil {
		fmt.Println(err)
		return Comment{}, ErrFetchingCommentByID
	}

	return comment, nil
}

func (s *Service) UpdateComment(ctx context.Context, content Comment) error {
	return ErrNotImplemented
}

func (s *Service) DeleteComment(ctx context.Context, id string) error {
	return ErrNotImplemented
}