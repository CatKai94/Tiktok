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
	log.Println("当前用户: ", curId, "   关注了目标用户：", userId, "   ", isFlowing)
	if isFlowing { // 如果当前用户关注了目标用户
		// 重新设置过期时间。
		redis.RdbFollowings.Expire(redis.Ctx, strconv.Itoa(int(curId)), time.Duration(2592000)*time.Second)
		return true, nil
	}
	log.Println("redis内没查到当前用户关注过目标用户")
	// 如果redis内没查到，再查一次mysql
	relation, err := models.QueryFollowInfo(userId, curId)
	if relation == (models.Follow{}) { // mysql数据库内也未查到关注信息，则直接返会错误
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

	// 如果redis内没有查到，则查询mysql
	followingsId, err := models.GetFollowingsId(userId)
	// 如果mysql内查到了，则回写至redis

	if err != nil {
		log.Println("查询mysql表时发生错误: ", err)
	}
	if followingsId == nil {
		log.Println("查询到的关注总量为0")
		return 0
	}

	// mysql表中查到了关注信息，更新Redis中的followings
	go func(userId int, followingsId []int64) {
		// 先加上默认值, 默认值为自己
		redis.RdbFollowings.SAdd(redis.Ctx, strconv.Itoa(userId), userId)
		for id := range followingsId {
			redis.RdbFollowings.SAdd(redis.Ctx, strconv.Itoa(userId), id)
		}
		// 更新following的过期时间
		redis.RdbFollowings.Expire(redis.Ctx, strconv.Itoa(userId), time.Duration(2592000)*time.Second)
	}(int(userId), followingsId)

	return int64(len(followingsId))
}

