package service

import (
	"log"
	"strconv"
	"sync"
	"tiktokbackend/middleware/redis"
	"tiktokbackend/models"
	"time"
)

// FollowServiceImp 该结构体继承FollowService接口。
type FollowServiceImp struct {
}

// IsFollowing 当前用户是否关注了目标用户
func (fsi *FollowServiceImp) IsFollowing(userId int64, curId int64) (bool, error) {
	// 先查Redis里面是否有此关系。
	isFlowing, err := redis.RdbFollowings.SIsMember(redis.Ctx, strconv.Itoa(int(curId)), userId).Result()
	if isFlowing { // 如果当前用户关注了目标用户
		// 重新设置过期时间。
		redis.RdbFollowings.Expire(redis.Ctx, strconv.Itoa(int(curId)), time.Duration(2592000)*time.Second)
		return true, nil
	}
	// 如果redis内没查到，再查一次mysql
	relation, err := models.FindRelation(userId, curId)
	if relation == nil { // mysql数据库内也未查到关注信息，则直接返会错误
		return false, err
	} else { // mysql中查到了关注信息，将其注入Redis中。
		go func(userId int, curId int) {
			// 第一次存入时，给该key添加一个-1为key，防止脏数据的写入。当然set可以去重，直接加，便于CPU。
			redis.RdbFollowings.SAdd(redis.Ctx, strconv.Itoa(int(curId)), -1)
			// 将查询到的关注关系注入Redis.
			redis.RdbFollowings.SAdd(redis.Ctx, strconv.Itoa(int(curId)), userId)
			// 更新过期时间。
			redis.RdbFollowings.Expire(redis.Ctx, strconv.Itoa(int(curId)), time.Duration(2592000)*time.Second)
		}(int(userId), int(curId))

		return true, nil
	}

}

// GetTotalFollowersCnt 给定用户id，查询其粉丝数量。
func (fsi *FollowServiceImp) GetTotalFollowersCnt(userId int64) int64 {
	// 查Redis中是否已经存在。
	cnt, err := redis.RdbFollowers.SCard(redis.Ctx, strconv.Itoa(int(userId))).Result()
	if cnt > 0 { // 如果该用户被人关注过
		// 更新过期时间。
		redis.RdbFollowers.Expire(redis.Ctx, strconv.Itoa(int(userId)), time.Duration(2592000)*time.Second)
		// 减去默认值后返回
		return cnt - 1
	}
	// 如果redis中没查到，则查询mysql
	followersId, err := models.GetFollowersId(userId)
	if err != nil {
		log.Println("查询粉丝关注表时发生了错误：", err)
	}
	if followersId == nil {
		log.Println("mysql中follow表内也未查到粉丝信息")
		return 0
	}

	// 如果在mysql表中查询到了粉丝信息，则更新redis表
	go func(userId int, followersId []int64) {
		// 先加入默认值
		redis.RdbFollowers.SAdd(redis.Ctx, strconv.Itoa(userId), -1)
		// 随后逐个将粉丝id加入到redis数据库中
		for id := range followersId {
			redis.RdbFollowers.SAdd(redis.Ctx, strconv.Itoa(userId), id)
		}
		// 更新followers的过期时间。
		redis.RdbFollowers.Expire(redis.Ctx, strconv.Itoa(userId), time.Duration(2592000)*time.Second)
	}(int(userId), followersId)

	return int64(len(followersId))
}

// GetTotalFollowingsCnt 给定当前用户id，查询该用户的关注总数量。
func (fsi *FollowServiceImp) GetTotalFollowingsCnt(userId int64) int64 {
	// 查看Redis中是否有关注的key。
	cnt, err := redis.RdbFollowings.SCard(redis.Ctx, strconv.Itoa(int(userId))).Result()
	if cnt > 0 { // RdbFollowings内存在userId的key
		// 更新过期时间。
		redis.RdbFollowings.Expire(redis.Ctx, strconv.Itoa(int(userId)), time.Duration(2592000)*time.Second)
		return cnt - 1
	}
	// 查询mysql
	followingsId, err := models.GetFollowingsId(userId)

	if err != nil {
		log.Println("查询mysql表时发生错误: ", err)
	}
	if followingsId == nil {
		log.Println("查询到的关注总量为0")
		return 0
	}

	// mysql表中查到了关注信息，更新Redis中的followings
	go func(userId int, followingsId []int64) {
		// 先加上默认值
		redis.RdbFollowings.SAdd(redis.Ctx, strconv.Itoa(userId), -1)
		for id := range followingsId {
			redis.RdbFollowings.SAdd(redis.Ctx, strconv.Itoa(userId), id)
		}
		// 更新following的过期时间
		redis.RdbFollowings.Expire(redis.Ctx, strconv.Itoa(userId), time.Duration(2592000)*time.Second)
	}(int(userId), followingsId)

	return int64(len(followingsId))
}

