package model

import "github.com/astaxie/beego/orm"

type User struct {
	Id           uint   `orm:"auto" json:"id"`
	UserName     string `orm:"size(32);index" json:"userName"`
	Email        string `orm:"size(64)" json:"email"`
	PasswordHash string `orm:"size(64)" json:"-"`
	CreateTime   uint64 `json:"createTime"`
	UpdateTime   uint64 `json:"updateTime"`
}

func init() {
	orm.RegisterModel(&User{})
}
