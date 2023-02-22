package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"tiktokbackend/models"
	"tiktokbackend/service"
)

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User service.FmtUser `json:"user"`
}

// Register douyin/user/register/ POST
// 用户注册
func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	usi := service.UserServiceImpl{}

	// 查询该用户名是否被注册过
	u := usi.GetUserByUsername(username)
	if username == u.Username { // 如果注册过
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User already exist"},
		})
	} else { // 如果没有注册过，就注册一个
		newUser := models.User{
			Username: username,
			Password: service.EnCoder(password),
		}
		if usi.InsertUser(&newUser) != true {
			println("Insert Data Fail")
		}
		u := usi.GetUserByUsername(username)
		token := service.GenerateToken(username)
		log.Println("注册返回的id: ", u.Id)
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   u.Id,
			Token:    token,
		})
	}
}

// Login POST /douyin/user/login/
// 用户登录
func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	encoderPassword := service.EnCoder(password)
	println("encoderPassword: ", encoderPassword)

	usi := service.UserServiceImpl{}

	u := usi.GetUserByUsername(username)

	if encoderPassword == u.Password {
		fmt.Println("密码相同")
		token := service.GenerateToken(username)
		fmt.Println("成功生成jwt令牌")
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   u.Id,
			Token:    token,
		})
		fmt.Println("成功发送消息")
	} else {
		fmt.Println("密码不同")
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "Username or Password Error"},
		})
	}
}

// UserInfo GET /douyin/user/
// 获取用户信息
func UserInfo(c *gin.Context) {
	userId := c.Query("user_id")
	id, _ := strconv.ParseInt(userId, 10, 64)

	userService := service.UserServiceImpl{}

	if fmtUser, err := userService.GetFmtUserById(id); err != nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User Doesn't Exist"},
		})
	} else {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 0},
			User:     fmtUser,
		})
	}
}
