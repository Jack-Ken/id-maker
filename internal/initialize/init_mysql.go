package initialize

import (
	"fmt"
	"gorm.io/gorm/schema"
	"id-maker/config"
	"id-maker/internal/entity"

	"go.uber.org/zap"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var SqlSession *gorm.DB

func Init_Mysql(cfg *config.MySqlConfig) (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DbName,
		cfg.Charset,
	)
	SqlSession, err = gorm.Open(mysql.New(mysql.Config{DSN: dsn}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		//PrepareStmt:            true, //缓存预编译命令
		SkipDefaultTransaction:                   true, //禁用默认事务操作
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		zap.L().Error("Connect DB failed", zap.Error(err))
		return
	}
	sqlDb, err := SqlSession.DB()

	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDb.SetMaxIdleConns(cfg.MaxIdleConnections)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDb.SetMaxOpenConns(cfg.MaxOpenConnections)

	err = SqlSession.AutoMigrate(&entity.Segments{})
	if err != nil {
		zap.L().Error("Create Mysql tables failed.Please check the config files....", zap.Error(err))
		return
	}
	return
}
