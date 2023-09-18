package user

import (
	"fmt"
	"microsvc/enums"
	"microsvc/protocol/svc/user"
	"strings"
	"time"
)

type User struct {
	Base
	ExtUid    int64     `gorm:"column:ext_uid" json:"ext_uid"` // 外部id
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

type Base struct {
	Uid        int64     `gorm:"column:uid" json:"uid"` // 内部id
	Nickname   string    `gorm:"column:nickname" json:"nickname"`
	Birthday   time.Time `gorm:"column:birthday" json:"birthday"`
	Sex        enums.Sex `gorm:"column:sex" json:"sex"`
	PasswdSalt string    `gorm:"column:passwd_salt" json:"passwd_salt"`
	Password   string    `gorm:"column:password" json:"password"`
}

func (u *User) TableName() string {
	return "user"
}

func (u *User) Check() error {
	if !(u.Uid > 0 && u.ExtUid > 0) {
		return fmt.Errorf("invalid uid or ext_uid")
	}
	if strings.TrimSpace(u.Nickname) == "" {
		return fmt.Errorf("invalid nickname")
	}
	if u.Birthday.IsZero() || !u.Sex.IsValid() {
		return fmt.Errorf("invalid birthday or sex")
	}
	if !(u.PasswdSalt != "" && u.Password != "") {
		return fmt.Errorf("invalid password")
	}
	return nil
}

func (u *User) Age() int32 {
	today := time.Now()
	age := today.Year() - u.Birthday.Year()
	if today.YearDay() < u.Birthday.YearDay() {
		age--
	}
	return int32(age)
}

func (u *User) ToPb() *user.User {
	return &user.User{
		Uid:      u.Uid,
		Nickname: u.Nickname,
		Age:      u.Age(),
		Sex:      u.Sex.Int32(),
	}
}

func (u *User) SetIntField(uid, extUid int64, sex enums.Sex) {
	u.Uid = uid
	u.ExtUid = extUid
	u.Sex = sex
}

func (u *User) SetStrField(nickname string) {
	u.Nickname = nickname
}

func (u *User) SetTimeField(birthday time.Time, password, salt string) {
	u.Birthday = birthday
	u.Password = password
	u.PasswdSalt = salt
}
