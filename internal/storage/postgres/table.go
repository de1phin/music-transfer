package postgres

import (
	"errors"
	"fmt"
	"strings"
)

type Table[Key comparable, T any] struct {
	psql    *Postgres
	columns []string
	name    string
	key     string
}

func NewTable[Key comparable, T any](psql *Postgres, name string, key string) *Table[Key, T] {
	var val T
	return &Table[Key, T]{
		psql:    psql,
		name:    name,
		key:     key,
		columns: reflectColumns(val),
	}
}

func (t *Table[Key, T]) Exist(key Key) (bool, error) {
	rows, err := t.psql.Queryx("SELECT 1 AS exist FROM "+t.name+" WHERE "+t.key+" = $1", &key)
	if err != nil {
		return false, err
	}
	ok := rows.Next()
	if !ok {
		return false, nil
	}
	return true, nil
}

func (t *Table[Key, T]) Get(key Key) (T, error) {
	result := []T{}
	query := "SELECT " + strings.Join(t.columns, ", ") + " FROM " + t.name + " WHERE " + t.key + " = $1"
	err := t.psql.Select(&result, query, &key)
	if err != nil {
		var r T
		return r, err
	}
	if result == nil || len(result) == 0 {
		var r T
		return r, errors.New(fmt.Sprintf("Postgres.Get: Row with specified key (%v) does not exist", key))
	}
	if len(result) > 1 {
		var r T
		return r, errors.New(fmt.Sprintf("Postgres.Get: Multiple entries for key %v", key))
	}
	return result[0], nil
}

func (t *Table[Key, T]) Put(key Key, val T) error {
	ok, err := t.Exist(key)
	if err != nil {
		return err
	}

	if !ok {
		_, err = t.psql.Exec("INSERT INTO "+t.name+" ("+t.key+") VALUES($1)", &key)
		if err != nil {
			return err
		}
	}
	query := "UPDATE " + t.name + " SET "
	for i, c := range t.columns {
		if i > 0 {
			query += ", "
		}
		query += c + "=:" + c
	}
	query += " WHERE " + t.key + " = " + fmt.Sprintf("%v", key)
	fmt.Println(query)
	fmt.Println("Put", val)
	_, err = t.psql.NamedExec(query, val)
	return err
}
