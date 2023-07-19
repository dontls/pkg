package orm

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// "root:Howenforever@tcp(172.16.60.219:35200)/jtdata?charset=utf8&parseTime=True&loc=Local"
// 去掉时区，使用钩子函数更新时间
const timeformat = "2006-01-02 15:04:05"

// 去掉parseTime&loc=Local
// jtime format json time field by myself
type CreatedAt struct {
	CreatedAt string `json:"createdAt" gorm:"type:datetime;"`
}

func (CreatedAt) BeforeCreate(tx *gorm.DB) error {
	tx.Statement.SetColumn("CreatedAt", time.Now().Format(timeformat))
	return nil
}

type UpdatedAt struct {
	UpdatedAt string `json:"UpdatedAt" gorm:"type:datetime;"`
}

func (UpdatedAt) BeforeUpdate(tx *gorm.DB) error {
	tx.Statement.SetColumn("UpdatedAt", time.Now().Format(timeformat))
	return nil
}

type StringArray []string

func (p StringArray) Value() (driver.Value, error) {
	return json.Marshal(p)
}

func (p *StringArray) Scan(data interface{}) error {
	return json.Unmarshal(data.([]byte), &p)
}

type UintArray []uint

// Value insert
func (j UintArray) Value() (driver.Value, error) {
	return json.Marshal(&j)
}

// Scan valueof
func (t *UintArray) Scan(v interface{}) error {
	return json.Unmarshal(v.([]byte), t)
}
