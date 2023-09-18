package dao

import (
	"microsvc/pkg/xerr"
	muser "microsvc/proto/model/user"
	"time"
)

func CreateUser(ent *muser.User) error {
	if err := ent.Check(); err != nil {
		return xerr.ErrInvalidRegisterInfo.AppendMsg(err.Error())
	}
	ent.CreatedAt = time.Now()
	ent.UpdatedAt = ent.CreatedAt
	return muser.Q.Create(ent).Error
}
