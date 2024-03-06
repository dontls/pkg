package orm

import (
	"time"
)

const (
	timeFormat = "2006-01-02"
	partFormat = "20060102"
)

type schemaPart struct {
	Name       string `gorm:"column:PARTITION_NAME;"`
	Desc       string `gorm:"column:PARTITION_DESCRIPTION;"`
	Expression string `gorm:"column:PARTITION_EXPRESSION;"`
	// Rows       uint   `gorm:"column:TABLE_ROWS;"`
}

type Partition struct {
	// 表名
	Table string
	// 分区字段
	Field string
	// 保留分区数
	Reserved int
}

func (o *Partition) schemaInfo(t time.Time) (string, string) {
	name := "p" + t.Format(partFormat)
	lessDay := t.AddDate(0, 0, 1).Format(timeFormat)
	return name, lessDay
}

func (p *Partition) queryTabelPart(data interface{}) error {
	return _db.Raw("SELECT PARTITION_NAME, PARTITION_DESCRIPTION, PARTITION_EXPRESSION "+
		"FROM information_schema.PARTITIONS WHERE table_name = ? ORDER BY PARTITION_NAME ASC;", p.Table).Scan(data).Error
}

// 返回当前最新分区日期，时间过滤
func (p *Partition) WithRange(rDays, interval int) string {
	parts, t0 := []schemaPart{}, time.Now()
	if err := p.queryTabelPart(&parts); err != nil {
		return ""
	}
	if parts[0].Name == "" {
		name, lessDay := p.schemaInfo(t0)
		_db.Exec("ALTER TABLE "+p.Table+" PARTITION BY RANGE (TO_DAYS("+p.Field+"))(PARTITION "+name+" VALUES LESS THAN (TO_DAYS(?)));", lessDay) // 初始化分区
		parts[0].Name = name
	}
	for {
		num := len(parts)
		t1, _ := time.Parse(partFormat, parts[num-1].Name[1:])
		if num*interval > rDays {
			_db.Exec("ALTER TABLE " + p.Table + " DROP PARTITION " + parts[0].Name) // 删除过时分区
			parts = parts[1:]
		}
		if t0.AddDate(0, 0, interval*p.Reserved).Compare(t1.AddDate(0, 0, 1)) <= 0 {
			return t1.Format(timeFormat)
		}
		// note, 这里再最新分区基础上interval添加新分区
		part, nextDay := p.schemaInfo(t1.AddDate(0, 0, interval))
		_db.Exec("ALTER TABLE "+p.Table+" ADD PARTITION(PARTITION "+part+" VALUES LESS THAN (TO_DAYS(?)));", nextDay)
		parts = append(parts, schemaPart{Name: part})
	}
}

func AlterPart(s, field string) *Partition {
	return &Partition{Table: s, Reserved: 2, Field: field}
}
