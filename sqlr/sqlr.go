package sqlr

import "database/sql"

var _db = &DB{}

func SetDefault(db *sql.DB) {
	_db.Db = db
}

func Exec(query string, args ...any) (sql.Result, error) {
	return _db.Exec(query, args...)
}

func Query(query string, args ...any) (*sql.Rows, error) {
	return _db.Query(query, args...)
}

func Count(dest any, query string, args ...any) error {
	return _db.Count(dest, query, args...)
}

func Get(dest any, query string, args ...any) error {
	return _db.Get(dest, query, args...)
}

func Find(dest any, query string, args ...any) error {
	return _db.Find(dest, query, args...)
}

func Insert(v any, table string) (sql.Result, error) {
	return _db.Insert(v, table)
}
