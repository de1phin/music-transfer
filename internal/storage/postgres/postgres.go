package postgres

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Postgres struct {
	*sqlx.DB
}

func NewPostgresDatabase(dataSourceName string) (*Postgres, error) {
	psql := new(Postgres)
	var err error
	psql.DB, err = sqlx.Open("postgres", dataSourceName)
	return psql, err
}
