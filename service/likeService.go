package service

// LikeService 定义点赞状态和点赞数量
type LikeService interface {
	// GetUserIsLike 根据当前视频id判断当前用户是否点赞了该视频。
	JudgeUserIsLike(videoId int64, userId int64) (bool, error)

	// GetLikeCount  根据当前视频id获取当前视频点赞数量。
	GetLikeCount(videoId int64) (int64, error)

	// GetUserTotalIsLikedCount 根据userId获取这个用户总共被点赞数量
	GetUserTotalIsLikedCount(userId int64) (int64, error)

	// GetLikeVideoCount  根据userId获取这个用户点赞视频数量
	GetLikeVideoCount(userId int64) (int64, error)

	// LikeAction 当前用户对视频的点赞操作 ,并把这个行为更新到like表中。
	//当前操作行为，1点赞，2取消点赞。
	LikeAction(userId int64, videoId int64, actionType int32) error
	// GetLikeVideoList 获取当前用户的所有点赞视频，调用videoService的方法
	GetLikeVideoList(userId int64, curId int64) ([]FmtVideo, error)
}