// FollowAction 当前用户关注目标用户
func (fsi *FollowServiceImp) FollowAction(userId int64, curId int64) error {
	strUserId := strconv.Itoa(int(userId))
	strCurId := strconv.Itoa(int(curId))

	cnt, _ := redis.RdbFollowers.Exists(redis.Ctx, strUserId).Result()
	log.Println("followService 120行    目标用户：", userId, "  当前用户：", curId)
	if cnt > 0 { // 如果RdbFollowers存在userIdStr的key，则直接增加value值curId
		// 目标用户增加一名粉丝(当前用户)
		if _, err := redis.RdbFollowers.SAdd(redis.Ctx, strUserId, curId).Result(); err != nil {
			log.Println("方法FollowAction RdbFollowers 为strUserId增加curId失败")
			redis.RdbFollowers.Del(redis.Ctx, strUserId)
			return err
		}
		if _, err := redis.RdbFollowers.Expire(redis.Ctx, strUserId, time.Duration(2592000)*time.Second).Result(); err != nil {
			log.Println("方法FollowAction RdbFollowers 为键值strUserId更新过期时间失败")
			redis.RdbFollowers.Del(redis.Ctx, strUserId)
			return err
		}
		// 更新mysql数据库
		// 先在数据库内查查是否有该数据
		followInfo, _ := models.QueryEverFollowInfo(userId, curId)
		if followInfo == (models.Follow{}) { // 没有这条数据就添加一个
			log.Println("数据库内没查到相关信息！！！！")
			if err1 := models.InsertFollowRelation(userId, curId); err1 != nil {
				log.Println("方法FollowAction 向数据库中插入记录时发生错误: ", err1)
				return err1
			}
		} else { // 如果有这条数据，就更新一下          cancel 1为关注，2为取消关注
			if err1 := models.UpdateFollowRelation(userId, curId, 1); err1 != nil {
				return err1
			}
		}

	} else { // 如果RdbFollowers不存在userIdStr的key
		// 添加默认值， 默认值为自己
		if _, err := redis.RdbFollowers.SAdd(redis.Ctx, strUserId, userId).Result(); err != nil {
			log.Println("方法FollowAction RdbFollowers 添加默认值失败")
			redis.RdbFollowers.Del(redis.Ctx, strUserId)
			return err
		}
		// 给键值设置有效期
		if _, err := redis.RdbFollowers.Expire(redis.Ctx, strUserId, time.Duration(2592000)*time.Second).Result(); err != nil {
			log.Println("方法FollowAction RdbFollowers 为键值设置有效期失败")
			redis.RdbFollowers.Del(redis.Ctx, strUserId)
			return err
		}
		// 查询mysql
		followerIdList, err := models.GetFollowersId(userId)
		if err != nil {
			log.Println("方法FollowAction 查询数据库时发生错误：", err)
			return err
		}
		// 如果在mysql表中查到了关注的记录，则要维护一下redis数据库
		for _, followerId := range followerIdList {
			if _, err := redis.RdbFollowers.SAdd(redis.Ctx, strUserId, followerId).Result(); err != nil {
				log.Println("方法FollowAction 维护redis数据库时发生错误：", err)
				redis.RdbFollowers.Del(redis.Ctx, strUserId)
				return err
			}
		}

		// 开始插入最新的关注记录
		if _, err1 := redis.RdbFollowers.SAdd(redis.Ctx, strUserId, curId).Result(); err1 != nil {
			log.Println("方法FollowAction 向redis数据库中插入数据失败：", err)
			return err1
		} else {
			//成功在redis内记录好数据后，重新再在mysql表中记录一下
			// 先在数据库中查一下有没有该数据
			followInfo, err2 := models.QueryFollowInfo(userId, curId)
			if err2 != nil {
				log.Println("方法FollowAction 查询数据时发生错误: ", err1)
				return err2
			}
			// 没有这条数据就添加进去
			if followInfo == (models.Follow{}) {
				if err3 := models.InsertFollowRelation(userId, curId); err3 != nil {
					log.Println("方法FollowAction 插入数据时发生错误：", err3)
					return err3
				}
			} else { //如果有该数据，就更新一下
				if err3 := models.UpdateFollowRelation(userId, curId, 1); err3 != nil {
					log.Println("方法FollowAction 更新数据时发生错误：", err3)
					return err3
				}
			}
		}
	}

	// 更新关注者RdbFollowings数据库部分
	n, err := redis.RdbFollowings.Exists(redis.Ctx, strCurId).Result()
	if err != nil {
		log.Println("方法FollowAction 查询RdbFollowings时发生错误：", err)
		return err
	}
	if n > 0 { // 若key strCurId存在
		//当前用户增加一位关注者(目标用户)
		if _, err := redis.RdbFollowings.SAdd(redis.Ctx, strCurId, userId).Result(); err != nil {
			log.Println("方法FollowAction RdbFollowings 为key strCurId插入value失败：", err)
			return err
		}
	} else { // 若key strCurId不存在
		// 为key设置默认值， 默认值为自己
		if _, err := redis.RdbFollowings.SAdd(redis.Ctx, strCurId, userId).Result(); err != nil {
			log.Println("方法FollowAction RdbFollowings 为key strCurId设置默认值失败：", err)
			return err
		}
		// 为键值更新有效期
		if _, err := redis.RdbFollowings.Expire(redis.Ctx, strCurId, time.Duration(2592000)*time.Second).Result(); err != nil {
			log.Println("方法FollowAction RdbFollowings 为key strCurId更新有效期失败：", err)
			redis.RdbFollowings.Del(redis.Ctx, strCurId)
			return err
		}
		// 查询mysql
		followingsIdList, err := models.GetFollowingsId(curId)
		if err != nil {
			log.Println("方法FollowAction 查询mysql followings失败")
			return err
		}
		// 如果查到了记录，则维护一次redis
		for _, followingId := range followingsIdList {
			if _, err := redis.RdbFollowings.SAdd(redis.Ctx, strCurId, followingId).Result(); err != nil {
				log.Println("方法FollowAction RdbFollowings 维护redis失败：", err)
				redis.RdbFollowings.Del(redis.Ctx, strCurId)
				return err
			}
		}
		// 维护过后再次插入最新的记录
		if _, err := redis.RdbFollowings.SAdd(redis.Ctx, strCurId, userId).Result(); err != nil {
			log.Println("方法FollowAction RdbFollowings 插入key strCurId失败：", err)
			return err
		}
	}

	return nil
}

