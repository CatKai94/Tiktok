package models

import "log"

// 映射字段名
type User struct {
	Id       int64
	Username string
	Password string
}

// 表名映射
func (u User) TableName() string {
	return "users"
}

// 增加
func SaveUser(user *User) {
	// 数据库操作
	err := DB.Create(user).Error
	if err != nil {
		log.Println("insert user ", err)
	}
}

// GetUserList 获取全部User对象
func GetUserList() ([]User, error) {
	users := []User{}
	if err := DB.Find(&users).Error; err != nil {
		log.Println(err.Error())
		return users, err
	}
	return users, nil
}

// GetUserByUsername 根据username获得User对象
func GetUserByUsername(name string) (User, error) {
	user := User{}
	if err := DB.Where("username = ?", name).First(&user).Error; err != nil {
		log.Println(err.Error())
		return user, err
	}
	return user, nil
}

// GetUserById 根据user_id获得User对象
func GetUserById(id int64) (User, error) {
	user := User{}
	if err := DB.Where("id = ?", id).First(&user).Error; err != nil {
		log.Println(err.Error())
		return user, err
	}
	return user, nil
}

// InsertUser 将user插入表内
func InsertUser(user *User) bool {
	if err := DB.Create(&user).Error; err != nil {
		log.Println(err.Error())
		return false
	}
	return true
}
