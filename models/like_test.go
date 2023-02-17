package models

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddLike(t *testing.T) {
	err := AddLike(203, 100)
	fmt.Println(err)
}

func TestQueryLikeInfo(t *testing.T) {
	likeInfo, err := QueryLikeInfo(203, 200)
	fmt.Println(likeInfo)
	fmt.Println(err)
}

func TestGetVideoIdList(t *testing.T) {
	videoIdList, err := GetVideoIdList(203)
	fmt.Println(videoIdList)
	fmt.Println(err)
}

func TestGetUserIdList(t *testing.T) {

	userIdList, err := GetUserIdList(200)
	fmt.Println(userIdList)
	fmt.Println(err)
}

func TestUpdateLikeAction(t *testing.T) {
	err := UpdateLikeAction(203, 100, 1)
	fmt.Println(err)
	likeInfo, _ := QueryLikeInfo(203, 100)
	fmt.Println(likeInfo)
	assert.Equal(t, 1, int(likeInfo.Cancel))
}
