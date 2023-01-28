package service

import (
	"log"
	"tiktokbackend/config"
	"tiktokbackend/models"
	"time"
)

type VideoServiceImpl struct {
	UserService
}

// Feed
// 通过传入时间戳，当前用户的id，返回对应的视频数组，以及视频数组中最早的发布时间
// 获取视频数组大小是可以控制的，在config中的videoCount变量
func (vsi *VideoServiceImpl) Feed(lastTime time.Time, userId int64) ([]FmtVideo, time.Time, error) {
	//创建对应返回视频的切片数组，提前将切片的容量设置好，可以减少切片扩容的性能
	videosList := make([]FmtVideo, 0, config.VideoCount)

	//根据传入的时间，获得传入时间前n个视频，可以通过config.videoCount来控制
	videos, err := models.GetVideosByLastTime(lastTime)
	log.Println("service层的videos: ", videos[0])
	if err != nil {
		log.Printf("方法models.GetVideosByLastTime(lastTime) 失败：%v", err)
		return nil, time.Time{}, err
	}
	log.Printf("方法models.GetVideosByLastTime(lastTime) 成功：%v", videos)
	//将数据通过copyVideos进行处理，在拷贝的过程中对数据进行组装    // 2023年1.28注释，晚点再写，先用假数据
	err = vsi.copyVideos(&videosList, &videos, userId)
	//if err != nil {
	//	log.Printf("方法videoService.copyVideos(&videos, &tableVideos, userId) 失败：%v", err)
	//	return nil, time.Time{}, err
	//}
	//log.Printf("方法videoService.copyVideos(&videos, &tableVideos, userId) 成功")
	//返回数据，同时获得视频中最早的时间返回
	return videosList, videos[len(videos)-1].PublishTime, nil
}

// 该方法可以将数据进行拷贝和转换，并从其他方法获取对应的数据
func (vsi *VideoServiceImpl) copyVideos(result *[]FmtVideo, data *[]models.Video, userId int64) error {

	//遍历查到的所有Video对象，将它们逐个包装成FmtVideo对象
	for _, temp := range *data {
		var fmtVideo FmtVideo
		//将fmtVideo进行组装，添加想要的信息,插入从数据库中查到的数据
		vsi.creatVideo(&fmtVideo, &temp, userId)
		*result = append(*result, fmtVideo)
	}
	return nil
}

// 将fmtVideo进行组装，添加想要的信息,插入从数据库中查到的数据
func (vsi *VideoServiceImpl) creatVideo(fmtVideo *FmtVideo, data *models.Video, userId int64) {
	////建立协程组，当这一组的携程全部完成后，才会结束本方法
	//var wg sync.WaitGroup
	//wg.Add(4)
	//var err error
	//fmtVideo.Video = *data
	////插入Author，这里需要将视频的发布者和当前登录的用户传入，才能正确获得isFollow，
	////如果出现错误，不能直接返回失败，将默认值返回，保证稳定
	//go func() {
	//	fmtVideo.Author, err = vsi.GetUserByIdWithCurId(data.AuthorId, userId)
	//	if err != nil {
	//		log.Printf("方法videoService.GetUserByIdWithCurId(data.AuthorId, userId) 失败：%v", err)
	//	} else {
	//		log.Printf("方法videoService.GetUserByIdWithCurId(data.AuthorId, userId) 成功")
	//	}
	//	wg.Done()
	//}()
	//
	////插入点赞数量，同上所示，不将nil直接向上返回，数据没有就算了，给一个默认就行了
	//go func() {
	//	fmtVideo.FavoriteCount, err = vsi.FavouriteCount(data.Id)
	//	if err != nil {
	//		log.Printf("方法videoService.FavouriteCount(data.ID) 失败：%v", err)
	//	} else {
	//		log.Printf("方法videoService.FavouriteCount(data.ID) 成功")
	//	}
	//	wg.Done()
	//}()
	//
	////获取该视屏的评论数字
	//go func() {
	//	fmtVideo.CommentCount, err = vsi.CountFromVideoId(data.Id)
	//	if err != nil {
	//		log.Printf("方法videoService.CountFromVideoId(data.ID) 失败：%v", err)
	//	} else {
	//		log.Printf("方法videoService.CountFromVideoId(data.ID) 成功")
	//	}
	//	wg.Done()
	//}()
	//
	////获取当前用户是否点赞了该视频
	//go func() {
	//	fmtVideo.IsFavorite, err = vsi.IsFavourite(video.Id, userId)
	//	if err != nil {
	//		log.Printf("方法videoService.IsFavourit(video.Id, userId) 失败：%v", err)
	//	} else {
	//		log.Printf("方法videoService.IsFavourit(video.Id, userId) 成功")
	//	}
	//	wg.Done()
	//}()
	//
	//wg.Wait()

	//这里我手写了一些假数据
	fmtUser := FmtUser{
		Id:            2,
		Name:          "辛弃疾",
		FollowCount:   17,
		FollowerCount: 32,
		IsFollow:      false,
		TotalFavorite: 245,
		FavoriteCount: 1756,
	}
	fmtVideo.Video = *data
	fmtVideo.Author = fmtUser
	fmtVideo.FavoriteCount = 40
	fmtVideo.CommentCount = 345
	fmtVideo.IsFavorite = false

}
