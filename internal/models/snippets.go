package models

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
)

// Define a Snippet type to hold the data for an individual snippet.
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// Define a SnippetModel type which wraps a sql.DB connection pool.
type SnippetModel struct {
	DB *pgx.Conn
}

// This will insert a new snippet into the database.
func (m *SnippetModel) Insert(title string, content string, expires int) error {
	stmt := `INSERT INTO snippets (title, content, created, expires) VALUES($1, $2, $3, $4)`

	expiresTime := time.Now().AddDate(0, 0, expires)

	_, err := m.DB.Exec(context.Background(), stmt, title, content, time.Now(), expiresTime)
	if err != nil {
		return err
	}

	return nil
}

// This will return a specific snippet based on its id.
func (m *SnippetModel) Get(id int) (*Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets WHERE expires > NOW() AND id = $1`
	row := m.DB.QueryRow(context.Background(), stmt, id)

	snippet := &Snippet{}

	err := row.Scan(&snippet.ID, &snippet.Title, &snippet.Content, &snippet.Created, &snippet.Expires)
	/*
		We can use shorthand single-record queries like this:
		err := m.DB.QueryRow("SELECT ...", id).Scan(&snippet.ID, &snippet.Title, 
		&snippet.Content, &snippet.Created, &snippet.Expires)
	*/
	if err != nil {
		// In go 1.13 or newer prefer using errors.Is() instead of '=='
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return snippet, nil
}

// This will return the 10 most recently created snippets.
func (m *SnippetModel) Latest() ([]*Snippet, error) {
	return nil, nil
}
