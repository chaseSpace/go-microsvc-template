package dao

import (
	"microsvc/infra/orm"
	"microsvc/proto/model/user"
)

func GetUser(uid ...int64) (list []*user.User, row user.User, err error) {
	if len(uid) == 1 {
		err = user.Q.Take(&row, "uid=?", uid[0]).Error
	} else {
		err = user.Q.Find(&list, "uid in (?)", uid).Error
	}
	err = orm.IgnoreNil(err)
	return
}
