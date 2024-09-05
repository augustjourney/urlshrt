// модуль storage отвечает за сохранение данных о ссылках.
// поддерживается сохранение в памяти, в файле и базе данных postgres.
package storage

import (
	"context"
	"errors"
)

// хранит информацию о ссылке
type URL struct {
	UUID      string `json:"uuid,omitempty"`
	Short     string `json:"short_url"`
	Original  string `json:"original_url"`
	UserUUID  string `json:"user_uuid,omitempty"`
	IsDeleted bool
}

// хранит информацию о статистике:
// количество сохраненных ссылок
// и количество пользователей
type Stats struct {
	UrlsCount  int `json:"urls" db:"urls_count"`
	UsersCount int `json:"users" db:"users_count"`
}

// описывает методы хранилища
type IRepo interface {
	Create(ctx context.Context, url URL) error
	Get(ctx context.Context, short string) (*URL, error)
	GetByOriginal(ctx context.Context, original string) (*URL, error)
	CreateBatch(ctx context.Context, urls []URL) error
	GetByUserUUID(ctx context.Context, userUUID string) (*[]URL, error)
	Delete(ctx context.Context, short []string, userID string) error
	GetStats(ctx context.Context) (Stats, error)
}

// ошибка если ссылка уже существует
var ErrAlreadyExists = errors.New("URL already exists")
