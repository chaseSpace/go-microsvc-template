package dao

import (
	"gorm.io/gorm"
	"microsvc/infra/orm"
	"microsvc/pkg/xerr"
	"microsvc/proto/model/user"
)

func GetMaxUid() (int64, error) {
	row := new(user.User)
	err := user.Q.Order("uid desc").Take(row).Error
	err = orm.IgnoreNil(err)
	return row.Uid, err
}

func IsUidExists(uid uint64) (bool, error) {
	if uid < 1 {
		return false, xerr.ErrParams.AppendMsg("invalid uid on insert")
	}
	exec := user.Q.Take(&user.User{}, "uid=?", uid)
	if exec.Error != nil && exec.Error != gorm.ErrRecordNotFound {
		return false, exec.Error
	}
	return exec.RowsAffected > 0, nil
}

func GetUser(uid ...int64) (list []*user.User, row user.User, err error) {
	if len(uid) == 1 {
		err = user.Q.Take(&row, "uid=?", uid[0]).Error
	} else {
		err = user.Q.Find(&list, "uid in (?)", uid).Error
	}
	err = orm.IgnoreNil(err)
	return
}

func GetUserByPhone(phone ...string) (list []*user.User, row user.User, err error) {
	if len(phone) == 1 {
		err = user.Q.Take(&row, "phone=?", phone[0]).Error
	} else {
		err = user.Q.Find(&list, "phone in (?)", phone).Error
	}
	err = orm.IgnoreNil(err)
	return
}
