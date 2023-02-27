package models

import (
	"context"
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
	return nil, nil
}

// This will return the 10 most recently created snippets.
func (m *SnippetModel) Latest() ([]*Snippet, error) {
	return nil, nil
}
