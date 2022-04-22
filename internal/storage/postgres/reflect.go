package postgres

import (
	"reflect"
)

func reflectColumns(val any) []string {
	fields := reflect.VisibleFields(reflect.TypeOf(val))
	columns := []string{}
	for _, f := range fields {
		if tag := f.Tag.Get("db"); tag != "" {
			columns = append(columns, tag)
		} else {
			columns = append(columns, f.Name)
		}
	}
	return columns
}
