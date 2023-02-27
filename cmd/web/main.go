package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/scul0405/let-s-go-book/internal/models"
)

// Define an application struct to hold the application-wide dependencies for the web application
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	snippets *models.SnippetModel
}

func main() {
	// Get command-line value
	addr := flag.String("addr", ":3000", "HTTP network address")
	dsn := flag.String("dsn", "postgres://web:password@localhost:5432/snippetbox", "PostgreSQL data source name")

	flag.Parse()

	// Create customize loggers
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	conn, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	infoLog.Print("Successfully connect to database")

	defer conn.Close(context.Background())

	// Initialize a new instance of our application struct, containing the dependencies.
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		snippets: &models.SnippetModel{DB: conn},
	}

	// Initialize a new http.Server struct with customize logger
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(dsn string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
