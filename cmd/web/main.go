package main

import (
	"context"
	"flag"
	"github.com/jackc/pgx/v4"
	"log"
	"net/http"
	"os"
	"se03.com/pkg/models/postgresql"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	snippets *postgresql.SnippetModel
	debug    bool
}

func main() {
	//Router
	addr := flag.String("addr", ":4000", "HTTP network address")
	//postgres://postgres:alimzhan125@localhost:5432/snippetbox
	flag.Parse()
	ctx := context.Background()
	connStr := "postgres://postgres:alimzhan125@localhost:5432/snippetbox"
	dsn := flag.String("dsn", connStr, "PostgreSQL data source name")
	db, err := openDB(*dsn, ctx)
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close(ctx)

	app := application{
		errorLog: errorLog,
		infoLog:  infoLog,
		snippets: &postgresql.SnippetModel{DB: db, Ctx: ctx},
	}

	srv := &http.Server{
		Addr:     *addr,
		Handler:  app.routes(),
		ErrorLog: errorLog,
	}
	infoLog.Printf("Starting  server on %v", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)

}

func openDB(dsn string, ctx context.Context) (*pgx.Conn, error) {
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return nil, err
	}
	if err = conn.Ping(ctx); err != nil {
		return nil, err
	}
	return conn, nil
}
