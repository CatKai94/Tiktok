package models

import (
	"log"
	"tiktokbackend/config"
	"time"
)

// 映射字段名
type Video struct {
	Id          int64 `json:"id"`
	AuthorId    int64
	PlayUrl     string `json:"play_url"`
	CoverUrl    string `json:"cover_url"`
	PublishTime time.Time
	Title       string `json:"title"`
}

// 表名映射
func (Video) TableName() string {
	return "videos"
}

// GetVideosByLastTime
// 依据一个时间，来获取这个时间之前的一些视频
func GetVideosByLastTime(lastTime time.Time) ([]Video, error) {
	videos := make([]Video, config.VideoCount)
	result := DB.Where("publish_time<?", lastTime).Order("publish_time desc").Limit(config.VideoCount).Find(&videos)
	if result.Error != nil {
		log.Println("查询videos有错误", result.Error)
		return videos, result.Error
	}
	log.Println("查询到视频：", result)
	log.Println("查询到videos：", videos[0])
	return videos, nil
}

// GetVideosByAuthorId
// 根据作者的id来查询对应数据库数据，返回Video切片
func GetVideosByAuthorId(authorId int64) ([]Video, error) {
	var data []Video
	result := DB.Where(&Video{AuthorId: authorId}).Find(&data)
	if result.Error != nil {
		return nil, result.Error
	}
	return data, nil
}

// GetVideoByVideoId
// 依据VideoId来获得视频信息
func GetVideoByVideoId(videoId int64) (Video, error) {
	var tableVideo Video
	tableVideo.Id = videoId
	result := DB.First(&tableVideo)
	if result.Error != nil {
		return tableVideo, result.Error
	}
	return tableVideo, nil
}

// Save 保存视频记录
func Save(videoName string, imageName string, authorId int64, title string) error {
	VideoUrl := config.IpUrl + "/static/videos/" //视频的根路由
	CoverUrl := config.IpUrl + "/static/images/" //封面截图的根路由
	newVideo := Video{
		PublishTime: time.Now(),
		PlayUrl:     VideoUrl + videoName + ".mp4",
		CoverUrl:    CoverUrl + imageName + ".jpg",
		AuthorId:    authorId,
		Title:       title,
	}
	result := DB.Create(&newVideo)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// GetVideoIdsByAuthorId
// 通过作者id来查询发布的视频id切片集合
func GetVideoIdsByAuthorId(authorId int64) ([]int64, error) {
	var id []int64
	//通过pluck来获得单独的切片
	result := DB.Model(&Video{}).Where("author_id", authorId).Pluck("id", &id)
	//如果出现问题，返回对应到空，并且返回error
	if result.Error != nil {
		return nil, result.Error
	}
	return id, nil
}
