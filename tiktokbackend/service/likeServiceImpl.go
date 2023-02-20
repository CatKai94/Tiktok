package service

import (
	"errors"
	"log"
	"strconv"
	"sync"
	"tiktokbackend/middleware/redis"
	"tiktokbackend/models"
	"time"
)

type LikeServiceImpl struct {
}

//redis有两个库
//第一个库为RdbUserToVideo,用key为userId,value为存储该userId的用户点赞过的视频的videoId的set
//第二个库为RdbVideoToUser，用key为videoId,value为存储所有点赞过该视频的用户的userId的set

// JudgeUserIsLike 根据userId,videoId,查询当前用户对该视频的点赞状态
func (likeService *LikeServiceImpl) JudgeUserIsLike(userId int64, videoId int64) (bool, error) {
	strUserId := strconv.FormatInt(userId, 10)
	strVideoId := strconv.FormatInt(videoId, 10)

	//查询RdbUserToVideo库,key为strUserId的set中是否存在
	count, err := redis.RdbUserToVideo.Exists(redis.Ctx, strUserId).Result()
	if count > 0 { //RdbUserToVideo库, 存在key为strUserId的set
		if err != nil {
			log.Println("方法JudgeUserIsLike RdbUserToVideo query key失败:", err)
			return false, err
		}
		//在key为strUserId的Set中用SIsMember判断是否有value:videoId存在
		result, err1 := redis.RdbUserToVideo.SIsMember(redis.Ctx, strUserId, videoId).Result()
		if err1 != nil {
			log.Println("方法JudgeUserIsLike RdbUserToVideo query value失败:", err1)
			return false, err1
		}

		return result, nil
	} else { // RdbUserToVideo不存在key为userId，则再查询RdbVideoToUser
		count, err := redis.RdbVideoToUser.Exists(redis.Ctx, strVideoId).Result()
		if count > 0 {
			//err不为空，查询失败
			if err != nil {
				log.Println("方法JudgeUserIsLike RdbVideoToUser query key失败:")
				return false, err
			}
			//在key为strVideo的set中用SIsMember判断是否有value: userId存在
			result, err1 := redis.RdbVideoToUser.SIsMember(redis.Ctx, strVideoId, userId).Result()
			if err1 != nil {
				log.Println("方法JudgeUserIsLike RdbVideoToUser query value失败:", err1)
				return false, err1
			}
			log.Println("方法JudgeUserIsLike RdbVideoToUser query value成功")
			return result, nil
		} else { //两个redis库中都没有对应的key
			//给优先查询的RdbUserToVideo库中加入key:setUserId,value:-1
			_, err := redis.RdbUserToVideo.SAdd(redis.Ctx, strUserId, -1).Result()
			if err != nil {
				log.Println("方法JudgeUserIsLike 添加-1默认值失败")
				redis.RdbUserToVideo.Del(redis.Ctx, strUserId)
				return false, err
			}
			//给默认值设置有效期
			_, err1 := redis.RdbUserToVideo.Expire(redis.Ctx, strUserId, time.Duration(2592000)*time.Second).Result()
			if err1 != nil {
				log.Println("方法JudgeUserIsLike 设置有效期失败")
				redis.RdbUserToVideo.Del(redis.Ctx, strUserId)
				return false, nil
			}
			//查询mysql 根据userId查询likes表，返回所有该用户点赞过的视频并加入redis中RdbUserToVideo key: strUserId中
			videoIdList, err2 := models.GetVideoIdList(userId)
			if err2 != nil {
				log.Println(err2.Error())
				return false, err2
			}
			//加入redis
			for _, videoIdResult := range videoIdList {
				redis.RdbUserToVideo.SAdd(redis.Ctx, strUserId, videoIdResult)
			}
			result, err3 := redis.RdbUserToVideo.SIsMember(redis.Ctx, strUserId, videoId).Result()
			if err3 != nil {
				log.Println("方法JudgeUserIsLike在维护redis后 查询RdbUserToVideo value失败")
				return false, err3
			}
			log.Println("方法JudgeUserIsLike在维护redis后 查询RdbUserToVideo value成功")
			return result, err3
		}
	}
}

