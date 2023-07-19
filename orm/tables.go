// Copyright 2021 utils. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Sub-tables

package orm

import (
	"fmt"
	"reflect"
)

// type Tabs struct {
// 	ID    uint `gorm:"primarykey"`
// 	Name  string
// 	TabID uint
// }

// var gTlbs = 5

// func (o Tabs) TableName() string {
// 	l := len(strconv.Itoa(gTlbs - 1))
// 	return fmt.Sprintf("t_table_%0*d", l, o.TabID%uint(gTlbs))
// }
// func (Tabs) TableNameOf(id uint) string {
// 	return Tabs{TabID: id}.TableName()
// }

// func (Tabs) TableCount() uint {
// 	return uint(gTlbs)
// }

type XTablers interface {
	TableName() string       // 默认创建表
	TableNameOf(uint) string // 获取分区表
	TableCount() uint
}

//
// 请自行定
func CreateTables(v interface{}) {
	s, ok := v.(XTablers)
	if !ok {
		panic(fmt.Errorf("%v TypeOf not XTablers", reflect.TypeOf(v)))
	}
	_db.AutoMigrate(v)
	var i uint = 1
	for ; i < s.TableCount(); i++ {
		tablename := s.TableNameOf(i)
		if ok := _db.Migrator().HasTable(tablename); !ok {
			_db.Exec(fmt.Sprintf("CREATE TABLE %s LIKE %s;", tablename, s.TableName()))
		}
	}
}
