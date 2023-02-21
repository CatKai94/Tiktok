package service

import (
	"log"
	"sort"
	"strconv"
	"sync"
	"tiktokbackend/middleware/redis"
	"tiktokbackend/models"
	"time"
)

type CommentServiceImpl struct {
	UserService
}

func (c CommentServiceImpl) CountFromVideoId(videoId int64) (int64, error) {
	// 先在缓存中查
	count, err := redis.RdbVideoToCommentId.SCard(redis.Ctx, strconv.FormatInt(videoId, 10)).Result()
	if err != nil {
		log.Println("查询redis出错 （CountFromVideoId）:", err)
	}
	log.Println("redis查询到的评论数量:", count)
	// 返回数量值-1（去除0 值）
	if count != 0 {
		return count - 1, nil
	}
	// redis中没有数据库查
	DaoCount, err1 := models.Count(videoId)
	log.Println("comment count dao :", DaoCount)
	if err1 != nil {
		log.Println("comment count dao err:", err1)
		return 0, nil
	}
	// 将评论id切片存入redis
	go func() {
		// 查询评论id list
		commentList, _ := models.CommentIdList(videoId)
		// 先在redis中存储一个-1值，防止脏读
		_, _err := redis.RdbVideoToCommentId.SAdd(redis.Ctx, strconv.Itoa(int(videoId)), -1).Result()
		if _err != nil {
			log.Println("存储redis失败 CountFromVideoId函数")
			return
		}
		// 设置key值过期时间
		_, err := redis.RdbVideoToCommentId.Expire(redis.Ctx, strconv.Itoa(int(videoId)),
			// 半个月
			time.Duration(60*60*24*15)*time.Second).Result()
		if err != nil {
			log.Println("key值过期时间设置失败")
		}
		// 评论id存入redis
		for _, commentId := range commentList {
			insertRedisVideoCommentId(strconv.Itoa(int(videoId)), commentId)
		}
		log.Println("redis存入成功：CountFromVideoId")
	}()
	return DaoCount, nil
}
func (c CommentServiceImpl) SendComment(comment models.Comment) (CommentInfo, error) {
	comment.Cancel = 0
	insertComment, err := models.InsertComment(comment)
	if err != nil {
		return CommentInfo{}, err
	}
	// 查询用户信息
	impl := UserServiceImpl{}
	user, err := impl.GetFmtUserByIdWithCurId(comment.UserId, comment.UserId)
	// 返回信息
	commentInfo := CommentInfo{
		CommentId:   insertComment.Id,
		UserInfo:    user,
		Content:     insertComment.CommentText,
		PublishDate: insertComment.CreateDate.Format("2006-01-02 15:04:05"),
	}
	// 存入redis
	go func() {
		insertRedisVideoCommentId(strconv.Itoa(int(comment.VideoId)), strconv.Itoa(int(insertComment.Id)))
	}()
	return commentInfo, nil
}
func (c CommentServiceImpl) DeleteComment(commentId int64) error {
	// 1.先查询redis，若有则删除，返回客户端-再go协程删除数据库；无则在数据库中删除，返回客户端。
	count, _ := redis.RdbCommentToVideoId.Exists(redis.Ctx, strconv.FormatInt(commentId, 10)).Result()
	if count > 0 {
		// 在缓存中有此值，则找出来删除，然后返回
		vid, _ := redis.RdbCommentToVideoId.Get(redis.Ctx, strconv.FormatInt(commentId, 10)).Result()
		// 删除，两个redis都要删除
		del1, _ := redis.RdbCommentToVideoId.Del(redis.Ctx, strconv.FormatInt(commentId, 10)).Result()
		del2, _ := redis.RdbVideoToCommentId.SRem(redis.Ctx, vid, strconv.FormatInt(commentId, 10)).Result()
		log.Println("redis删除评论成功", del1, del2) // del1、del2代表删除了几条数据
	}
	err := models.DeleteComment(commentId)
	return err
}
func (c CommentServiceImpl) GetCommentsList(videoId int64, userId int64) ([]CommentInfo, error) {
	// 查询评论列表信息
	list, err := models.GetCommentList(videoId)
	// 组装
	commentInfoList := make([]CommentInfo, len(list))

	// 创建协程组
	var wg sync.WaitGroup
	wg.Add(len(list))
	// idx := 0
	var userService UserServiceImpl

	for i := 0; i < len(list); i++ {
		go func(pos int) {
			defer wg.Done()
			comment := list[pos]
			var commentInfo CommentInfo
			// 根据当前用户id和目标用户id，查询目标用户信息
			var err error

			commentInfo.CommentId = comment.Id
			commentInfo.Content = comment.CommentText
			commentInfo.PublishDate = comment.CreateDate.Format("2026-01-02 15:04:05")
			// userId查询用户信息
			commentInfo.UserInfo, err = userService.GetFmtUserByIdWithCurId(comment.UserId, userId)
			if err != nil {
				log.Println("CommentService-GetList: GetUserByIdWithCurId return err: " + err.Error()) // 函数返回提示错误信息
			}
			commentInfoList[pos] = commentInfo
		}(i)
	}

	wg.Wait()
	// 排序
	sort.Sort(CommentSlice(commentInfoList))

	// 协程查询redis中是否有此记录
	go func() {
		// 缓存中查此视频是否已有评论列表
		cnt, _ := redis.RdbVideoToCommentId.SCard(redis.Ctx, strconv.FormatInt(videoId, 10)).Result()
		// 数据正常，不用更新缓存
		if cnt > 0 {
			return
		}
		// 数据不正确，更新缓存：
		// 先在redis中存储一个-1 值，防止脏读
		_, _err := redis.RdbVideoToCommentId.SAdd(redis.Ctx, strconv.Itoa(int(videoId)), -1).Result()
		if _err != nil { // 若存储redis失败，则直接返回
			log.Println("存储redis失败")
			return
		}
		// key过期时间设置
		_, err2 := redis.RdbVideoToCommentId.Expire(redis.Ctx, strconv.Itoa(int(videoId)),
			// 半个月
			time.Duration(60*60*24*15)*time.Second).Result()
		if err2 != nil {
			log.Println("设置key值过期时间失败")
		}
		// 将评论id存入redis
		for _, comment1 := range commentInfoList {
			insertRedisVideoCommentId(strconv.Itoa(int(videoId)), strconv.Itoa(int(comment1.CommentId)))
		}
		log.Println("评论id存入redis")
	}()
	return commentInfoList, err
}

// 在redis中存储video_id对应的comment_id 、 comment_id对应的video_id
func insertRedisVideoCommentId(videoId string, commentId string) {
	// 在redis-RdbVCid中存储video_id对应的comment_id
	_, err := redis.RdbVideoToCommentId.SAdd(redis.Ctx, videoId, commentId).Result()
	if err != nil {
		// 若存储redis失败-有err，则直接删除key
		redis.RdbVideoToCommentId.Del(redis.Ctx, videoId)
		return
	}
	// 在redis-RdbCVid中存储comment_id对应的video_id
	_, err = redis.RdbCommentToVideoId.Set(redis.Ctx, commentId, videoId, 0).Result()
	if err != nil {
		log.Println("存储评论对应的视频ID失败")
	}
}

// CommentSlice 此变量以及以下三个函数都是做排序-准备工作
type CommentSlice []CommentInfo

func (a CommentSlice) Len() int { // 重写Len()方法
	return len(a)
}
func (a CommentSlice) Swap(i, j int) { // 重写Swap()方法
	a[i], a[j] = a[j], a[i]
}
func (a CommentSlice) Less(i, j int) bool { // 重写Less()方法
	return a[i].CommentId > a[j].CommentId
}
