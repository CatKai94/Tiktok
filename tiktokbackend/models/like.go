package models

import (
	"errors"
	"log"
)

// Like 表的结构。
type Like struct {
	Id      int64 //自增主键
	UserId  int64 //点赞用户id
	VideoId int64 //视频id
	Cancel  int8  //是否点赞，1为点赞，2为取消赞
}

// TableName 修改表名映射
func (Like) TableName() string {
	return "likes"
}

// AddLike 添加点赞数据
func AddLike(userId int64, videoId int64) error {
	like := Like{
		UserId:  userId,
		VideoId: videoId,
		Cancel:  1,
	}
	err := DB.Create(&like).Error
	if err != nil {
		log.Println(err.Error())
		return errors.New("添加点赞数据失败")
	}
	return nil
}

// QueryLikeInfo 根据usrId,videoId查询具体的一条点赞信息
func QueryLikeInfo(userId int64, videoId int64) (Like, error) {
	var result Like
	err := DB.Model(Like{}).Where(map[string]interface{}{
		"user_id":  userId,
		"video_id": videoId,
	}).First(&result).Error
	if err != nil {
		//数据库中没有该数据
		if err.Error() == "record not found" {
			log.Println("数据库中没有该数据")
			return Like{}, nil
		} else {
			log.Println(err.Error())
			return result, errors.New("查询失败")
		}
	}
	return result, nil
}

// GetVideoIdList GetVideoList 根据userId,查询该用户点赞过的全部videoId
func GetVideoIdList(userId int64) ([]int64, error) {
	var result []int64
	err := DB.Model(Like{}).Where(map[string]interface{}{
		"user_id": userId,
		"cancel":  1,
	}).Pluck("video_id", &result).Error
	if err != nil {
		if err.Error() == "record not found" {
			log.Println("该用户没有点赞过任何视频")
			return result, nil
		} else {
			log.Println(err.Error())
			return result, errors.New("查询失败")
		}
	}
	return result, nil
}

// GetUserIdList 根据videoId查询所有点赞过该视频的用户userId
func GetUserIdList(videoId int64) ([]int64, error) {
	var result []int64
	err := DB.Model(Like{}).Where(map[string]interface{}{
		"video_id": videoId,
		"cancel":   1,
	}).Pluck("user_id", &result).Error
	if err != nil {
		if err.Error() == "record not found" {
			log.Println("该视频没有被任何用户点赞过")
			return result, nil
		} else {
			log.Println(err.Error())
			return result, errors.New("查询失败")
		}
	}
	return result, nil
}

// UpdateLikeAction 根据userId,videoId,actionType修改用户点赞的状态
func UpdateLikeAction(userId int64, videoId int64, actionType int32) error {
	err := DB.Model(Like{}).Where(map[string]interface{}{
		"user_id":  userId,
		"video_id": videoId,
	}).Update("cancel", actionType).Error
	if err != nil {
		log.Println(err.Error())
		return errors.New("更新失败")
	}
	return nil
}
