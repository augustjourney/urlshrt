package postgres

import (
	"context"
	"database/sql"

	"github.com/augustjourney/urlshrt/internal/logger"
	"github.com/augustjourney/urlshrt/internal/storage"
	"github.com/google/uuid"
)

type Repo struct {
	db *sql.DB
}

func (r *Repo) Init(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS urls (
		    id SERIAL PRIMARY KEY NOT NULL,
		    uuid VARCHAR(50) NOT NULL,
		    short VARCHAR(50) NOT NULL,
		    original VARCHAR NOT NULL,
			UNIQUE(short)
		)
	`)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) Create(ctx context.Context, short string, original string) error {
	uuid, err := uuid.NewRandom()

	if err != nil {
		logger.Log.Error("Could not create uuid ", err)
		return err
	}

	url := storage.URL{
		UUID:     uuid.String(),
		Short:    short,
		Original: original,
	}

	_, err = r.db.ExecContext(ctx, `
		insert into urls (uuid, short, original)
		values ($1, $2, $3)
	`, url.UUID, url.Short, url.Original)

	if err != nil {
		return err
	}

	return nil
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
