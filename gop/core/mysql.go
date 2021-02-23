package core

import (
	"fmt"
	"main/global"
	"main/model"
	"os"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// MySQL ...
func MySQL() *gorm.DB {
	cfg := global.CONFIG.MySQL
	dbDSN := fmt.Sprintf("%s:%s@%s(%s)/%s?%s", cfg.Username, cfg.Password, cfg.Protocol, cfg.Addr, cfg.Database, cfg.Config)
	dialector := mysql.New(mysql.Config{
		DSN:                       dbDSN, // data source name, refer https://github.com/go-sql-driver/mysql#dsn-data-source-name
		DefaultStringSize:         256,   // add default size for string fields, by default, will use db type `longtext` for fields without size, not a primary key, no index defined and don't have default values
		DisableDatetimePrecision:  true,  // disable datetime precision support, which not supported before MySQL 5.6
		DontSupportRenameIndex:    true,  // drop & create index when rename index, rename index not supported before MySQL 5.7, MariaDB
		DontSupportRenameColumn:   true,  // use change when rename column, rename rename not supported before MySQL 8, MariaDB
		SkipInitializeWithVersion: false, // smart configure based on used version
	})
	config := &gorm.Config{
		Logger:                                   logger.Default.LogMode(logger.Info),
		DisableForeignKeyConstraintWhenMigrating: true,
	}
	db, err := gorm.Open(dialector, config)
	if err != nil {
		global.LOG.Error("MySQL连接失败", zap.Any("err", err))
		return nil
	}
	global.LOG.Info("MySQL连接成功", zap.Any("dbDSN", dbDSN))
	return db
}

// MySQLTables ...
func MySQLTables(db *gorm.DB) {
	err := db.AutoMigrate(
		model.User{},
		model.Role{},
		model.UserRole{},
	)
	if err != nil {
		global.LOG.Error("创建表失败", zap.Any("err", err))
		os.Exit(0)
	}
	global.LOG.Info("创建表成功")
}
