package models

// Like 表的结构。
type Like struct {
	Id      int64 //自增主键
	UserId  int64 //点赞用户id
	VideoId int64 //视频id
	Cancel  int8  //是否点赞，0为点赞，1为取消赞
}

// TableName 修改表名映射
func (Like) TableName() string {
	return "likes"
}
