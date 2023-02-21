package service

import (
	"log"
	"tiktokbackend/models"
)

type MessageServiceImpl struct {
}

func (messageService *MessageServiceImpl) ActionMessage(fromUserId int64, toUserId int64, content string, actionType int64) error {
	var err error
	if actionType == 1 {
		err = models.SendMessage(fromUserId, toUserId, content)
	} else {
		log.Println("actionType != 1")
		return err
	}
	return err
}
func (messageService *MessageServiceImpl) MessageChat(loginUserId int64, targetUserId int64) ([]Message, error) {
	messages := make([]Message, 0, 6)
	DbMessages, err := models.MessageChat(loginUserId, targetUserId)
	if err != nil {
		log.Println("MessageChat Service:", err)
		return nil, err
	}
	err = messageService.getRespMessage(&messages, &DbMessages)
	if err != nil {
		log.Println("getRespMessage:", err)
		return nil, err
	}
	return messages, nil
}

// 返回 message list 接口所需的 Message 结构体
func (messageService *MessageServiceImpl) getRespMessage(messages *[]Message, DbMessages *[]models.Message) error {
	for _, DbMessage := range *DbMessages {
		var message Message
		message.Id = DbMessage.Id
		message.ReceiverId = DbMessage.ReceiverId
		message.UserId = DbMessage.UserId
		message.MsgContent = DbMessage.MsgContent
		message.CreatedAt = DbMessage.CreatedAt
		*messages = append(*messages, message)
	}
	return nil
}

func (messageService *MessageServiceImpl) LatestMessage(loginUserId int64, targetUserId int64) (LatestMessage, error) {
	pMessage, err := models.LatestMessage(loginUserId, targetUserId)
	if err != nil {
		log.Println("models.LatestMessage", err)
		return LatestMessage{}, err
	}
	var latestMessage LatestMessage
	latestMessage.message = pMessage.MsgContent
	//0 => 当前请求用户接收的消息， 1 => 当前请求用户发送的消息
	if pMessage.UserId == loginUserId {
		latestMessage.msgType = 1
	} else {
		latestMessage.msgType = 0
	}
	return latestMessage, nil
}
