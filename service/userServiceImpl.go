package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
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
		Id:            0,
		Name:          "",
		FollowCount:   0,
		FollowerCount: 0,
		IsFollow:      false,
		TotalFavorite: 0,
		FavoriteCount: 0,
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
	fmtUser := FmtUser{
		Id:            0,
		Name:          "",
		FollowCount:   0,
		FollowerCount: 0,
		IsFollow:      false,
		TotalFavorite: 0,
		FavoriteCount: 0,
	}

	//user, err := models.GetUserById(id)
	//if err != nil {
	//	log.Println("Err:", err.Error())
	//	log.Println("User Not Found")
	//	return fmtUser, err
	//}
	//
	//log.Println("Query User Success")
	//followCount, err := usi.GetFollowingCnt(id)
	//if err != nil {
	//	log.Println("Err:", err.Error())
	//}
	//followerCount, err := usi.GetFollowerCnt(id)
	//if err != nil {
	//	log.Println("Err:", err.Error())
	//}
	//isfollow, err := usi.IsFollowing(curId, id)
	//if err != nil {
	//	log.Println("Err:", err.Error())
	//}
	//u := GetLikeService() //解决循环依赖
	//totalFavorited, _ := u.TotalFavourite(id)
	//favoritedCount, _ := u.FavouriteVideoCount(id)
	//fmtUser = FmtUser{
	//	Id:            id,
	//	Name:          user.Username,
	//	FollowCount:   followCount,
	//	FollowerCount: followerCount,
	//	IsFollow:      isfollow,
	//	TotalFavorite: totalFavorited,
	//	FavoriteCount: favoritedCount,
	//}
	return fmtUser, nil
}

// GenerateToken 根据username生成一个token
func GenerateToken(username string) string {
	u := UserService.GetUserByUsername(new(UserServiceImpl), username)
	fmt.Printf("generatetoken: %v\n", u)
	token := NewToken(u)
	println(token)
	return token
}

// NewToken 根据信息创建token
func NewToken(u models.User) string {
	expiresTime := time.Now().Unix() + 60*60*24
	fmt.Printf("expiresTime: %v\n", expiresTime)
	id64 := u.Id
	fmt.Printf("id: %v\n", strconv.FormatInt(id64, 10))
	claims := jwt.StandardClaims{
		Audience:  u.Username,
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
	fmt.Println("Result: " + sha)
	return sha
}
