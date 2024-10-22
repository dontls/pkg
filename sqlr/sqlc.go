package sqlr

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

func (s *sqlc) Where(query string) *sqlc {
	if query != "" {
		s.wheres = append(s.wheres, query)
	}
	return s
}

type String string

func (s *sqlc) Equal(query string, v any) *sqlc {
	s1 := ""
	switch v := v.(type) {
	case string:
		if v != "" {
			s1 = fmt.Sprintf("%s='%s'", query, v)
		}
	case String:
		v1 := string(v)
		if v1 != "" {
			s1 = fmt.Sprintf("%s='%s'", query, v1)
		}
	case int:
		s1 = fmt.Sprintf("%s=%d", query, v)
	}
	if s1 != "" {
		s.wheres = append(s.wheres, s1)
	}
	return s
}

// 1,2,3
// ab,cd,ef，类型使用sqlr.String
// ["ab", "cd"]
func (s *sqlc) In(query string, v any) *sqlc {
	s1 := ""
	switch v := v.(type) {
	case string:
		if v != "" {
			s1 = (query + " IN(" + v + ")")
		}
	case String:
		v1 := string(v)
		if v1 != "" {
			s1 = (query + " IN('" + strings.ReplaceAll(v1, ",", "','") + "')")
		}
	case []string:
		if len(v) > 0 {
			s1 = (query + " IN('" + strings.Join(v, "','") + "')")
		}
	}
	if s1 != "" {
		s.wheres = append(s.wheres, s1)
	}
	return s
}

func (s *sqlc) Between(query string, ss ...string) *sqlc {
	if len(ss) == 2 {
		s1 := (query + " BETWEEN '" + ss[0] + "' AND '" + ss[1] + "'")
		s.wheres = append(s.wheres, s1)
	}
	return s
}

func (s *sqlc) Limit(num, size int) *sqlc {
	if size > 0 {
		s.others = append(s.others, fmt.Sprintf("LIMIT %d,%d", (num-1)*size, size))
	}
	return s
}

func (s *sqlc) OrderBy(query string, v string) *sqlc {
	s.others = append(s.others, "ORDER BY "+query+" "+v)
	return s
}

func (s *sqlc) GroupBy(query string) *sqlc {
	s.others = append(s.others, "GROUP BY "+query)
	return s
}

func (s *sqlc) Select(v string) string {
	if len(s.wheres) > 0 {
		v += (" WHERE " + strings.Join(s.wheres, " AND "))
	}
	if len(s.others) > 0 {
		v += (" " + strings.Join(s.others, " "))
	}
	return v
}

func (s *sqlc) String() string {
	return s.Select(s.sqls)
}
