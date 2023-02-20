package service

type followService interface {

	// IsFollowing 根据当前用户id和目标用户id来判断当前用户是否关注了目标用户
	IsFollowing(userId int64, curId int64) (bool, error)

	// GetFollowerCnt 根据用户id查询粉丝数量
	GetFollowerCnt(userId int64) (int64, error)

	// GetFollowingCnt 根据用户id查询该用户的关注数量
	GetFollowingCnt(userId int64) (int64, error)

	// AddFollowRelation 当前用户关注目标用户 操作
	AddFollowRelation(userId int64, curId int64) (bool, error)

	// DeleteFollowRelation 当前用户取消对目标用户的关注 操作
	DeleteFollowRelation(userId int64, curId int64) (bool, error)

	// GetFollowing 获取当前用户的关注列表
	GetFollowing(userId int64) ([]FmtUser, error)

	// GetFollowers 获取当前用户的粉丝列表
	GetFollowers(userId int64) ([]FmtUser, error)
}
