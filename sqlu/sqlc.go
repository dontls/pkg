package sqlu

import (
	"fmt"
	"strings"
)

type sqlc struct {
	sqls   string
	wheres []string // sqlc集合
	others []string // 附加条件 ORDER BY/LIMIT etc
}

func NewSQL(s string) *sqlc {
	return &sqlc{sqls: s}
}

func (o *sqlc) Where(query string) *sqlc {
	if query != "" {
		o.wheres = append(o.wheres, query)
	}
	return o
}

func (o *sqlc) Equal(query string, v interface{}) *sqlc {
	switch v := v.(type) {
	case string:
		if v != "" {
			o.wheres = append(o.wheres, fmt.Sprintf("%s = '%s'", query, v))
		}
	case int:
		o.wheres = append(o.wheres, fmt.Sprintf("%s = %d", query, v))
	}
	return o
}

func (o *sqlc) In(query string, v interface{}) *sqlc {
	switch v := v.(type) {
	case string:
		if v != "" {
			s := (query + " IN(" + v + ")")
			o.wheres = append(o.wheres, s)
		}
	}
	return o
}

func (o *sqlc) Between(query string, ss ...string) *sqlc {
	if len(ss) == 2 {
		s := (query + " BETWEEN '" + ss[0] + "' AND '" + ss[1] + "'")
		o.wheres = append(o.wheres, s)
	}
	return o
}

func (o *sqlc) Limit(num, size int) *sqlc {
	if size > 0 {
		o.others = append(o.others, fmt.Sprintf("LIMIT %d,%d", (num-1)*size, size))
	}
	return o
}

func (o *sqlc) OrderBy(query string, s string) *sqlc {
	o.others = append(o.others, "ORDER BY "+query+" "+s)
	return o
}

func (o *sqlc) GroupBy(query string) *sqlc {
	o.others = append(o.others, "GROUP BY "+query)
	return o
}

func (o *sqlc) String() string {
	if len(o.wheres) > 0 {
		o.sqls += (" WHERE " + strings.Join(o.wheres, " AND "))
	}
	if len(o.others) > 0 {
		o.sqls += (" " + strings.Join(o.others, " "))
	}
	return o.sqls
}
