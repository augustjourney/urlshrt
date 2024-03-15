package postgres

import (
	"context"
	"database/sql"

	"github.com/augustjourney/urlshrt/internal/storage"
)

type Repo struct {
	db *sql.DB
}

func (r *Repo) Init(ctx context.Context) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS urls (
			id SERIAL PRIMARY KEY NOT NULL,
			uuid VARCHAR(50) NOT NULL,
			short VARCHAR(50) NOT NULL,
			original VARCHAR NOT NULL,
			UNIQUE(short)
		)
	`)

	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.ExecContext(ctx, `
		CREATE INDEX IF NOT EXISTS short_idx ON urls (short);
	`)

	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *Repo) Create(ctx context.Context, url storage.URL) error {

	_, err := r.db.ExecContext(ctx, `
		insert into urls (uuid, short, original)
		values ($1, $2, $3)
	`, url.UUID, url.Short, url.Original)

	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) CreateBatch(ctx context.Context, urls []storage.URL) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	for _, url := range urls {
		_, err = tx.ExecContext(ctx, `
			insert into urls (uuid, short, original)
			values ($1, $2, $3)
		`, url.UUID, url.Short, url.Original)

		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (r *Repo) Get(ctx context.Context, short string) (*storage.URL, error) {
	var url storage.URL

	row := r.db.QueryRowContext(ctx, `
		select uuid, short, original 
		from urls
		where short = $1

	`, short)

	err := row.Scan(&url.UUID, &url.Short, &url.Original)

	if err != nil {
		return nil, err
	}

	return &url, nil
}

func New(ctx context.Context, db *sql.DB) (*Repo, error) {
	repo := Repo{
		db: db,
	}
	err := repo.Init(ctx)
	if err != nil {
		return nil, err
	}
	return &repo, nil
}
