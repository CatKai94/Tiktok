package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"tiktokbackend/service"
)

type ChatResponse struct {
	Response
	MessageList []service.Message `json:"message_list"`
}

// MessageAction
func MessageAction(c *gin.Context) {
	//loginUserId, _ := strconv.ParseInt(c.Query("userId"), 10, 64)
	UserId := c.Value("userId")
	loginUserId, _ := strconv.ParseInt(UserId.(string), 10, 64)
	content := c.Query("content")
	toUserId, _ := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	ActionType, _ := strconv.ParseInt(c.Query("action_type"), 10, 64)

	messageService := service.MessageServiceImpl{}
	err := messageService.ActionMessage(loginUserId, toUserId, content, ActionType)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "MessageAction error"})
	}
	c.JSON(http.StatusOK, Response{StatusCode: 0})
}

// MessageChat
func MessageChat(c *gin.Context) {
	UserId := c.Value("userId")
	loginUserId, _ := strconv.ParseInt(UserId.(string), 10, 64)
	targetUserId, _ := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	messageService := service.MessageServiceImpl{}
	messages, err := messageService.MessageChat(loginUserId, targetUserId)
	log.Println(messages)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	} else {
		c.JSON(http.StatusOK, ChatResponse{Response: Response{StatusCode: 0, StatusMsg: "success"}, MessageList: messages})
	}
}
