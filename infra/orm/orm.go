package orm

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"microsvc/deploy"
)

var instMap = make(map[deploy.DBname]*gorm.DB)

func InitGorm(must bool) func(func(must bool, err error)) {
	return func(onEnd func(must bool, err error)) {

		gconf := &gorm.Config{
			Logger:          logger.Default.LogMode(logger.Info),
			CreateBatchSize: 100, // 批量插入时，分批进行
		}
		var db *gorm.DB
		var err error
		if len(deploy.XConf.Mysql) == 0 {
			fmt.Println("### there is no mysql config found")
		} else {
			for _, v := range deploy.XConf.Mysql {
				db, err = gorm.Open(mysql.Open(v.Dsn()), gconf)
				if err != nil {
					fmt.Printf("\n****** failed to connect to mysql: err:%v\n", err)
					fmt.Printf("****** mysql.dsn: %s\n\n", v.Dsn())
					break
				}
				instMap[v.DBname] = db
			}
		}
		onEnd(must, err)
	}
}

func GetMysqlInstance(dbname string) *gorm.DB {
	return instMap[deploy.DBname(dbname)]
}
