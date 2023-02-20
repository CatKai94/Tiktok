package service

import "tiktokbackend/models"

// CommentInfo CommentService 接口定义
// 发表评论-使用的结构体-service层引用dao层↑的Comment。
type CommentInfo struct {
	CommentId   int64   `json:"id"`
	UserInfo    FmtUser `json:"user"`
	Content     string  `json:"content"`
	PublishDate string  `json:"create_date"`
}

type CommentService interface {
	// CountFromVideoId 根据videoId获取视频评论数量的接口
	CountFromVideoId(videoId int64) (int64, error)
	// SendComment 发送评论
	SendComment(comment models.Comment) (CommentInfo, error)
	// DeleteComment 删除评论
	DeleteComment(commentId int64) error
	// GetCommentsList 获取视频的所有评论
	GetCommentsList(videoId int64, userId int64) ([]CommentInfo, error)
}
