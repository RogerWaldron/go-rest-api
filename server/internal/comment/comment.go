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

type CommentStore interface {
	GetComments(context.Context) ([]Comment, error)
	GetCommentByID(context.Context, string) (Comment, error)
	PostComment(context.Context, Comment) (Comment, error)
	UpdateComment(context.Context, string, Comment) (Comment, error)
	DeleteComment(context.Context, string) error
	Ping(context.Context) error
}

type Comment struct {
	ID 			string 	`db:"id"`
	Slug 		string	`db:"slug"`
	Body 		string	`db:"body"`
	Author 	string	`db:"author"`
}

type Service struct{
	Store CommentStore
}

func NewService(store CommentStore) *Service {
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

func (s *Service) UpdateComment(ctx context.Context, content Comment) (Comment, error) {
	return Comment{}, nil
}

func (s *Service) DeleteComment(ctx context.Context, id string) error {
	return ErrNotImplemented
}