// GetLikeCount 根据videoId获取当前视频总共的喜爱量
func (likeService *LikeServiceImpl) GetLikeCount(videoId int64) (int64, error) {
	strVideoId := strconv.FormatInt(videoId, 10)
	n, err := redis.RdbVideoToUser.Exists(redis.Ctx, strVideoId).Result()
	if n > 0 {
		if err != nil {
			log.Println("方法GetLikeCount Redis查询key失败")
			return 0, err
		}
		result, err1 := redis.RdbVideoToUser.SCard(redis.Ctx, strVideoId).Result()
		if err1 != nil {
			log.Println("方法GetLikeCount Redis查询value失败")
			return 0, err1
		}
		log.Println("方法GetLikeCount Redis查询value成功")
		return result - 1, err1
	} else {
		//RdbVideoToUser中没有对应的key,需要查询mysql
		//加入defaultRedisValue
		_, err := redis.RdbVideoToUser.SAdd(redis.Ctx, strVideoId, -1).Result()
		if err != nil {
			log.Println("方法GetLikeCount redis 添加默认值失败")
			redis.RdbVideoToUser.Del(redis.Ctx, strVideoId)
			return 0, err
		}
		//设置有效期
		_, err = redis.RdbVideoToUser.Expire(redis.Ctx, strVideoId, time.Duration(2592000)*time.Second).Result()
		if err != nil {
			log.Println("方法GetLikeCount redis 设置有效期失败")
			redis.RdbVideoToUser.Del(redis.Ctx, strVideoId)
			return 0, err
		}
		//查询mysql
		userIdList, err1 := models.GetUserIdList(videoId)
		if err1 != nil {
			log.Println(err1.Error())
			return 0, err1
		}
		for _, userIdResult := range userIdList {
			redis.RdbVideoToUser.SAdd(redis.Ctx, strVideoId, userIdResult)
		}
		result, err2 := redis.RdbVideoToUser.SCard(redis.Ctx, strVideoId).Result()
		if err2 != nil {
			log.Println("方法GetLikeCount 在维护redis后查询redis失败")
			return 0, err2
		}
		log.Println("方法GetLikeCount 在维护redis后查询成功")
		return result - 1, nil
	}
}

// GetUserTotalIsLikedCount 根据当前用户userId,查询当前用户总共被点赞个数
func (likeService *LikeServiceImpl) GetUserTotalIsLikedCount(userId int64) (int64, error) {
	//调用model.video中的方法
	videoIdList, err := models.GetVideoIdsByAuthorId(userId)
	if err != nil {
		log.Println("方法GetUserTotalIsLikedCount发生了错误： ", err.Error())
		return 0, err
	}
	var result int64
	videoLikeCountList := new([]int64)
	i := len(videoIdList)
	var wg sync.WaitGroup
	wg.Add(i)
	for j := 0; j < i; j++ {
		go func(pos int) {
			defer wg.Done()
			//调用GetLikeCount:根据videoId,获取点赞数
			count, err := likeService.GetLikeCount(videoIdList[pos])
			if err != nil {
				//如果有错误，输出错误信息，并不加入该视频点赞数
				log.Printf(err.Error())
				return
			}
			*videoLikeCountList = append(*videoLikeCountList, count)
		}(j)
	}
	wg.Wait()
	for _, count := range *videoLikeCountList {
		result += count
	}
	return result, nil
}

