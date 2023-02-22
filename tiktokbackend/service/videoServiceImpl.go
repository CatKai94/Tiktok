package service

import (
	"log"
	"sync"
	"tiktokbackend/config"
	"tiktokbackend/models"
	"time"
)

type VideoServiceImpl struct {
}

// Feed
// 通过传入时间戳，当前用户的id，返回对应的视频数组，以及视频数组中最早的发布时间
// 获取视频数组大小是可以控制的，在config中的videoCount变量
func (vsi *VideoServiceImpl) Feed(lastTime time.Time, userId int64) ([]FmtVideo, time.Time, error) {
	//创建对应返回视频的切片数组，提前将切片的容量设置好，可以减少切片扩容的性能
	videosList := make([]FmtVideo, 0, config.VideoCount)

	//根据传入的时间，获得传入时间前n个视频，可以通过config.videoCount来控制
	videos, err := models.GetVideosByLastTime(lastTime)

	if err != nil {
		log.Println("方法Feed获取按时间排序的视频失败: ", err)
		return nil, time.Time{}, err
	}

	//将models层内的video加工成FmtVideo
	err = vsi.refactorVideos(&videosList, &videos, userId)

	return videosList, videos[len(videos)-1].PublishTime, nil
}

// GetVideo
// 传入视频id获得对应的视频对象，注意还需要传入当前登录用户id
func (vsi *VideoServiceImpl) GetVideo(videoId int64, userId int64) (FmtVideo, error) {
	//初始化video对象
	var fmtVideo FmtVideo

	//从数据库中查询数据，如果查询不到数据，就直接失败返回，后续流程就不需要执行了
	video, err := models.GetVideoByVideoId(videoId)
	if err != nil {
		log.Printf("方法dao.GetVideoByVideoId(videoId) 失败：%v", err)
		return fmtVideo, err
	}

	//插入从数据库中查到的数据
	vsi.creatFmtVideo(&fmtVideo, &video, userId)
	return fmtVideo, nil
}

func (vsi *VideoServiceImpl) Publish(fileName string, userId int64, title string) error {
	//生成视频名称
	videoName := fileName
	imageName := fileName

	err := models.Save(videoName, imageName, userId, title)

	if err != nil {
		log.Println("视频信息入库失败：", err)
		return err
	}

	log.Println("视频信息入库成功")
	return nil
}

func (vsi *VideoServiceImpl) List(userId int64, curId int64) ([]FmtVideo, error) {
	videos, err := models.GetVideosByAuthorId(userId)

	fmtVideoList := make([]FmtVideo, 0, len(videos))
	err = vsi.refactorVideos(&fmtVideoList, &videos, curId)

	if err != nil {
		log.Println("格式化Videos时发生错误：", err)
		return nil, err
	}
	return fmtVideoList, nil
}

// 将Video格式化为fmtVideo
func (vsi *VideoServiceImpl) refactorVideos(result *[]FmtVideo, data *[]models.Video, userId int64) error {
	//遍历查到的所有Video对象，将它们逐个包装成FmtVideo对象
	for _, temp := range *data {
		var fmtVideo FmtVideo
		//将fmtVideo进行组装，添加想要的信息,插入从数据库中查到的数据
		vsi.creatFmtVideo(&fmtVideo, &temp, userId)
		*result = append(*result, fmtVideo)
	}
	return nil
}

// 将fmtVideo进行组装，添加想要的信息,插入从数据库中查到的数据
func (vsi *VideoServiceImpl) creatFmtVideo(fmtVideo *FmtVideo, video *models.Video, userId int64) {

	var likeService LikeServiceImpl
	var userService UserServiceImpl
	var commentService CommentServiceImpl

	//创建协程组
	var wg sync.WaitGroup
	wg.Add(4) // 加上评论后再改成4
	var err error

	fmtVideo.Video = *video

	//插入Author，这里需要将视频的发布者和当前登录的用户传入
	go func() {
		fmtVideo.Author, err = userService.GetFmtUserByIdWithCurId(video.AuthorId, userId)
		if err != nil {
			log.Printf("方法videoService.GetUserByIdWithCurId(data.AuthorId, userId) 失败：%v", err)
		} else {
			log.Printf("方法videoService.GetUserByIdWithCurId(data.AuthorId, userId) 成功")
		}
		wg.Done()
	}()

	//获取视频的点赞总数
	go func() {
		fmtVideo.FavoriteCount, err = likeService.GetLikeCount(video.Id)
		if err != nil {
			log.Printf("方法videoService.FavouriteCount(data.ID) 失败：%v", err)
		} else {
			log.Printf("方法videoService.FavouriteCount(data.ID) 成功")
		}
		wg.Done()
	}()

	//获取当前用户是否点赞了该视频
	go func() {
		fmtVideo.IsFavorite, err = likeService.JudgeUserIsLike(userId, video.Id)
		log.Println("当前用户是否点赞了该视频，", fmtVideo.IsFavorite)
		if err != nil {
			log.Printf("方法videoService.IsFavourit(video.Id, userId) 失败：%v", err)
		} else {
			log.Printf("方法videoService.IsFavourit(video.Id, userId) 成功")
		}
		wg.Done()
	}()

	//获取该视频的评论总数
	go func() {
		fmtVideo.CommentCount, err = commentService.CountFromVideoId(video.Id)
		if err != nil {
			log.Printf("方法videoService.CountFromVideoId(data.ID) 失败：%v", err)
		} else {
			log.Printf("方法videoService.CountFromVideoId(data.ID) 成功")
		}
		wg.Done()
	}()

	wg.Wait()
}

// GetFmtVideo 通过视频id获得FmtVideo对象
func (vsi *VideoServiceImpl) GetFmtVideo(videoId int64, userId int64) (FmtVideo, error) {
	var fmtVideo FmtVideo

	video, err := models.GetVideoByVideoId(videoId)
	if err != nil {
		log.Println("查询视频失败")
		return fmtVideo, err
	}

	vsi.creatFmtVideo(&fmtVideo, &video, userId)
	return fmtVideo, nil
}

func (vsi *VideoServiceImpl) GetVideoCntByUserId(userId int64) int64 {
	ids, err := models.GetVideosByAuthorId(userId)
	if err != nil {
		log.Println("发生了错误, ", err)
	}

	return int64(len(ids))
}