// AddFollowRelation 当前用户关注目标用户
func (fsi *FollowServiceImp) AddFollowRelation(userId int64, curId int64) (bool, error) {
	userIdStr := strconv.Itoa(int(userId))
	curIdStr := strconv.Itoa(int(curId))

	cnt, _ := redis.RdbFollowers.SCard(redis.Ctx, userIdStr).Result()
	if cnt > 0 { // 如果RdbFollowers存在userIdStr的key，则直接增加value值curId
		// 目标用户增加一名粉丝(当前用户)
		redis.RdbFollowers.SAdd(redis.Ctx, userIdStr, curId)
		// 更新过期时间
		redis.RdbFollowers.Expire(redis.Ctx, userIdStr, time.Duration(2592000)*time.Second)
	}

	cnt1, _ := redis.RdbFollowings.SCard(redis.Ctx, curIdStr).Result()
	if cnt1 > 0 {
		redis.RdbFollowings.SAdd(redis.Ctx, curIdStr, userId) //当前用户增加一位关注者(目标用户)
		// 更新过期时间
		redis.RdbFollowings.Expire(redis.Ctx, curIdStr, time.Duration(2592000)*time.Second)
	}

	return true, nil

}

// DeleteFollowRelation 当前用户取消关注目标用户
func (fsi *FollowServiceImp) DeleteFollowRelation(userId int64, curId int64) (bool, error) {
	userIdStr := strconv.Itoa(int(userId))
	curIdStr := strconv.Itoa(int(curId))

	cnt, _ := redis.RdbFollowers.SCard(redis.Ctx, userIdStr).Result()
	if cnt > 0 { // 如果RdbFollowers存在userIdStr的key，则直接减去value值curId
		// 目标用户失去一名粉丝(当前用户)
		redis.RdbFollowers.SRem(redis.Ctx, userIdStr, curId)
		// 更新过期时间
		redis.RdbFollowers.Expire(redis.Ctx, userIdStr, time.Duration(2592000)*time.Second)
	}

	cnt1, _ := redis.RdbFollowings.SCard(redis.Ctx, curIdStr).Result()
	if cnt1 > 0 { // 如果RdbFolloings存在curIdStr的key，则直接减去value值userId
		// 当前用户失去一位关注者(目标用户)
		redis.RdbFollowings.SRem(redis.Ctx, curIdStr, userId)
		// 更新过期时间
		redis.RdbFollowings.Expire(redis.Ctx, curIdStr, time.Duration(2592000)*time.Second)
	}

	return true, nil
}

// GetFollowingsList 查询目标用户的关注列表
func (fsi *FollowServiceImp) GetFollowingsList(userId int64) ([]FmtUser, error) {
	// 先查询redis
	strUserId := strconv.FormatInt(userId, 10)
	var followingsIdList []int64

	// redis的RdbFollowings库中是否存在key值strUserId
	n, err := redis.RdbFollowings.Exists(redis.Ctx, strUserId).Result()
	if err != nil {
		log.Println("方法getFollowingsList中redis查询RdbFollowings库的key值失败")
		return nil, err
	}
	if n > 0 { // 如果存在
		// 获取该key所对应的所有value, 并拼装成followingsIdList
		strFollowingsIdList, err1 := redis.RdbFollowings.SMembers(redis.Ctx, strUserId).Result()
		for _, strFollowingId := range strFollowingsIdList {
			followingId, _ := strconv.ParseInt(strFollowingId, 10, 64)
			followingsIdList = append(followingsIdList, followingId)
		}
		if err1 != nil {
			log.Println("方法getFollowingsList中redis查询RdbFollowings库的value值失败")
		}
	} else {
		// key不存在
		// 为key值设置默认值
		if _, err := redis.RdbFollowings.SAdd(redis.Ctx, strUserId, -1).Result(); err != nil {
			log.Println("方法getFollowingsList中redis为RdbFollowings库的key值设置默认值失败")
			redis.RdbFollowings.Del(redis.Ctx, strUserId)
			return nil, err
		}
		// 为key值设置有效期
		if _, err := redis.RdbFollowings.Expire(redis.Ctx, strUserId, time.Duration(2592000)*time.Second).Result(); err != nil {
			log.Println("方法getFollowingsList中redis为RdbFollowings库的key值设置有效期失败")
			redis.RdbFollowings.Del(redis.Ctx, strUserId)
			return nil, err
		}
		// 查询mysql
		followingsIdList, err := models.GetFollowingsId(userId)
		if err != nil {
			log.Println("方法getFollowingsList中mysql查询失败")
			return nil, err
		}
		// 维护redis
		for _, followingId := range followingsIdList {
			if _, err := redis.RdbFollowings.SAdd(redis.Ctx, strUserId, followingId).Result(); err != nil {
				log.Println("方法getFollowingsList中redis维护失败")
				return nil, err
			}
		}

	}

	var userService UserServiceImpl
	followingsList := new([]FmtUser)

	// 根据每个id来查询用户信息。
	idListLen := len(followingsIdList) - 1
	if idListLen == 0 {
		return *followingsList, nil
	}

	// 创建协程组
	var wg sync.WaitGroup
	wg.Add(idListLen)
	for i := 0; i < idListLen; i++ {
		if followingsIdList[i] == -1 {
			continue
		}
		go func(pos int) {
			defer wg.Done()
			fmtUser, _ := userService.GetFmtUserByIdWithCurId(userId, followingsIdList[pos])
			log.Println("方法getFollowingsList输出fmtUser: ", fmtUser)
			*followingsList = append(*followingsList, fmtUser)
		}(i)
	}
	wg.Wait()
	// 返回关注对象列表。
	return *followingsList, nil
}

