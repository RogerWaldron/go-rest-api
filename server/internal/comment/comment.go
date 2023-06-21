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
	GetComments(context.Context, int, int) ([]Comment, error)
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

func (s *Service) GetComments(ctx context.Context, limit int, offset int) ([]Comment, error) {
	comments, err := s.Store.GetComments(ctx, limit, offset)
	if err != nil {
		return []Comment{}, err
	}

	return comments, nil
}

func (s *Service) GetCommentByID(ctx context.Context, id string) (Comment, error) {
	comment, err := s.Store.GetCommentByID(ctx, id)
	if err != nil {
		fmt.Println(err)
		return Comment{}, ErrFetchingCommentByID
	}

	return comment, nil
}

func (s *Service) PostComment(ctx context.Context, newComment Comment) (Comment, error) {
	result, err := s.Store.PostComment(ctx, newComment)
	if err != nil {
		return Comment{}, err
	}

	return result, nil
}


func (s *Service) UpdateComment(ctx context.Context, id string, updated Comment) (Comment, error) {
	result, err := s.Store.UpdateComment(ctx, id, updated)
	if err != nil {
		return Comment{}, err
	}

	return result, nil
}

func (s *Service) DeleteComment(ctx context.Context, id string) error {
	return s.Store.DeleteComment(ctx, id)
}