// UnFollowAction 当前用户取消关注目标用户
func (fsi *FollowServiceImp) UnFollowAction(userId int64, curId int64) error {
	strUserId := strconv.Itoa(int(userId))
	strCurId := strconv.Itoa(int(curId))

	// 查询RdbFollowers
	n, err := redis.RdbFollowers.Exists(redis.Ctx, strUserId).Result()
	if err != nil {
		log.Println("方法UnFollowAction 查询RdbFollower的key值userIdstr发生错误: ", err)
		return err
	}
	if n > 0 { // RdbFollowers中存在该key值
		// 在redis中删除该关注记录
		if _, err := redis.RdbFollowers.SRem(redis.Ctx, strUserId, curId).Result(); err != nil {
			log.Println("方法UnFollowAction 删除RdbFollower的key值失败: ", err)
			return err
		} else { // 在redis中删除成功后，对mysql表进行该操作
			if err := models.UpdateFollowRelation(userId, curId, 2); err != nil {
				log.Println("方法UnFollowAction 更改mysql的记录失败: ", err)
				return err
			}
		}
	} else { // RdbFollowers中不存在key值strUserId
		// 对RdbFollower进行系列操作
		// 为RdbFollowers设置默认值， 默认值为自己的id
		if _, err := redis.RdbFollowers.SAdd(redis.Ctx, strUserId, userId).Result(); err != nil {
			log.Println("方法UnFollowAction RdbFollowers设置默认键值失败：", err)
			return err
		}
		// 为键值更新有效时间
		if _, err := redis.RdbFollowers.Expire(redis.Ctx, strUserId, time.Duration(2592000)*time.Second).Result(); err != nil {
			log.Println("方法UnFollowAction RdbFollowers 为strUserId更新默认时间：", err)
			redis.RdbFollowers.Del(redis.Ctx, strUserId)
			return err
		}
		// 查询mysql
		followerIdList, err := models.GetFollowersId(userId)
		if err != nil {
			log.Println("方法UnFollowAction 查询mysql发生错误: ", err)
			return err
		}
		// 维护redis
		for _, followerId := range followerIdList {
			if _, err := redis.RdbFollowers.SAdd(redis.Ctx, strUserId, followerId).Result(); err != nil {
				log.Println("方法UnFollowAction 维护redis失败: ", err)
				return err
			}
		}
		// 维护结束后，开始取消关注操作
		if _, err := redis.RdbFollowers.SRem(redis.Ctx, strUserId, curId).Result(); err != nil {
			log.Println("方法UnFollowAction RdbFollowers 删除value失败: ", err)
		} else { // 如果redis内成功完成了取关操作，则对mysql进行取关操作
			if err := models.UpdateFollowRelation(userId, curId, 2); err != nil {
				log.Println("方法UnFollowAction 取消关注操作失败: ", err)
				return err
			}
		}
	}

	// 对RdbFollowing进行系列操作
	n1, err := redis.RdbFollowings.Exists(redis.Ctx, strCurId).Result()
	if err != nil {
		log.Println("方法UnFollowAction 查询RdbFollowings失败: ", err)
		return err
	}
	if n1 > 0 { //如果存在键值
		if _, err := redis.RdbFollowings.SRem(redis.Ctx, strCurId, userId).Result(); err != nil {
			log.Println("方法UnFollowAction RdbFollowings删除value值失败: ", err)
			return err
		}
	} else { // 如果不存在键值
		// 为key设置默认值, 默认值为自己的id
		if _, err := redis.RdbFollowings.SAdd(redis.Ctx, strCurId, userId).Result(); err != nil {
			log.Println("方法UnFollowAction RdbFollowings 添加默认值失败: ", err)
			return err
		}
		// 为key更新过期时间
		if _, err := redis.RdbFollowings.Expire(redis.Ctx, strCurId, time.Duration(2592000)*time.Second).Result(); err != nil {
			log.Println("方法UnFollowAction RdbFollowings 为key更新默认值失败: ", err)
			return err
		}
		// 查询mysql
		followingIdList, err := models.GetFollowingsId(curId)
		if err != nil {
			log.Println("方法UnFollowAction查询mysql following记录失败")
			return err
		}
		// 维护redis
		for _, followingId := range followingIdList {
			if _, err := redis.RdbFollowings.SAdd(redis.Ctx, strCurId, followingId).Result(); err != nil {
				log.Println("方法UnFollowAction RdbFollowings 维护redis失败")
				return err
			}
		}
		// 维护完redis后，进行最新的取消关注操作
		if _, err := redis.RdbFollowings.SRem(redis.Ctx, strCurId, userId).Result(); err != nil {
			log.Println("方法UnFollowAction RdbFollowings 删除vlaue失败: ", err)
			return err
		}

	}

	return nil
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
		log.Println("用户关注列表：", strFollowingsIdList)
		for _, strFollowingId := range strFollowingsIdList {
			followingId, _ := strconv.ParseInt(strFollowingId, 10, 64)
			log.Println("用户关注的用户id: ", followingId)
			if followingId != userId { // 过滤掉自身
				followingsIdList = append(followingsIdList, followingId)
			}
		}
		if err1 != nil {
			log.Println("方法getFollowingsList中redis查询RdbFollowings库的value值失败")
		}
	} else {
		// key不存在
		// 为key值设置默认值, 默认值为该用户的id
		if _, err := redis.RdbFollowings.SAdd(redis.Ctx, strUserId, userId).Result(); err != nil {
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
	idListLen := len(followingsIdList)
	if idListLen == 0 {
		return *followingsList, nil
	}

	log.Println("用户关注的Id列表：", followingsIdList)

	// 创建协程组
	var wg sync.WaitGroup
	wg.Add(idListLen)
	for i := 0; i < idListLen; i++ {
		go func(pos int) {
			defer wg.Done()
			log.Println("当前被查的用户具体id是：", followingsIdList[pos])
			fmtUser, _ := userService.GetFmtUserByIdWithCurId(followingsIdList[pos], userId)
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
			if followerId != userId {
				followersIdList = append(followersIdList, followerId)
			}
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
	idListLen := len(followersIdList)
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
			fmtUser, _ := userService.GetFmtUserByIdWithCurId(followersIdList[pos], userId)
			log.Println("方法getFollowingsList输出fmtUser: ", fmtUser)
			*followersList = append(*followersList, fmtUser)
		}(i)
	}
	wg.Wait()
	// 返回粉丝对象列表。
	return *followersList, nil
}

// GetFollowersIdList 为下面获取好友列表服务
func (fsi *FollowServiceImp) GetFollowersIdList(userId int64) ([]int64, error) {
	strUserId := strconv.FormatInt(userId, 10)
	followersIdList := new([]int64)

	strFollowersIdList, err := redis.RdbFollowers.SMembers(redis.Ctx, strUserId).Result()
	if err != nil {
		log.Println("方法GetFollowersIdList从RdbFollowers获取粉丝id列表失败")
		return *followersIdList, err
	}
	log.Println("方法GetFollowersIdList查询到用户", userId, "的粉丝Id列表：", strFollowersIdList)
	// 将所有的粉丝id加入列表中
	for _, strFollowerId := range strFollowersIdList {
		followerId, _ := strconv.ParseInt(strFollowerId, 10, 64)
		if followerId != userId { // 去掉自己
			*followersIdList = append(*followersIdList, followerId)
		}
	}

	return *followersIdList, nil
}

func (fsi *FollowServiceImp) GetFriendList(userId int64) ([]FmtFriend, error) {
	userService := UserServiceImpl{}

	friendList := new([]FmtFriend)
	followerIdList, _ := fsi.GetFollowersIdList(userId)
	log.Println("方法GetFriendList查询到粉丝列表：", followerIdList)
	for i := 0; i < len(followerIdList); i++ {
		isFriend, _ := fsi.IsFollowing(followerIdList[i], userId)
		if isFriend { // 如果互关过，则是friends
			log.Println("用户", followerIdList[i], "和用户", userId, "是好友")
			// 获取该用户的信息
			fmtFriend := FmtFriend{}
			fmtFriend.FmtUser, _ = userService.GetFmtUserByIdWithCurId(followerIdList[i], userId)
			log.Println("fmtFriend.FmtUser：", fmtFriend.FmtUser)
			message, _ := models.LatestMessage(followerIdList[i], userId)
			if message == (models.Message{}) { //如果没有消息记录
				fmtFriend.Message = ""
				fmtFriend.MsgType = 0
				log.Println("没有查到消息记录！！！！")
			} else {
				// 判断最近一条信息的类型 0 => 当前用户接收的信息， 1 => 当前用户发送的信息
				if message.ReceiverId != userId {
					fmtFriend.MsgType = 1
				} else {
					fmtFriend.MsgType = 0
				}
				fmtFriend.Message = message.MsgContent
			}

			*friendList = append(*friendList, fmtFriend)
		}
	}
	log.Println("方法GetFriendList查询到粉丝信息：", *friendList)
	return *friendList, nil

}
