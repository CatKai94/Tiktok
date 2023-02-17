package controller

type User struct {
	Id            int64  `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	FollowCount   int64  `json:"follow_count,omitempty"`
	FollowerCount int64  `json:"follower_count,omitempty"`
	IsFollow      bool   `json:"is_follow,omitempty"`
}

//	var DemoUser = User{
//		Id:            20052,
//		Name:          "20052",
//		FollowCount:   0,
//		FollowerCount: 0,
//		IsFollow:      false,
//	}
var DemoUser = User{
	Id:            20000,
	Name:          "20000",
	FollowCount:   0,
	FollowerCount: 0,
	IsFollow:      false,
}
