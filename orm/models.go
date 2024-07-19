package orm

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// 去掉时区，使用钩子函数更新时间
const TimeFormat = "2006-01-02 15:04:05"

// Time format json time field by myself
type Time time.Time

// MarshalJSON on JSONTime format Time field with %Y-%m-%d %H:%M:%S
func (t Time) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%s\"", time.Time(t).Format(TimeFormat))
	return []byte(formatted), nil
}

// Value insert timestamp into mysql need this function.
func (t Time) Value() (driver.Value, error) {
	var zeroTime time.Time
	lt := time.Time(t)
	if lt.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return lt, nil
}

type ModelTime struct {
	CreatedAt Time           `json:"createdAt"`
	UpdatedAt Time           `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" sql:"index"`
}

type ModelAuth struct {
	UpdatedBy string `json:"updatedBy" gorm:"comment:更新者;"`
	CreatedBy string `json:"createdBy" gorm:"comment:创建者;"`
}

// 去掉parseTime&loc=Local
type CreatedAt struct {
	CreatedAt string `json:"createdAt" gorm:"type:datetime;"`
}

func (CreatedAt) BeforeCreate(tx *gorm.DB) error {
	tx.Statement.SetColumn("CreatedAt", time.Now().Format(TimeFormat))
	return nil
}

type UpdatedAt struct {
	UpdatedAt string `json:"UpdatedAt" gorm:"type:datetime;"`
}

func (UpdatedAt) BeforeUpdate(tx *gorm.DB) error {
	tx.Statement.SetColumn("UpdatedAt", time.Now().Format(TimeFormat))
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
