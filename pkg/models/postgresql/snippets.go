package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgx/v4"
	"se03.com/pkg/models"
	"strconv"
	"time"
)

type SnippetModel struct {
	DB  *pgx.Conn
	Ctx context.Context
}

func (m *SnippetModel) Insert(title, content, expires string) (int, error) {

	stmt := "INSERT INTO snippets (title, content, created, expires) VALUES ( $1,$2, $3 ,$4) RETURNING id"
	days, _ := strconv.Atoi(expires)
	created := time.Now()
	id := 0
	err := m.DB.QueryRow(context.Background(), stmt, title, content, created, time.Now().AddDate(0, 0, days)).Scan(&id)
	if err != nil {
		return 0, nil
	}
	return id, nil
}

func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	s := &models.Snippet{}
	err := m.DB.QueryRow(m.Ctx, "Select id, title, content, created, expires FROM snippets WHERE expires > NOW() AND id = $1", id).Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}

func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
       WHERE expires > NOW() ORDER BY created DESC LIMIT 10`
	rows, err := m.DB.Query(m.Ctx, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	snippets := []*models.Snippet{}
	for rows.Next() {
		s := &models.Snippet{}
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return snippets, nil

}
