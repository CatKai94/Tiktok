package models

import (
	"log"
	"time"
)

type Message struct {
	Id         int64  `json:"id" gorm:"id"`
	UserId     int64  `json:"user_id" gorm:"user_id"`
	ReceiverId int64  `json:"receiver_id" gorm:"receiver_id"`
	MsgContent string `json:"msg_content" gorm:"msg_content"`
	//CreatedAt  time.Time `json:"created_at" gorm:"created_at"`
	CreatedAt int64 `json:"created_at" gorm:"created_at"`
	HaveGet   int64 `gorm:"have_get"`
}

func (Message) TableName() string {
	return "message"
}

// SendMessage fromUserId 发送消息 content 给 toUserId
func SendMessage(fromUserId int64, toUserId int64, content string) error {
	var message Message
	message.UserId = fromUserId
	message.ReceiverId = toUserId
	message.MsgContent = content
	message.CreatedAt = time.Now().Unix()
	message.HaveGet = 0 //如果使用的话，0代表未被客户端拉取过，1代表被拉取过
	return SaveMessage(message)
}
func SaveMessage(msg Message) error {
	if err := DB.Save(&msg).Error; err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

// MessageChat 当前登录用户和其他指定用户的聊天记录
func MessageChat(loginUserId int64, targetUserId int64) ([]Message, error) {
	messages := make([]Message, 0, 5)
	result := DB.Where(&Message{UserId: loginUserId, ReceiverId: targetUserId}).
		Or(&Message{UserId: targetUserId, ReceiverId: loginUserId}).
		Order("created_at asc").
		Find(&messages)
	if result.Error != nil {
		log.Println("获取聊天记录失败！")
		return nil, result.Error
	}
	return messages, nil
}

// DeleteChat 查完就删除记录，防止轮询的重复显示，但就不满足friend list里最新消息的实现了
func DeleteChat(loginUserId int64, targetUserId int64) {
	result := DB.Delete(&Message{UserId: loginUserId, ReceiverId: targetUserId}).
		Or(&Message{UserId: targetUserId, ReceiverId: loginUserId})
	if result.Error != nil {
		log.Println("删除记录失败！")
	}
}

// LatestMessage 返回 loginUserId 和 targetUserId 最近的一条聊天记录
func LatestMessage(loginUserId int64, targetUserId int64) (Message, error) {
	var message Message
	result := DB.Where(&Message{UserId: loginUserId, ReceiverId: targetUserId}).
		Or(&Message{UserId: targetUserId, ReceiverId: loginUserId}).
		Order("created_at desc").Limit(1).Take(&message)
	if result.Error != nil {
		log.Println("获取最近一条聊天记录失败！")
		return Message{}, result.Error
	}
	return message, nil
}
