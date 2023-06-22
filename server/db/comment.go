package db

import (
	"context"
	"database/sql"
	"fmt"

	// "github.com/RogerWaldron/go-rest-api/server/internal/comment"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type Comment struct {
	ID 			string 	`db:"id"`
	Slug 		string	`db:"slug"`
	Body 		string	`db:"body"`
	Author 	string	`db:"author"`
}

type CommentRow struct {
	ID 			string					
	Slug 		sql.NullString 	
	Body 		sql.NullString	
	Author 	sql.NullString
}

func convertCommentRowToComment(c CommentRow) Comment {
	return Comment{
		ID: c.ID,
		Slug: c.Slug.String,
		Body: c.Body.String,
		Author: c.Author.String,
	}
}

func (d *Database) GetCommentByID(
	ctx context.Context, 
	uuid string,
	) (Comment, error) {
		var cRow CommentRow
		  
		row := d.Client.QueryRowContext(
			ctx,
			`SELECT id, slug, body, author
			FROM comments
			WHERE id = $1`,
			uuid,
		)

		err := row.Scan(&cRow.ID, &cRow.Slug, &cRow.Body, &cRow.Author)
		if err != nil {
			log.Error().Err(err).Str("uuid", uuid).Msg("failed to get commment for uuid")
			return Comment{}, fmt.Errorf("failed fetching comment by uuid: %w", err)
		}

		return convertCommentRowToComment(cRow), nil
}

func (d *Database) GetComments(ctx context.Context, limit int, offset int) ([]Comment, error) {
	var (
		comments []Comment
	)

	rows, err := d.Client.QueryxContext(ctx, "SELECT * FROM comments LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return comments, fmt.Errorf("no comments found: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var row CommentRow
		err = rows.StructScan(&row)
		if err != nil {
			return comments, err
		}		
		comments = append(comments, convertCommentRowToComment(row))
	}

	return comments, rows.Err()
}

func (d *Database) PostComment(ctx context.Context, newComment Comment) (Comment, error) {
	newComment.ID = uuid.New().String()
	newEntry := CommentRow{
		ID: newComment.ID,
		Slug: sql.NullString{String: newComment.Slug, Valid: true},
		Body: sql.NullString{String: newComment.Body, Valid: true},
		Author: sql.NullString{String: newComment.Author, Valid: true},
	}

	rows, err := d.Client.NamedQueryContext(
		ctx,
		`INSERT INTO comments
		(id, slug, body, author) VALUES
		(:id, :slug, :body, :author)`,
		newEntry,
	)

	if err != nil {
		log.Error().
			Err(err).
			Str("id", newComment.ID).
			Str("slug", newComment.Slug).
			Str("body", newComment.Body).
			Str("author", newComment.Author).
			Msg("failed inserting comment")

		return Comment{}, fmt.Errorf("failed to insert comment: %w", err)
	}
	defer rows.Close()
	
	return convertCommentRowToComment(newEntry), nil
}

func (d *Database) UpdateComment(ctx context.Context, id string, newComment Comment) (Comment, error) {
	newEntry := CommentRow{
		ID: id,
		Slug: sql.NullString{String: newComment.Slug, Valid: true},
		Body: sql.NullString{String: newComment.Body, Valid: true},
		Author: sql.NullString{String: newComment.Author, Valid: true},
	}

	row, err := d.Client.NamedQueryContext(
		ctx,
		`UPDATE comments SET
		slug = :slug, 
		body = :body, 
		author = :author
		WHERE id = :id`,
		newEntry,
	)

	if err != nil {
		log.Error().
			Err(err).
			Str("id", newComment.ID).
			Str("slug", newComment.Slug).
			Str("body", newComment.Body).
			Str("author", newComment.Author).
			Msg("failed inserting comment")

		return Comment{}, fmt.Errorf("failed to insert comment: %w", err)
	}
	defer row.Close()

	return convertCommentRowToComment(newEntry), nil
}

func (d *Database) DeleteComment(ctx context.Context, deleteID string) error {
	_, err := d.Client.ExecContext(
		ctx,
		`DELETE FROM comments WHERE id = $1`,
		deleteID,
	 )
	if err != nil {
		log.Error().Err(err).Str("id", deleteID).Msg(`Delete failed`)
		return fmt.Errorf("delete failed: %w", err)
	}
	
	return nil
}