// GetLikeVideoCount 根据userId获得该用户点赞过的视频数量
func (likeService *LikeServiceImpl) GetLikeVideoCount(userId int64) (int64, error) {
	strUserId := strconv.FormatInt(userId, 10)
	n, err := redis.RdbUserToVideo.Exists(redis.Ctx, strUserId).Result()
	if n > 0 { // RdbUserToVideo
		if err != nil {
			log.Println("方法GetLikeVideoCount redis查询key失败")
			return 0, err
		}
		result, err1 := redis.RdbUserToVideo.SCard(redis.Ctx, strUserId).Result()
		if err1 != nil {
			log.Println("方法GetLikeVideoCount redis查询value失败")
			return 0, err1
		}
		log.Println("方法GetLikeVideoCount redis查询key成功")
		// 减去默认值-1
		return result - 1, nil
	} else {
		//RdbUserToVideo中没有对应的Key
		if _, err := redis.RdbUserToVideo.SAdd(redis.Ctx, strUserId, -1).Result(); err != nil {
			log.Printf("方法GetLikeVideoCount 添加默认值失败")
			redis.RdbUserToVideo.Del(redis.Ctx, strUserId)
			return 0, err
		}
		//给键值设置有效期
		_, err := redis.RdbUserToVideo.Expire(redis.Ctx, strUserId,
			time.Duration(259200)*time.Second).Result()
		if err != nil {
			log.Printf("方法GetLikeVideoCount 设置有效期失败")
			redis.RdbUserToVideo.Del(redis.Ctx, strUserId)
			return 0, err
		}
		videoIdList, err1 := models.GetVideoIdList(userId)
		if err1 != nil {
			log.Println(err1.Error())
			return 0, err1
		}
		//维护redis
		for _, videoIdResult := range videoIdList {
			redis.RdbUserToVideo.SAdd(redis.Ctx, strUserId, videoIdResult)
		}
		//再查询redis
		result, err2 := redis.RdbUserToVideo.SCard(redis.Ctx, strUserId).Result()
		if err2 != nil {
			log.Println("方法GetLikeVideoCount 在维护redis后查询value失败")
			return 0, err2
		}
		log.Println("方法GetLikeVideoCount 在维护redis后查询value成功")
		return result - 1, nil
	}
}

