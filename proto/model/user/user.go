package user

import (
	"microsvc/consts"
	"time"
)

type Userbase struct {
	Uid     int64      `gorm:"column:uid" json:"uid"`           // 内部id
	AliasId int64      `gorm:"column:alias_id" json:"alias_id"` // 可做靓号id/外部id，若不需要可不设置
	Nick    string     `gorm:"column:nick" json:"nick"`
	Age     int32      `gorm:"column:age" json:"age"`
	Sex     consts.Sex `gorm:"column:sex" json:"sex"`
}

type User struct {
	Userbase
	Salt      string    `gorm:"column:salt" json:"salt"`
	Password  string    `gorm:"column:password" json:"password"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (User) TableName() string {
	return "user"
}
