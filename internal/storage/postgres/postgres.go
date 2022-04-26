package postgres

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Postgres struct {
	*sqlx.DB
}

func NewPostgresDatabase(config Config) (psql *Postgres, err error) {
	psql = new(Postgres)
	dataSourceName := "user=" + config.User + " password=" + config.Password + " dbname=" + config.DBName + " sslmode=" + config.SSLMode
	psql.DB, err = sqlx.Open("postgres", dataSourceName)
	return psql, err
}
