package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"github.com/dgrijalva/jwt-go"
	"log"
	"strconv"
	"tiktokbackend/config"
	"tiktokbackend/models"
	"time"
)

type UserServiceImpl struct {
}

func (usi *UserServiceImpl) GetUserList() []models.User {
	users, err := models.GetUserList()
	if err != nil {
		log.Println("Err:", err.Error())
		return users
	}
	return users
}

func (usi *UserServiceImpl) GetUserByUsername(username string) models.User {
	user, err := models.GetUserByUsername(username)
	if err != nil {
		log.Println("Err:", err.Error())
		log.Println("User Not Found")
		return user
	}
	return user
}

// GetUserById 根据user_id获得TUser对象
func (usi *UserServiceImpl) GetUserById(id int64) models.User {
	user, err := models.GetUserById(id)
	if err != nil {
		log.Println("Err:", err.Error())
		log.Println("User Not Found")
		return user
	}
	log.Println("Query User Success")
	return user
}

// InsertUser 将tableUser插入表内
func (usi *UserServiceImpl) InsertUser(user *models.User) bool {
	flag := models.InsertUser(user)
	if flag == false {
		log.Println("插入失败")
		return false
	}
	return true
}

// GetFmtUserById 未登录情况下,根据user_id获得User对象
func (usi *UserServiceImpl) GetFmtUserById(id int64) (FmtUser, error) {
	fmtUser := FmtUser{
		Id:             0,
		Name:           "",
		FollowCount:    0,
		FollowerCount:  0,
		IsFollow:       false,
		TotalFavorited: 0,
		FavoriteCount:  0,
	}
	//user, err := models.GetUserById(id)
	//if err != nil {
	//	log.Println("Err:", err.Error())
	//	log.Println("User Not Found")
	//	return fmtUser, err
	//}
	//log.Println("Query User Success")
	//followCount, _ := usi.GetFollowingCnt(id)
	//if err != nil {
	//	log.Println("Err:", err.Error())
	//}
	//followerCount, _ := usi.GetFollowerCnt(id)
	//if err != nil {
	//	log.Println("Err:", err.Error())
	//}
	//u := GetLikeService() //解决循环依赖
	//totalFavorited, _ := u.TotalFavourite(id)
	//favoritedCount, _ := u.FavouriteVideoCount(id)
	//fmtUser = FmtUser{
	//	Id:             id,
	//	Name:           user.Username,
	//	FollowCount:    followCount,
	//	FollowerCount:  followerCount,
	//	IsFollow:       false,
	//	TotalFavorite: totalFavorited,
	//	FavoriteCount:  favoritedCount,
	//}
	return fmtUser, nil
}

// GetFmtUserByIdWithCurId 已登录(curID)情况下,根据user_id获得User对象
func (usi *UserServiceImpl) GetFmtUserByIdWithCurId(id int64, curId int64) (FmtUser, error) {
	var likeService LikeServiceImpl
	var followService FollowServiceImp
	var videoService VideoServiceImpl

	fmtUser := FmtUser{
		Id:              0,
		Name:            "",
		FollowCount:     0,
		FollowerCount:   0,
		IsFollow:        false,
		Avatar:          config.DefaultAvatar,
		BackgroundImage: config.DefaultBGI,
		Signature:       config.DefaultSign,
		TotalFavorited:  0, // 获赞数量
		WorkCount:       0, // 作品数量
		FavoriteCount:   0, // 点赞数量

	}

	user, err := models.GetUserById(id)
	if err != nil {
		log.Println("产生错误:", err.Error())
		log.Println("没有查到用户id为", id, "的用户")
		return fmtUser, err
	}
	fmtUser.Name = user.Username
	fmtUser.Id = user.Id

	// 获取关注的人数
	followCount := followService.GetTotalFollowingsCnt(id)
	fmtUser.FollowCount = followCount

	// 获取粉丝的人数
	followerCount := followService.GetTotalFollowersCnt(id)
	fmtUser.FollowerCount = followerCount

	// 判断是否关注了该用户
	isFollow, err := followService.IsFollowing(id, curId)
	fmtUser.IsFollow = isFollow

	// 获取作品数量
	workCount := videoService.GetVideoCntByUserId(id)
	fmtUser.WorkCount = workCount

	// 获取用户获得的点赞总数和获赞总数
	favoriteCount, _ := likeService.GetLikeVideoCount(id)         // 点赞
	totalFavorited, _ := likeService.GetUserTotalIsLikedCount(id) // 被点赞
	fmtUser.TotalFavorited = totalFavorited
	fmtUser.FavoriteCount = favoriteCount

	return fmtUser, nil
}

// GenerateToken 根据username生成一个token
func GenerateToken(username string) string {
	u := UserService.GetUserByUsername(new(UserServiceImpl), username)
	token := NewToken(u)
	println(token)
	return token
}

// NewToken 根据用户信息创建token
func NewToken(user models.User) string {
	expiresTime := time.Now().Unix() + 60*60*24
	id64 := user.Id
	claims := jwt.StandardClaims{
		Audience:  user.Username,
		ExpiresAt: expiresTime,
		Id:        strconv.FormatInt(id64, 10),
		IssuedAt:  time.Now().Unix(),
		Issuer:    "tiktok",
		NotBefore: time.Now().Unix(),
		Subject:   "token",
	}
	var jwtSecret = []byte(config.Secret)
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	if token, err := tokenClaims.SignedString(jwtSecret); err == nil {
		token = "Bearer " + token
		println("generate token success!\n")
		return token
	} else {
		println("generate token fail\n")
		return "fail"
	}
}

// EnCoder 密码加密
func EnCoder(password string) string {
	h := hmac.New(sha256.New, []byte(password))
	sha := hex.EncodeToString(h.Sum(nil))
	return sha
}
