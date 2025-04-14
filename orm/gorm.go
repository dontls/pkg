package orm

import (
	"reflect"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var _db *gorm.DB

// H 多列处理
type H map[string]any

// gorm对象
func DB() *gorm.DB {
	return _db
}

func DbModel(v any) *gorm.DB {
	return _db.Model(v)
}

// DbCount 数目
func DbCount(model, where any) int64 {
	var count int64
	db := _db.Model(model)
	if where != nil {
		db = db.Where(where)
	}
	db.Count(&count)
	return count
}

// DbCreate 创建
func DbCreate(model any) error {
	return _db.Create(model).Error
}

// DbSave 保存
func DbSave(value any) error {
	return _db.Save(value).Error
}

// DbUpdateModel 更新
// 默认匹配primary_key
func DbUpdateModel(model any) error {
	return _db.Model(model).Updates(model).Error
}

// DbUpdateModelBy 条件更新
func DbUpdateModelBy(model any, where string, args ...any) error {
	return _db.Where(where, args...).Updates(model).Error
}

// DbUpdateModelByID 更新
func DbUpdateByID(model, id any) error {
	return _db.Where("ID = ?", id).Updates(model).Error
}

// DbUpdateFields 更新指定列，
// 默认匹配primary_key
func DbUpdateFields(model any, fields ...string) error {
	return _db.Model(model).Select(fields).Updates(model).Error
}

// DbUpdateByID 更新
// 如果id是数组，则批量更新, 这是不使用H, gorm不识别
func DbUpdateValues(model any, ids any, value map[string]any) error {
	return _db.Model(model).Where("ID in (?)", ids).Updates(value).Error
}

// DbDeletes 批量删除
func DbDeletes(value any) error {
	return _db.Delete(value).Error
}

// DbDeleteByIds 批量删除 id数组[]
func DbDeleteByIDs(model, ids any) error {
	return _db.Delete(model, ids).Error
}

// DbDeleteBy 删除
func DbDeleteBy(model any, where string, args ...any) (count int64, err error) {
	db := _db.Where(where, args...).Delete(model)
	err = db.Error
	if err != nil {
		return
	}
	count = db.RowsAffected
	return
}

// DbFirstBy 指定条件查找
func DbFirstBy(out any, where string, args ...any) (err error) {
	err = _db.Where(where, args...).First(out).Error
	return
}

// DbFirstByID 查找
func DbFirstByID(out any, id any) error {
	return _db.First(out, id).Error
}

// DbFirstWhere 查找
func DbFirstWhere(out, where any) error {
	return _db.Where(where).First(out).Error
}

// DbFind 多个查找
func DbFind(out any, orders ...string) error {
	db := _db
	if len(orders) > 0 {
		for _, order := range orders {
			db = db.Order(order)
		}
	}
	return db.Find(out).Error
}

// DbFindBy 多个条件查找
func DbFindBy(out any, where string, args ...any) (int64, error) {
	db := _db.Where(where, args...).Find(out)
	return db.RowsAffected, db.Error
}

// DbFindPageRaw obj必须是数组类型
func DbFindPage(query string, obj any, page, size int) (int64, error) {
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

var _models []any

func RegisterModel(dst ...any) {
	_models = append(_models, dst...)
}

var _engineModels = make(map[string][]any)

func RegisterEngineModel(engine string, dst ...any) {
	v := _engineModels[engine]
	v = append(v, dst...)
	_engineModels[engine] = v
}

var gconf = gorm.Config{
	NamingStrategy: schema.NamingStrategy{
		SingularTable: true,
		NoLowerCase:   true,
	},
	DisableForeignKeyConstraintWhenMigrating: true, // 禁用外键约束
}

// 在NewDB之前调用
func SetLogger(w logger.Writer) {
	gconf.Logger = logger.New(w, logger.Config{
		SlowThreshold:             time.Second, // Slow SQL threshold
		LogLevel:                  logger.Info, // Log level
		IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
		Colorful:                  false,       // Disable color
	})
}

func CreateDB(dialector gorm.Dialector, debug bool) (err error) {
	_db, err = gorm.Open(dialector, &gconf)
	if err != nil {
		return err
	}
	sqldb, err := _db.DB()
	if err != nil {
		return err
	}
	if debug {
		_db = _db.Debug()
	}
	// SetMaxIdleCons 设置连接池中的最大闲置连接数。
	sqldb.SetMaxIdleConns(10)
	// SetMaxOpenCons 设置数据库的最大连接数量。
	sqldb.SetMaxOpenConns(100)
	// auto migrate
	_db.AutoMigrate(_models...)
	for k, v := range _engineModels {
		_db.Set("gorm:table_options", "ENGINE="+k).AutoMigrate(v...)
	}
	return nil
}
