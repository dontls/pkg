package sqlr

import (
	"database/sql"
	"errors"
	"reflect"
)

var errDbOpened = errors.New("db open error")

type DB struct {
	Db *sql.DB
}

func (o *DB) Exec(query string, args ...any) (sql.Result, error) {
	if o.Db == nil {
		return nil, errDbOpened
	}
	return o.Db.Exec(query, args...)
}

func (o *DB) Query(query string, args ...any) (*sql.Rows, error) {
	if o.Db == nil {
		return nil, errDbOpened
	}
	return o.Db.Query(query, args...)
}

func (o *DB) Count(dest any, query string, args ...any) error {
	if o.Db == nil {
		return errDbOpened
	}
	rows, err := o.Db.Query(query, args...)
	if err != nil {
		return err
	}
	for rows.Next() {
		return rows.Scan(dest)
	}
	return nil
}

func (o *DB) Get(dest any, query string, args ...any) error {
	if o.Db == nil {
		return errDbOpened
	}
	rows, err := o.Db.Query(query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()
	fields := SqlValuesAddr(dest)
	for rows.Next() {
		return rows.Scan(fields...)
	}
	return nil
}

func (o *DB) Find(dest any, query string, args ...any) error {
	if o.Db == nil {
		return errDbOpened
	}
	rows, err := o.Db.Query(query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()
	value := _reflectValue(dest)
	t := reflect.New(value.Type().Elem())
	cols, _ := rows.Columns()
	fields := SqlValuesAddr(t.Interface())[:len(cols)]
	for rows.Next() {
		if err := rows.Scan(fields...); err != nil {
			continue
		}
		// fmt.Printf("%#v\n", t)
		value.Set(reflect.Append(value, t.Elem()))
	}
	return nil
}

func (o *DB) Insert(v any, table string) (sql.Result, error) {
	if o.Db == nil {
		return nil, errDbOpened
	}
	fields := SqlValueNames(v)
	sqls := SqlValues(v, nil)
	return o.Db.Exec("INSERT INTO " + table + " " + fields + " VALUES " + sqls)
}
