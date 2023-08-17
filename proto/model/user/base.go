package user

import (
	"fmt"
	"gorm.io/gorm"
	"microsvc/infra/orm"
)

const dbName = "microsvc"

func MustReady() {
	if Q() == nil {
		panic(fmt.Sprintf("dbname [%s] is not found in instance map", dbName))
	}
}

func Q() *gorm.DB {
	return orm.GetMysqlInstance(dbName)
}
