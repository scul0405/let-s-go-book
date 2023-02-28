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
	stmt := `SELECT id, title, content, created, expires FROM snippets 
	WHERE expires > NOW() ORDER BY id DESC LIMIT 10`

	rows, err := m.DB.Query(context.Background(), stmt)
	if err != nil {
		return nil, err
	}

	// Closing a resultset with defer rows.Close() is critical. As long as 
	// a resultset is open it will keep the underlying database connection open… so if
	// something goes wrong in this method and the resultset isn’t closed, it can rapidly lead
	// to all the connections in your pool being used up.
	defer rows.Close()

	// Empty slice to hold snippet structs
	snippets := []*Snippet{}

	for rows.Next() {
		snippet := &Snippet{}

		err = rows.Scan(&snippet.ID, &snippet.Title, &snippet.Content, &snippet.Created, &snippet.Expires)
		if err != nil {
			return nil, err
		}
		// Append it to the slice of snippets.
		snippets = append(snippets, snippet)
	}

	// When the rows.Next() loop has finished we call rows.Err() to retrieve any
	// error that was encountered during the iteration. It's important to
	// call this - don't assume that a successful iteration was completed
	// over the whole resultset.
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
