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
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

type Base struct {
	Uid        int64     `gorm:"column:uid" json:"uid"` // 内部id
	Nickname   string    `gorm:"column:nickname" json:"nickname"`
	Birthday   time.Time `gorm:"column:birthday" json:"birthday"`
	Sex        enums.Sex `gorm:"column:sex" json:"sex"`
	PasswdSalt string    `gorm:"column:password_salt" json:"password_salt"`
	Password   string    `gorm:"column:password" json:"password"`
	Phone      string    `gorm:"column:phone" json:"phone"`
}

func (u *User) TableName() string {
	return "user"
}

func (u *User) Check() error {
	if !(u.Uid > 0) {
		return fmt.Errorf("invalid uid")
	}
	u.Nickname = strings.TrimSpace(u.Nickname)
	if u.Nickname == "" || len([]rune(u.Nickname)) > 10 {
		return fmt.Errorf("无效昵称或超出长度")
	}
	if u.Sex.IsInvalid() {
		return fmt.Errorf("请设置有效的性别")
	}
	if u.Birthday.IsZero() {
		return fmt.Errorf("无效的生日信息")
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
		Birthday: u.Birthday.Format(time.DateOnly),
		Sex:      u.Sex.Int32(),
	}
}

func (u *User) SetIntField(uid int64, sex enums.Sex) {
	u.Uid = uid
	u.Sex = sex
}

func (u *User) SetStrField(nickname, password, salt string) {
	u.Nickname = nickname
	u.Password = password
	u.PasswdSalt = salt
}

func (u *User) SetTimeField(birthday time.Time) {
	u.Birthday = birthday
}
