package sqlu

import (
	"database/sql"
	"errors"
	"reflect"
)

var errDbOpened = errors.New("db open error")

type DB struct {
	db *sql.DB
}

func (o *DB) Exec(query string, args ...interface{}) (sql.Result, error) {
	if o.db == nil {
		return nil, errDbOpened
	}
	return o.db.Exec(query, args...)
}

func (o *DB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	if o.db == nil {
		return nil, errDbOpened
	}
	return o.db.Query(query, args...)
}

func (o *DB) Count(dest interface{}, query string, args ...interface{}) error {
	if o.db == nil {
		return errDbOpened
	}
	rows, err := o.db.Query(query, args...)
	if err != nil {
		return err
	}
	for rows.Next() {
		return rows.Scan(dest)
	}
	return nil
}

func (o *DB) Get(dest interface{}, query string, args ...interface{}) error {
	if o.db == nil {
		return errDbOpened
	}
	rows, err := o.db.Query(query, args...)
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

// 反射生成slice的基本类型对象
func (o *DB) SliceTypeValue(dest interface{}) interface{} {
	value := _reflectValue(dest)
	if value.Type().Kind() == reflect.Slice {
		return reflect.New(value.Type().Elem()).Interface()
	}
	return nil
}

func (o *DB) Find(dest interface{}, query string, args ...interface{}) error {
	if o.db == nil {
		return errDbOpened
	}
	value := _reflectValue(dest)
	rows, err := o.db.Query(query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()
	t := reflect.New(value.Type().Elem())
	fields := SqlValuesAddr(t.Interface())
	for rows.Next() {
		if err := rows.Scan(fields...); err != nil {
			continue
		}
		// fmt.Printf("%#v\n", t)
		value.Set(reflect.Append(value, t.Elem()))
	}
	return nil
}

func (o *DB) Insert(v interface{}, table string) (sql.Result, error) {
	if o.db == nil {
		return nil, errDbOpened
	}
	fields := SqlValueNames(v)
	sqls := SqlValues(v, nil)
	return o.db.Exec("INSERT INTO " + table + " " + fields + " VALUES " + sqls)
}
