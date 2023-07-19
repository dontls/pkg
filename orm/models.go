package orm

import (
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
