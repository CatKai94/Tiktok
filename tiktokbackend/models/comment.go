package models

import (
	"errors"
	"log"
	"time"
)

// Comment
// 评论信息-数据库中的结构体-dao层使用
type Comment struct {
	Id          int64     // 评论id
	UserId      int64     // 评论用户id
	VideoId     int64     // 视频id
	CommentText string    // 评论内容
	CreateDate  time.Time // 评论发布的日期mm-dd
	Cancel      int32     // 取消评论为1，发布评论为0
}

// TableName 修改表名映射
func (Comment) TableName() string {
	return "comments"
}

// Count
// 1、使用video id 查询Comment数量
func Count(videoId int64) (int64, error) {
	var count int64
	// 数据库中查询评论数量
	err := DB.Model(Comment{}).Where(map[string]interface{}{"video_id": videoId, "cancel": 0}).Count(&count).Error
	if err != nil {
		return -1, errors.New("查询评论数量失败")
	}
	return count, nil
}

// CommentIdList 根据视频id获取评论id 列表
func CommentIdList(videoId int64) ([]string, error) {
	var commentIdList []string
	err := DB.Model(Comment{}).Select("id").Where("video_id = ?", videoId).Find(&commentIdList).Error
	if err != nil {
		log.Println("CommentIdList:", err)
		return nil, err
	}
	return commentIdList, nil
}

// InsertComment
// 2、发表评论
func InsertComment(comment Comment) (Comment, error) {
	// 数据库中插入一条评论信息
	err := DB.Model(Comment{}).Create(&comment).Error
	if err != nil {
		return Comment{}, errors.New("发表评论存入数据失败")
	}
	return comment, nil
}

// DeleteComment
// 3、删除评论，传入评论id
func DeleteComment(id int64) error {
	var commentInfo Comment
	// 先查询是否有此评论
	result := DB.Model(Comment{}).Where(map[string]interface{}{"id": id, "cancel": 0}).First(&commentInfo)
	if result.RowsAffected == 0 { // 查询到此评论数量为0则返回无此评论
		return errors.New("此评论不存在")
	}
	// 数据库中删除评论-更新评论状态为-1
	err := DB.Model(Comment{}).Where("id = ?", id).Update("cancel", 1).Error
	if err != nil {
		return errors.New("删除评论失败")
	}
	return nil
}

// GetCommentList
// 4.根据视频id查询所属评论全部列表信息
func GetCommentList(videoId int64) ([]Comment, error) {
	// 数据库中查询评论信息list
	var commentList []Comment
	result := DB.Model(Comment{}).Where(map[string]interface{}{"video_id": videoId, "cancel": 0}).
		Order("create_date desc").Find(&commentList)
	// 若此视频没有评论信息，返回空列表
	if result.RowsAffected == 0 {
		log.Println("此视频赞无评论") // 函数返回提示无评论
		return nil, nil
	}
	// 若获取评论列表出错
	if result.Error != nil {
		log.Println(result.Error.Error())
		return commentList, errors.New("获取评论列表失败")
	}
	return commentList, nil
}
