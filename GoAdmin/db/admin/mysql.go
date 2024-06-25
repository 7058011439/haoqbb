package admin

import (
	"fmt"
	"github.com/7058011439/haoqbb/GoAdmin/config"
	"github.com/7058011439/haoqbb/GoAdmin/db/common"
	hLog "github.com/7058011439/haoqbb/Log"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)

func init() {
	dbName := config.MysqlDBName()
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=True&loc=Local", config.MysqlUserName(), config.MysqlPassWord(), config.MysqlHost(), config.MysqlPort(), dbName)
	if mysql, err := gorm.Open("mysql", connStr); err != nil {
		hLog.ErrorLog("连接mysql数据库错误, err = %v", err)
	} else {
		mysql.LogMode(config.MysqlLog())
		mysql.DB().SetMaxOpenConns(20)
		mysql.DB().SetMaxIdleConns(10)
		mysql.DB().SetConnMaxLifetime(time.Second * 300)
		if err = mysql.DB().Ping(); err != nil {
			panic(err)
		}
		// todo
		// mysql.Exec("ALTER DATABASE " + dbName + " CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci")
		if config.MysqlInit() {
			initModel(mysql)
		}
		dbManager = common.NewManager(mysql)
	}
}

// 初始化表结构以及定义
func initModel(mysql *gorm.DB) {
	var tabList []interface{}
	tabList = append(tabList, &Api{})
	tabList = append(tabList, &Config{})
	tabList = append(tabList, &User{})
	tabList = append(tabList, &Dept{})
	tabList = append(tabList, &DictData{})
	tabList = append(tabList, &DictType{})
	tabList = append(tabList, &LoginLog{})
	tabList = append(tabList, &Menu{})
	tabList = append(tabList, &OperateLog{})
	tabList = append(tabList, &Post{})
	tabList = append(tabList, &Role{})
	if err := mysql.AutoMigrate(tabList...).Error; err != nil {
		hLog.ErrorLog("初始化表失败，err = %v", err)
	}
}