// GetFollowersList 查询用户的粉丝列表
func (fsi *FollowServiceImp) GetFollowersList(userId int64) ([]FmtUser, error) {
	// 先查询redis
	strUserId := strconv.FormatInt(userId, 10)

	var followersIdList []int64

	n, err := redis.RdbFollowers.Exists(redis.Ctx, strUserId).Result()
	if n > 0 {
		if err != nil {
			log.Println("方法getFollowingsList中redis查询RdbFollowings库的key值失败")
			return nil, err
		}
		strFollowersIdList, err1 := redis.RdbFollowers.SMembers(redis.Ctx, strUserId).Result()
		for _, strFollowingId := range strFollowersIdList {
			followerId, _ := strconv.ParseInt(strFollowingId, 10, 64)
			followersIdList = append(followersIdList, followerId)
		}
		if err1 != nil {
			log.Println("方法getFollowingsList中redis查询RdbFollowings库的value值失败")
		}
	} else {
		// key不存在
		// 为key值设置默认值
		if _, err := redis.RdbFollowers.SAdd(redis.Ctx, strUserId, -1).Result(); err != nil {
			log.Println("方法getFollowingsList中redis为RdbFollowings库的key值设置默认值失败")
			redis.RdbFollowings.Del(redis.Ctx, strUserId)
			return nil, err
		}
		// 为key值设置有效期
		if _, err := redis.RdbFollowers.Expire(redis.Ctx, strUserId, time.Duration(2592000)*time.Second).Result(); err != nil {
			log.Println("方法getFollowingsList中redis为RdbFollowings库的key值设置有效期失败")
			redis.RdbFollowings.Del(redis.Ctx, strUserId)
			return nil, err
		}
		// 查询mysql
		followersIdList, err := models.GetFollowingsId(userId)
		if err != nil {
			log.Println("方法getFollowingsList中mysql查询失败")
			return nil, err
		}
		// 维护redis
		for _, followerId := range followersIdList {
			if _, err := redis.RdbFollowers.SAdd(redis.Ctx, strUserId, followerId).Result(); err != nil {
				log.Println("方法getFollowingsList中redis维护失败")
				return nil, err
			}
		}

	}

	var userService UserServiceImpl
	followersList := new([]FmtUser)

	// 根据每个id来查询用户信息。
	idListLen := len(followersIdList) - 1
	if idListLen == 0 {
		return *followersList, nil
	}

	// 创建协程组
	var wg sync.WaitGroup
	wg.Add(idListLen)
	for i := 0; i < idListLen; i++ {
		if followersIdList[i] == -1 {
			continue
		}
		go func(pos int) {
			defer wg.Done()
			fmtUser, _ := userService.GetFmtUserByIdWithCurId(userId, followersIdList[pos])
			log.Println("方法getFollowingsList输出fmtUser: ", fmtUser)
			*followersList = append(*followersList, fmtUser)
		}(i)
	}
	wg.Wait()
	// 返回粉丝对象列表。
	return *followersList, nil
}
