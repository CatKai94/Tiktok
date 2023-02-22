package service

import "tiktokbackend/models"

type UserService interface {
	// GetUserList 获得全部User对象
	GetUserList() []models.User
	// GetUserByUsername 根据username获得User对象
	GetUserByUsername(name string) models.User
	// InsertUser 将user插入表内
	InsertUser(user *models.User) bool
	// GetFmtUserById 未登录情况下,根据user_id获得User对象
	GetFmtUserById(id int64) (FmtUser, error)
	// GetFmtUserByIdWithCurId 已登录(curID)情况下,根据user_id获得User对象
	GetFmtUserByIdWithCurId(id int64, curId int64) (FmtUser, error)
}

// FmtUser 最终封装后,controller返回的FmtUser结构体
type FmtUser struct {
	Id              int64  `json:"id,omitempty"`
	Name            string `json:"name,omitempty"`
	FollowCount     int64  `json:"follow_count"`
	FollowerCount   int64  `json:"follower_count"`
	IsFollow        bool   `json:"is_follow"`
	Avatar          string `json:"avatar,omitempty"`
	BackgroundImage string `json:"background_image,omitempty"`
	Signature       string `json:"signature,omitempty"`
	TotalFavorited  int64  `json:"total_favorited,omitempty"`
	WorkCount       int64  `json:"work_count,omitempty"`
	FavoriteCount   int64  `json:"favorite_count,omitempty"`
}
