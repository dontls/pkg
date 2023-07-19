package orm

import (
	"reflect"

	"gorm.io/gorm"
)

var _db *gorm.DB

// H 多列处理
type H map[string]interface{}

// SetDB gorm对象
func SetDB(db *gorm.DB) {
	_db = db
}

// SetDB gorm对象
func DB() *gorm.DB {
	return _db
}

func Model(v interface{}) *gorm.DB {
	return _db.Model(v)
}

// DbCount 数目
func DbCount(model, where interface{}) int64 {
	var count int64
	db := _db.Model(model)
	if where != nil {
		db = db.Where(where)
	}
	db.Count(&count)
	return count
}

// DbCreate 创建
func DbCreate(model interface{}) error {
	return _db.Create(model).Error
}

// DbSave 保存
func DbSave(value interface{}) error {
	return _db.Save(value).Error
}

// DbUpdateModel 更新
// 默认匹配primary_key
func DbUpdateModel(model interface{}) error {
	return _db.Model(model).Updates(model).Error
}

// DbUpdateFields 更新指定列，
// 默认匹配primary_key
func DbUpdateFields(model interface{}, fields ...string) error {
	return _db.Model(model).Select(fields).Updates(model).Error
}

// DbUpdateModelBy 条件更新
func DbUpdateModelBy(model interface{}, where string, args ...interface{}) error {
	return _db.Where(where, args...).Updates(model).Error
}

// DbUpdateModelByID 更新
func DbUpdateModelByID(model, id interface{}) error {
	return _db.Where("ID = ?", id).Updates(model).Error
}

// DbUpdateByID 更新
// 如果id是数组，则批量更新
func DbUpdateByID(model interface{}, ids interface{}, value map[string]interface{}) error {
	return _db.Model(model).Where("ID in (?)", ids).Updates(value).Error
}

// DbDeletes 批量删除
func DbDeletes(value interface{}) error {
	return _db.Delete(value).Error
}

// DbDeleteByIds 批量删除
// ids id数组 []
func DbDeleteByIDs(model, ids interface{}) error {
	return _db.Delete(model, ids).Error
}

// DbDeleteBy 删除
func DbDeleteBy(model interface{}, where string, args ...interface{}) (count int64, err error) {
	db := _db.Where(where, args...).Delete(model)
	err = db.Error
	if err != nil {
		return
	}
	count = db.RowsAffected
	return
}

// DbFirstBy 指定条件查找
func DbFirstBy(out interface{}, where string, args ...interface{}) (err error) {
	err = _db.Where(where, args...).First(out).Error
	return
}

// DbFirstByID 查找
func DbFirstByID(out interface{}, id uint) error {
	return _db.First(out, id).Error
}

// DbFirstWhere 查找
func DbFirstWhere(out, where interface{}) error {
	return _db.Where(where).First(out).Error
}

// DbFind 多个查找
func DbFind(out interface{}, orders ...string) error {
	db := _db
	if len(orders) > 0 {
		for _, order := range orders {
			db = db.Order(order)
		}
	}
	return db.Find(out).Error
}

// DbFindBy 多个条件查找
func DbFindBy(out interface{}, where string, args ...interface{}) (int64, error) {
	db := _db.Where(where, args...).Find(out)
	return db.RowsAffected, db.Error
}

// DbFindPageRaw obj必须是数组类型
func DbFindPage(query string, obj interface{}, page, size int) (int64, error) {
	s := reflect.ValueOf(obj)
	if s.Kind() == reflect.Ptr {
		s = s.Elem()
	}
	if s.Kind() != reflect.Slice {
		return 0, nil
	}
	db := _db.Raw(query)
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return 0, err
	}
	if page > 0 {
		db = db.Offset((page - 1) * size).Limit(size)
	}
	return total, db.Scan(obj).Error
}

