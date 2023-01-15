package models

import "log"

// 映射字段名
type User struct {
	ID         int64
	Username   string `gorm:"column:username"`
	Password   string `gorm:"column:password"`
	CreateTime int64  `gorm:"column:createtime"`
}

// 映射表名
func (u User) TableName() string {
	return "user"
}

// 增加
func SaveUser(user *User) {
	// 数据库操作
	err := DB.Create(user).Error
	if err != nil {
		log.Println("insert user ", err)
	}

}

// 查询
func GetById(id int64) User {
	var user User
	err := DB.Where("id=?", id).First(&user).Error
	if err != nil {
		log.Println("insert user error ", err)
	}

	return user
}

// 更新
func UpdateUser(id int64) {
	err := DB.Model(&User{}).Where("id=?", id).Update("username", "lisi")
	if err != nil {
		log.Println("uopdate user by id is error ", err)
	}
}

// 删除
func DeleteUser(id int64) {
	err := DB.Where("id=?", id).Delete(&User{})
	if err != nil {
		log.Println("delete user by id is error ", err)
	}
}
