package orm

import (
	"fmt"
	"reflect"
	"time"

	"gorm.io/gorm/schema"
)

const (
	timeFormat = "2006-01-02"
	partFormat = "20060102"
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
	// 保留分区数
	Reserved int
}

func (o *partition) schemaName(t time.Time) (string, string) {
	name := "p" + t.Format(partFormat)
	lessDay := t.AddDate(0, 0, 1).Format(timeFormat)
	return name, lessDay
}

const queryPart = `SELECT PARTITION_NAME, PARTITION_DESCRIPTION, PARTITION_EXPRESSION FROM information_schema.PARTITIONS WHERE table_name = '%s' ORDER BY PARTITION_NAME DESC;`
const initPart = `ALTER TABLE %s PARTITION BY RANGE (TO_DAYS(%s))(PARTITION %s VALUES LESS THAN (TO_DAYS('%s')));`
const dropPart = `ALTER TABLE %s DROP PARTITION %s;`
const createPart = `ALTER TABLE %s ADD PARTITION(PARTITION %s VALUES LESS THAN (TO_DAYS('%s')));`

func (p *partition) queryAll(data interface{}) error {
	return _db.Raw(fmt.Sprintf(queryPart, p.Table)).Scan(data).Error
}

// 返回当前最新分区日期，时间过滤
func (p *partition) AlterRange(rDays, interval int, field string) string {
	data, t0 := []schemaPartition{}, time.Now()
	p.queryAll(&data)
	if data[0].Name == "" {
		name, lessDay := p.schemaName(t0)
		_db.Exec(fmt.Sprintf(initPart, p.Table, field, name, lessDay)) // 初始化分区
		data[0].Name = name
	}
	num := len(data)
	if num*interval > rDays {
		_db.Exec(fmt.Sprintf(dropPart, p.Table, data[num-1].Name)) // 删除过时分区
	}
	t1, _ := time.Parse(partFormat, data[0].Name[1:])
	if t0.AddDate(0, 0, interval*p.Reserved).Compare(t1.AddDate(0, 0, 1)) <= 0 {
		return t1.Format(timeFormat) // 返回最新分区日期
	}
	// note, 这里再最新分区基础上interval添加新分区
	nextp, nextLessDay := p.schemaName(t1.AddDate(0, 0, interval))
	_db.Exec(fmt.Sprintf(createPart, p.Table, nextp, nextLessDay))
	return nextLessDay
}

func Partition(m interface{}) *partition {
	if v, ok := m.(schema.Tabler); ok {
		return &partition{Table: v.TableName(), Reserved: 2}
	}
	panic(fmt.Errorf("%v typeOf not gorm schema.Tabler", reflect.TypeOf(m)))
}
