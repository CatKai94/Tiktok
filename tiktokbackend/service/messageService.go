package service

import "tiktokbackend/models"

type Message struct {
	Id         int64  `json:"id"`
	UserId     int64  `json:"from_user_id"`
	ReceiverId int64  `json:"to_user_id"`
	MsgContent string `json:"content"`
	CreatedAt  int64  `json:"create_time" db:"created_at"`
}

type LatestMessage struct {
	message string `json:"message"`
	msgType int64  `json:"msg_type"`
}

type MessageService interface {
	// ActionMessage 发送消息，即向数据库中保存消息
	ActionMessage(fromUserId int64, toUserId int64, content string) error

	// MessageChat 用来查询数据库中的消息记录，
	MessageChat(loginUserId int64, targetUserId int64) ([]Message, error)

	// getRespMessage 将数据库Db中的message映射为要返回的message结构体
	getRespMessage(messages *[]Message, plainMessages *[]models.Message) error

	// LatestMessage 返回两个 loginUserId 和好友 targetUserId 最近的一条聊天记录
	LatestMessage(loginUserId int64, targetUserId int64) (LatestMessage, error)
}
