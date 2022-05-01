package postgres

import (
	"errors"
	"fmt"
	"strings"
)

type Table[Key comparable, T any] struct {
	psql          *Postgres
	name          string
	key           string
	columnsJoined string
	updateColumns string
	insertColumns string
}

func NewTable[Key comparable, T any](psql *Postgres, name string, key string) *Table[Key, T] {
	var val T
	columns := reflectColumns(val)
	updateColumns := ""
	for i, c := range columns {
		if i > 0 {
			updateColumns += ", "
		}
		updateColumns += c + "=:" + c
	}
	insertColumns := ":" + strings.Join(columns, ", :")
	return &Table[Key, T]{
		psql:          psql,
		name:          name,
		key:           key,
		columnsJoined: strings.Join(columns, ", "),
		updateColumns: updateColumns,
		insertColumns: insertColumns,
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
	query := "SELECT " + t.columnsJoined + " FROM " + t.name + " WHERE " + t.key +
		" = $1"
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

func (t *Table[Key, T]) Set(key Key, val T) error {
	exists, err := t.Exist(key)
	if err != nil {
		return err
	}

	if exists {
		query := "UPDATE " + t.name + " SET " + t.updateColumns + " WHERE " + t.key + " = " + fmt.Sprintf("%v", key)
		_, err = t.psql.NamedExec(query, val)
		if err != nil {
			return err
		}
	} else {
		query := "INSERT INTO " + t.name + "(" + t.key + ", " + t.columnsJoined +
			") VALUES(" + fmt.Sprintf("%v", key) + ", " + t.insertColumns + ")"
		_, err = t.psql.NamedExec(query, &val)
		if err != nil {
			return err
		}
	}

	return nil
}
