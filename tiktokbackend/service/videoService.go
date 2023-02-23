package service

import (
	"mime/multipart"
	"tiktokbackend/models"
	"time"
)

type FmtVideo struct {
	models.Video
	Author        FmtUser `json:"author"`
	FavoriteCount int64   `json:"favorite_count"`
	CommentCount  int64   `json:"comment_count"`
	IsFavorite    bool    `json:"is_favorite"`
}

type VideoService interface {
	// Feed 通过传入时间戳，当前用户的id，返回对应的视频切片数组，以及视频数组中最早的发布时间
	Feed(lastTime time.Time, userId int64) ([]FmtVideo, time.Time, error)
	// GetVideo 传入视频id获得对应的视频对象
	GetVideo(videoId int64, userId int64) (FmtVideo, error)
	// PublishAction 将传入的视频流保存在文件服务器中，并存储在mysql表中
	PublishAction(data *multipart.FileHeader, userId int64, title string) error
	// GetVideoList 通过userId来查询对应用户发布的视频，并返回对应的视频切片数组
	GetVideoList(userId int64, curId int64) ([]FmtVideo, error)
}
