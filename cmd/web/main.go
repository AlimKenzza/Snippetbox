package main

import (
	"context"
	"flag"
	"github.com/golangcollege/sessions"
	"github.com/jackc/pgx/v4"
	"html/template"
	"log"
	"net/http"
	"os"
	"se03.com/pkg/models/postgresql"
	"time"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	session       *sessions.Session
	snippets      *postgresql.SnippetModel
	debug         bool
	templateCache map[string]*template.Template
}

func main() {
	//Router
	addr := flag.String("addr", ":4000", "HTTP network address")
	secret := flag.String("secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret key")
	//postgres://postgres:alimzhan125@localhost:5432/snippetbox
	flag.Parse()
	ctx := context.Background()
	connStr := "postgres://postgres:alimzhan125@host.docker.internal:5432/snippetbox"
	dsn := flag.String("dsn", connStr, "PostgreSQL data source name")
	db, err := openDB(*dsn, ctx)
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close(ctx)

	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}
	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour
	app := application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		session:       session,
		snippets:      &postgresql.SnippetModel{DB: db, Ctx: ctx},
		templateCache: templateCache,
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
