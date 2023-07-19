package orm

import (
	"errors"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var gconf = gorm.Config{
	NamingStrategy: schema.NamingStrategy{
		SingularTable: true,
		NoLowerCase:   true,
	},
	DisableForeignKeyConstraintWhenMigrating: true, // 禁用外键约束
}

type Options struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

var _models []interface{}

func RegisterModel(dst ...interface{}) {
	_models = append(_models, dst...)
}

var _engineModels = make(map[string][]interface{})

func RegisterEngineModel(engine string, dst ...interface{}) {
	v := _engineModels[engine]
	v = append(v, dst...)
	_engineModels[engine] = v
}

func NewDB(o *Options) (db *gorm.DB, err error) {
	switch o.Name {
	case "mysql":
		db, err = gorm.Open(mysql.New(mysql.Config{
			DSN: o.Address,
			// DefaultStringSize:         64,    // string 类型字段的默认长度
			DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
			DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
			DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
			SkipInitializeWithVersion: false, // 根据版本自动配置
		}), &gconf)
	case "sqlite3":
		// db, err = gorm.Open(sqlite.Open(address), &gconf)
	case "postgresql":
		// db, err = gorm.Open(postgres.Open(address), &gconf)
	default:
	}
	if err != nil {
		return nil, err
	}
	if db == nil {
		return nil, errors.New("db invalid")
	}
	sqldb, err := db.DB()
	if err != nil {
		return nil, err
	}
	// SetMaxIdleCons 设置连接池中的最大闲置连接数。
	sqldb.SetMaxIdleConns(10)
	// SetMaxOpenCons 设置数据库的最大连接数量。
	sqldb.SetMaxOpenConns(100)
	// auto migrate
	db.AutoMigrate(_models...)
	for k, v := range _engineModels {
		db.Set("gorm:table_options", "ENGINE="+k).AutoMigrate(v...)
	}
	_db = db
	return db, nil
}
