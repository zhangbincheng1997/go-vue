package core

import (
	"fmt"
	"main/global"

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
		global.LOG.Errorf("MySQL连接失败：%v", err)
		return nil
	}
	global.LOG.Infof("MySQL连接成功：%v", dbDSN)
	return db
}
