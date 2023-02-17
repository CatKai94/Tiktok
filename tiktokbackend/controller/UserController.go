package controller

import (
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

// douyin/user/register/ POST
// 用户注册
func Register(c *gin.Context) {
	//jl
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
			log.Println("Insert Data Fail")
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
	//jl
	username := c.Query("username")
	password := c.Query("password")
	encoderPassword := service.EnCoder(password)
	//println("encoderPassword: ", encoderPassword)
	usi := service.UserServiceImpl{}
	u := usi.GetUserByUsername(username)
	if username != u.Username {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist..."},
		})
	} else {
		if encoderPassword == u.Password {
			//fmt.Println("密码相同")
			token := service.GenerateToken(username)
			//fmt.Println("成功生成jwt令牌")
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: Response{StatusCode: 0},
				UserId:   u.Id,
				Token:    token,
			})
			//fmt.Println("成功发送消息")
		} else {
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: Response{
					StatusCode: 1,
					StatusMsg:  "Password Error",
				},
			})
		}
	}
}

// UserInfo GET /douyin/user/
// 获取用户信息
func UserInfo(c *gin.Context) {
	user_id := c.Query("user_id")
	id, _ := strconv.ParseInt(user_id, 10, 64)

	usi := service.UserServiceImpl{}
	//usi := service.UserServiceImpl{  //这部分等把点赞，关注，喜欢，视频部分写完再来完善
	//	FollowService: &service.FollowServiceImp{},
	//	LikeService:   &service.LikeServiceImpl{},
	//}

	if fmtUser, err := usi.GetFmtUserById(id); err != nil {
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