func (likeService *LikeServiceImpl) LikeAction(userId int64, videoId int64, actionType int32) error {
	strUserId := strconv.FormatInt(userId, 10)
	strVideoId := strconv.FormatInt(videoId, 10)

	//点赞 actionType == 1
	if actionType == 1 {
		//查询RbdUserToVideo
		n, err := redis.RdbUserToVideo.Exists(redis.Ctx, strUserId).Result()
		if n > 0 {
			if err != nil {
				log.Println("方法LikeAction RdbUserToVideo redis查询key失败")
				return err
			}
			//key: strUserId存在,则添加value:videoId
			_, err1 := redis.RdbUserToVideo.SAdd(redis.Ctx, strUserId, videoId).Result()
			if err1 != nil {
				log.Println("方法LikeAction RdbUserToVideo redis添加value失败")
				return err1
			} else {
				//先在数据库查找有没有该数据
				likeInfo, err2 := models.QueryLikeInfo(userId, videoId)
				if err2 != nil {
					log.Println(err2.Error())
					return err2
				}
				//没有该条数据就添加
				if likeInfo == (models.Like{}) {
					if err3 := models.AddLike(userId, videoId); err3 != nil {
						log.Println(err3.Error())
						return err3
					}
				} else {
					//有该数据，则更新
					if err3 := models.UpdateLikeAction(userId, videoId, actionType); err3 != nil {
						log.Println(err3.Error())
						return err3
					}
				}
			}

		} else {
			//若key:strUserId不存在
			if _, err := redis.RdbUserToVideo.SAdd(redis.Ctx, strUserId, -1).Result(); err != nil {
				log.Println("方法LikeAction RdbUserToVideo 添加默认值失败")
				redis.RdbUserToVideo.Del(redis.Ctx, strUserId)
				return err
			}
			//给键值设置有效期
			_, err := redis.RdbUserToVideo.Expire(redis.Ctx, strUserId,
				time.Duration(259200)*time.Second).Result()
			if err != nil {
				log.Println("方法LikeAction RdbUserToVideo 设置有效期失败")
				redis.RdbUserToVideo.Del(redis.Ctx, strUserId)
				return err
			}
			//查询mysql
			videoIdList, err1 := models.GetVideoIdList(userId)
			if err1 != nil {
				log.Println(err1.Error())
				return err1
			}
			//遍历videoIdList,添加进key的集合中，若失败，删除key，并返回错误信息
			for _, likeVideoId := range videoIdList {
				if _, err1 := redis.RdbUserToVideo.SAdd(redis.Ctx, strUserId, likeVideoId).Result(); err1 != nil {
					log.Println("方法LikeAction RdbUserToVideo 维护redis失败")
					redis.RdbUserToVideo.Del(redis.Ctx, strUserId)
					return err1
				}
			}
			if _, err2 := redis.RdbUserToVideo.SAdd(redis.Ctx, strUserId, videoId).Result(); err2 != nil {
				log.Println("方法LikeAction RdbUserToVideo add value失败")
				return err2
			} else {
				//先在数据库查找有没有该数据
				likeInfo, err2 := models.QueryLikeInfo(userId, videoId)
				if err2 != nil {
					log.Println(err2.Error())
					return err2
				}
				//没有该条数据就添加
				if likeInfo == (models.Like{}) {
					if err3 := models.AddLike(userId, videoId); err3 != nil {
						log.Println(err3.Error())
						return err3
					}
				} else {
					//有该数据，则更新
					if err3 := models.UpdateLikeAction(userId, videoId, actionType); err3 != nil {
						log.Println(err3.Error())
						return err3
					}
				}
			}
		}
		//查询RdbVideoToUser
		n, err = redis.RdbVideoToUser.Exists(redis.Ctx, strVideoId).Result()
		if n > 0 {
			if err != nil {
				log.Println("方法LikeAction RdbVideoToUser 查询key失败")
				return err
			}
			//key:strVideoId存在,添加value:userId
			_, err1 := redis.RdbVideoToUser.SAdd(redis.Ctx, strVideoId, userId).Result()
			if err1 != nil {
				log.Println("方法LikeAction RdbVideoToUser 添加value失败")
				return err1
			}
		} else {
			//若Key:strVideoId不存在
			if _, err := redis.RdbVideoToUser.SAdd(redis.Ctx, strVideoId, -1).Result(); err != nil {
				log.Println("方法LikeAction RdbVideoToUser 添加默认值失败")
				redis.RdbVideoToUser.Del(redis.Ctx, strVideoId)
				return err
			}
			//给键值设置有效期，类似于gc机制
			_, err := redis.RdbVideoToUser.Expire(redis.Ctx, strVideoId,
				time.Duration(259200)*time.Second).Result()
			if err != nil {
				log.Println("方法LikeAction RdbVideoToUser设置有效期失败")
				redis.RdbVideoToUser.Del(redis.Ctx, strVideoId)
				return err
			}
			//查询mysql
			userIdList, err1 := models.GetUserIdList(videoId)
			//如果有问题，说明查询失败，返回错误信息："get likeUserIdList failed"
			if err1 != nil {
				return err1
			}
			//遍历userIdList,添加进key的集合中，若失败，删除key，并返回错误信息，这么做的原因是防止脏读，
			//保证redis与mysql数据一致性
			for _, likeUserId := range userIdList {
				if _, err1 := redis.RdbVideoToUser.SAdd(redis.Ctx, strVideoId, likeUserId).Result(); err1 != nil {
					log.Println("方法LikeAction RdbVideoToUser 维护redis失败")
					redis.RdbVideoToUser.Del(redis.Ctx, strVideoId)
					return err1
				}
			}
			//这样操作理由同上
			if _, err2 := redis.RdbVideoToUser.SAdd(redis.Ctx, strVideoId, userId).Result(); err2 != nil {
				log.Println("方法LikeAction RdbVideoToUser 在维护redis后添加value失败", err2)
				return err2
			}
		}
	} else {
		//actionType == 2,进行取消赞操作
		//查询RbdUserToVideo
		n, err := redis.RdbUserToVideo.Exists(redis.Ctx, strUserId).Result()
		if n > 0 {
			if err != nil {
				log.Println("方法LikeAction RdbUserToVideo 查询key失败")
				return err
			}
			//在redis中删除该点赞记录
			if _, err1 := redis.RdbUserToVideo.SRem(redis.Ctx, strUserId, videoId).Result(); err1 != nil {
				log.Println("方法LikeAction RdbUserToVideo 删除value失败")
				return err1
			} else {
				//在redis中删除成功，进而对mysql操作
				//取消赞必定在点赞状态下取消,直接更新即可
				if err1 := models.UpdateLikeAction(userId, videoId, actionType); err1 != nil {
					log.Println(err1.Error())
					return err1
				}
			}
		} else {
			//Rdb UserToVideo 中不存在key:strUserId
			if _, err := redis.RdbUserToVideo.SAdd(redis.Ctx, strUserId, -1).Result(); err != nil {
				log.Println("方法LikeAction RdbUserToVideo 添加默认值失败")
				redis.RdbUserToVideo.Del(redis.Ctx, strUserId)
				return err
			}
			//给键值设置有效期，类似于gc机制
			_, err := redis.RdbUserToVideo.Expire(redis.Ctx, strUserId,
				time.Duration(259200)*time.Second).Result()
			if err != nil {
				log.Println("方法LikeAction RdbUserToVideo 设置有效期失败")
				redis.RdbUserToVideo.Del(redis.Ctx, strUserId)
				return err
			}
			//查询mysql
			videoIdList, err1 := models.GetVideoIdList(userId)
			//如果有问题，说明查询失败，返回错误信息："get likeVideoIdList failed"
			if err1 != nil {
				return err1
			}
			//遍历videoIdList,添加进key的集合中，若失败，删除key，并返回错误信息，这么做的原因是防止脏读，
			//保证redis与mysql 数据原子性
			for _, likeVideoId := range videoIdList {
				if _, err1 := redis.RdbUserToVideo.SAdd(redis.Ctx, strUserId, likeVideoId).Result(); err1 != nil {
					log.Println("方法LikeAction RdbUserToVideo 维护redis失败")
					redis.RdbUserToVideo.Del(redis.Ctx, strUserId)
					return err1
				}
			}
			//这样操作理由同上
			if _, err2 := redis.RdbUserToVideo.SRem(redis.Ctx, strUserId, videoId).Result(); err2 != nil {
				log.Println("方法LikeAction RdbUserToVideo 删除value失败")
				return err2
			} else {
				//在redis中删除成功，进而对mysql操作
				//取消赞必定在点赞状态下取消,直接更新即可
				if err1 := models.UpdateLikeAction(userId, videoId, actionType); err1 != nil {
					log.Println(err1.Error())
					return err1
				}
			}
		}
		//查询Rdb RdbVideoToUserId 是否有key:strVideoId
		if n, err := redis.RdbUserToVideo.Exists(redis.Ctx, strVideoId).Result(); n > 0 {
			//如果有问题，说明查询redis失败,返回错误信息
			if err != nil {
				log.Println("方法LikeAction RdbUserToVideo 查询key失败")
				return err
			} //如果加载过此信息key:strVideoId，则删除value:userId
			//如果redis LikeVideoId 删除失败，返回错误信息
			if _, err1 := redis.RdbUserToVideo.SRem(redis.Ctx, strVideoId, userId).Result(); err1 != nil {
				log.Println("方法LikeAction RdbUserToVideo 删除value失败")
				return err1
			}
		} else {
			//不存在key
			if _, err := redis.RdbUserToVideo.SAdd(redis.Ctx, strVideoId, -1).Result(); err != nil {
				log.Println("方法LikeAction RdbUserToVideo 添加默认值失败")
				redis.RdbUserToVideo.Del(redis.Ctx, strVideoId)
				return err
			}
			_, err = redis.RdbUserToVideo.Expire(redis.Ctx, strVideoId,
				time.Duration(259200)*time.Second).Result()
			if err != nil {
				log.Println("方法LikeAction RdbUserToVideo 设置有效期失败")
				redis.RdbUserToVideo.Del(redis.Ctx, strVideoId)
				return err
			}
			//查询mysql
			userIdList, err1 := models.GetUserIdList(videoId)
			//如果有问题，说明查询失败，返回错误信息："get likeUserIdList failed"
			if err1 != nil {
				redis.RdbUserToVideo.Del(redis.Ctx, strVideoId)
				return err1
			}
			//遍历userIdList,添加进key的集合中，若失败，删除key，并返回错误信息，这么做的原因是防止脏读，
			//保证redis与mysql数据一致性
			for _, likeUserId := range userIdList {
				if _, err1 := redis.RdbUserToVideo.SAdd(redis.Ctx, strVideoId, likeUserId).Result(); err1 != nil {
					log.Println("方法LikeAction RdbUserToVideo 添加默认值失败")
					redis.RdbUserToVideo.Del(redis.Ctx, strVideoId)
					return err1
				}
			}
			//这样操作理由同上
			if _, err2 := redis.RdbUserToVideo.SRem(redis.Ctx, strVideoId, userId).Result(); err2 != nil {
				log.Printf("方法:FavouriteAction RedisLikeVideoId del value失败：%v", err2)
				return err2
			}
		}
	}
	return nil
}

