package orm

import (
	"fmt"
	"reflect"
	"time"

	"gorm.io/gorm/schema"
)

const (
	timeFormat = "2006-01-02 15:04:05"
	partFormat = "20060102 150405"
)

type schemaPartition struct {
	Name       string `gorm:"column:PARTITION_NAME;"`
	Desc       string `gorm:"column:PARTITION_DESCRIPTION;"`
	Expression string `gorm:"column:PARTITION_EXPRESSION;"`
	// Rows       uint   `gorm:"column:TABLE_ROWS;"`
}

type partition struct {
	// 表名
	Table string
}

func (o *partition) schemaName(t time.Time) (string, string) {
	name := "p" + t.Format(partFormat)[:8]
	where := t.AddDate(0, 0, 1).Format(timeFormat)[:10]
	return name, where
}

const queryPart = `SELECT PARTITION_NAME, PARTITION_DESCRIPTION, PARTITION_EXPRESSION FROM information_schema.PARTITIONS WHERE table_name = '%s' ORDER BY PARTITION_NAME DESC;`
const initPart = `ALTER TABLE %s PARTITION BY RANGE (TO_DAYS(%s))(PARTITION %s VALUES LESS THAN (TO_DAYS('%s')));`
const dropPart = `ALTER TABLE %s DROP PARTITION %s;`
const createPart = `ALTER TABLE %s ADD PARTITION(PARTITION %s VALUES LESS THAN (TO_DAYS('%s')));`

func (o *partition) queryAll(data interface{}) error {
	return _db.Raw(fmt.Sprintf(queryPart, o.Table)).Scan(data).Error
}

func (o *partition) AlterRange(rDays, interval int, field string) {
	t := time.Now()
	pname, lessDay := o.schemaName(t)
	var data []schemaPartition
	o.queryAll(&data)
	if data[0].Name == "" {
		_db.Exec(fmt.Sprintf(initPart, o.Table, field, pname, lessDay))
		data[0].Name = pname
		return
	}
	num := len(data)
	if num*interval > rDays {
		_db.Exec(fmt.Sprintf(dropPart, o.Table, data[num-1].Name)) // 删除过时分区
	}
	if data[0].Name > pname {
		return
	}
	// 中途启动创建当前分区
	if data[0].Name < pname {
		_db.Exec(fmt.Sprintf(createPart, o.Table, pname, lessDay))
	}
	nextp, nextLessDay := o.schemaName(t.AddDate(0, 0, interval))
	_db.Exec(fmt.Sprintf(createPart, o.Table, nextp, nextLessDay))
}

func Partition(m interface{}) *partition {
	v, ok := m.(schema.Tabler)
	if ok {
		return &partition{Table: v.TableName()}
	}
	panic(fmt.Errorf("%v typeOf not gorm schema.Tabler", reflect.TypeOf(v)))
}
