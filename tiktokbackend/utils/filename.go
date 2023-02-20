package utils

import (
	"strconv"
	"time"
)

func CreateFileName(userId int64) string {
	//将用户ID和时间戳拼接成独一无二的命名
	timeUnix := time.Now().Unix()
	strTimeUnix := strconv.FormatInt(timeUnix, 10)
	strId := strconv.FormatInt(userId, 10)
	fileName := strId + strTimeUnix

	return fileName
}