func (likeService *LikeServiceImpl) GetLikeVideoList(userId int64, curId int64) ([]FmtVideo, error) {
	strUserId := strconv.FormatInt(userId, 10)
	// var videoIdList []int64
	// var ss = make([]int, n)
	var videoIdList = make([]int64, 0)
	n, err := redis.RdbUserToVideo.Exists(redis.Ctx, strUserId).Result()
	if n > 0 {
		log.Println("到这里了111")
		if err != nil {
			log.Println("方法GetLikeVideoList redis查询key失败")
			return nil, err
		}
		strVideoIdList, err1 := redis.RdbUserToVideo.SMembers(redis.Ctx, strUserId).Result()
		for _, strVideoId := range strVideoIdList {
			videoId, _ := strconv.ParseInt(strVideoId, 10, 64)
			videoIdList = append(videoIdList, videoId)
			log.Println("502行  videoIdList: ", videoIdList)
		}
		if err1 != nil {
			log.Println("方法GetLikeVideoList redis查询value失败")
		}
	} else {

		log.Println("到这里了222")
		//key不存在
		if _, err := redis.RdbUserToVideo.SAdd(redis.Ctx, strUserId, -1).Result(); err != nil {
			log.Println("方法GetLikeVideoList redis添加默认值失败")
			redis.RdbUserToVideo.Del(redis.Ctx, strUserId)
			return nil, err
		}
		_, err := redis.RdbUserToVideo.Expire(redis.Ctx, strUserId, time.Duration(259200)*time.Second).Result()
		if err != nil {
			log.Println("方法GetLikeVideoList redis设置有效期失败")
			redis.RdbUserToVideo.Del(redis.Ctx, strUserId)
			return nil, err
		}
		//查询mysql
		videoIdList, err1 := models.GetVideoIdList(userId)
		if err1 != nil {
			log.Println(err1.Error())
			return nil, err1
		}
		//维护redis
		for _, videIdResult := range videoIdList {
			if _, err2 := redis.RdbUserToVideo.SAdd(redis.Ctx, strUserId, videIdResult).Result(); err2 != nil {
				log.Println("方法GetLikeVideoList 维护redis失败")
				return nil, err2
			}
		}
	}

	likeVideoList := new([]FmtVideo)
	i := len(videoIdList) - 1
	log.Println("537行  len videoIdList", i+1)

	if i == 0 {
		return *likeVideoList, nil
	}

	var videoService VideoServiceImpl
	// 创建协程计数器
	var wg sync.WaitGroup
	wg.Add(i)
	for j := 0; j <= i; j++ {
		if videoIdList[j] == -1 {
			continue
		}
		log.Println("创建协程: ", j)
		go func(pos int) {
			defer wg.Done()
			fmtVideo, err := videoService.GetVideo(videoIdList[pos], curId)

			log.Println("video: ", fmtVideo)
			if err != nil {
				log.Println(errors.New("该喜欢的视频丢失"))
				return
			}
			*likeVideoList = append(*likeVideoList, fmtVideo)
		}(j)
	}
	wg.Wait()

	return *likeVideoList, nil
}
