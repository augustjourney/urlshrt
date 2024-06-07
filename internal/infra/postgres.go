// пакет infra хранит в себе подключения к внешним зависимостям: базы данных, кэш и т.п.
package infra

import (
	"database/sql"

	"github.com/augustjourney/urlshrt/internal/config"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// Создает подключение в postgres
func InitPostgres(config *config.Config) (*sql.DB, error) {
	db, err := sql.Open("pgx", config.DatabaseDSN)
	if err != nil {
		return nil, err
	}
	return db, nil
}
