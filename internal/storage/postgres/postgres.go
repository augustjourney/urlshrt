package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/augustjourney/urlshrt/internal/storage"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
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

	_, err = tx.ExecContext(ctx, `
		CREATE UNIQUE INDEX IF NOT EXISTS original_unique_idx ON urls (original);
	`)

	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.ExecContext(ctx, `
		ALTER TABLE urls ADD COLUMN IF NOT EXISTS user_uuid VARCHAR;
	`)

	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.ExecContext(ctx, `
		ALTER TABLE urls ADD COLUMN IF NOT EXISTS is_deleted BOOLEAN DEFAULT false;
	`)

	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *Repo) Create(ctx context.Context, url storage.URL) error {

	_, err := r.db.ExecContext(ctx, `
		insert into urls (uuid, short, original, user_uuid)
		values ($1, $2, $3, $4)
	`, url.UUID, url.Short, url.Original, url.UserUUID)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			// Если такой url уже есть
			// Возвращаем ошибку ErrAlreadyExists
			if pgErr.Code == pgerrcode.UniqueViolation {
				return storage.ErrAlreadyExists
			}
		}
		return err
	}

	return err
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

func (r *Repo) DeleteBatch(ctx context.Context, shortIds []string, userId string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	for _, shortId := range shortIds {
		_, err := tx.ExecContext(ctx, `
			update urls
			set is_deleted = true
			where user_uuid = $1 and short = $2
		`, userId, shortId)

		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (r *Repo) GetByOriginal(ctx context.Context, original string) (*storage.URL, error) {
	var url storage.URL

	row := r.db.QueryRowContext(ctx, `
		select uuid, short, original 
		from urls
		where original = $1

	`, original)

	err := row.Scan(&url.UUID, &url.Short, &url.Original)

	if err != nil {
		return nil, err
	}

	return &url, nil
}

func (r *Repo) Get(ctx context.Context, short string) (*storage.URL, error) {
	var url storage.URL

	row := r.db.QueryRowContext(ctx, `
		select uuid, short, original, is_deleted
		from urls
		where short = $1

	`, short)

	err := row.Scan(&url.UUID, &url.Short, &url.Original, &url.IsDeleted)

	if err != nil {
		return nil, err
	}

	return &url, nil
}

func (r *Repo) GetByUserUUID(ctx context.Context, userUUID string) (*[]storage.URL, error) {
	var urls []storage.URL

	rows, err := r.db.QueryContext(ctx, `
		select short, original 
		from urls
		where user_uuid = $1

	`, userUUID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var url storage.URL
		err = rows.Scan(&url.Short, &url.Original)
		if err != nil {
			return nil, err
		}

		urls = append(urls, url)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return &urls, nil
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
