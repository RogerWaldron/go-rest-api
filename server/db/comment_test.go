package db

import (
	"context"
	"testing"

	"github.com/RogerWaldron/go-rest-api/server/internal/comment"
	"github.com/stretchr/testify/assert"
)

const (
	DB_CONNECTION = "host=0.0.0.0 port=5432 user=postgres dbname=postgres password=postgres sslmode=disable"
)

func TestCommentDatabase(t *testing.T) {
	t.Run("Can create comment", func(t *testing.T){
		db, err := NewDatabase(DB_CONNECTION)
		assert.NoError(t, err)

		cmt, err := db.PostComment(context.Background(), 
			comment.Comment{
				Slug: "postContent",
				Author: "author",
				Body: "body",
		}) 
		assert.NoError(t, err)

		newCmt, err := db.GetCommentByID(context.Background(), cmt.ID)
		assert.NoError(t, err)
		assert.Equal(t, cmt.Slug, newCmt.Slug)
	})

	t.Run("Can update a comment", func(t *testing.T){
			db, err := NewDatabase(DB_CONNECTION)
			assert.NoError(t, err)

			cmt, err := db.PostComment(context.Background(), comment.Comment{
				Slug: "updateContent",
				Author: "author",
				Body: "old body",
			})
			assert.NoError(t, err)

			cmt.Body = "new body"
			cmt, err = db.UpdateComment(context.Background(), cmt.ID, cmt)
			assert.NoError(t, err)

			newCmt, err := db.GetCommentByID(context.Background(), cmt.ID)
			assert.NoError(t, err)
			assert.Equal(t, cmt.Body, newCmt.Body)
	})

	t.Run("Can get all Comments by ID", func(t *testing.T) {
		db, err := NewDatabase(DB_CONNECTION)
		assert.NoError(t, err)

		cmt, err := db.PostComment(context.Background(), comment.Comment{
			Slug: "GetCommentsByID",
			Author: "author",
			Body: "body",
		})
		assert.NoError(t, err)

		newCmt, err := db.GetCommentByID(context.Background(), cmt.ID)
		assert.Equal(t, cmt.Slug, newCmt.Slug)
		assert.NoError(t, err)
	})

	t.Run("Can get all Comments", func(t *testing.T){
		db, err := NewDatabase(DB_CONNECTION)
		assert.NoError(t, err)

		existingComments, err := db.GetComments(context.Background(), 0, 0)
		assert.NoError(t, err)

		newCmt, err := db.PostComment(context.Background(), comment.Comment{
			Slug: "GetComments",
			Author: "author",
			Body: "body",
		})
		assert.NoError(t, err)
		assert.Equal(t, "author", newCmt.Author)

		currentComments, err := db.GetComments(context.Background(), 0, 0)
		assert.NoError(t, err)
		assert.Greater(t, len(currentComments), len(existingComments))
	})

	t.Run("Can delete a comment", func(t *testing.T) {
		db, err := NewDatabase(DB_CONNECTION)
		assert.NoError(t, err)

		newCmt, err := db.PostComment(context.Background(), comment.Comment{
			Slug: "DeleteComment",
			Author: "author",
			Body: "body",
		})
		assert.NoError(t, err)
		cmt, err := db.GetCommentByID(context.Background(), newCmt.ID)
		assert.NoError(t, err)
		
		err = db.DeleteComment(context.Background(), cmt.ID)
		assert.NoError(t, err)

		cmt, err = db.GetCommentByID(context.Background(), cmt.ID)
		assert.Error(t, err)
		assert.Empty(t, cmt)
	})
}