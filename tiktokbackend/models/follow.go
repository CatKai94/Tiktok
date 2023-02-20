package models

import "log"

// Follow 用户关系结构，对应用户关系表。
type Follow struct {
	Id         int64
	UserId     int64
	FollowerId int64
	Cancel     int8
}

// TableName 设置Follow结构体对应数据库表名。
func (Follow) TableName() string {
	return "follows"
}

// FindRelation 给定当前用户和目标用户id，查询follow表中相应的记录。
func FindRelation(userId int64, curId int64) (*Follow, error) {
	// follow变量用于后续存储数据库查出来的用户关系。
	follow := Follow{}
	//当查询出现错误时，日志打印err msg，并return err.
	err := DB.Where("user_id = ?", userId).
		Where("follower_id = ?", curId).
		Where("cancel = ?", 0).
		Take(&follow).Error

	if err != nil {
		log.Println("查询关注关系时发生错误：", err.Error())
		return &follow, err
	}
	//正常情况，返回取到的值和空err.
	return &follow, nil
}

// GetTotalFollowerCnt 给定当前用户id，查询follow表中该用户的粉丝数。
func GetTotalFollowerCnt(userId int64) (int64, error) {
	// 用于存储当前用户粉丝数的变量
	var cnt int64
	// 当查询出现错误的情况，日志打印err msg，并返回err.
	err := DB.Model(Follow{}).
		Where("user_id = ?", userId).
		Where("cancel = ?", 0).
		Count(&cnt).Error
	if err != nil {
		log.Println("查询用户粉丝总数时发生错误", err.Error())
		return 0, err
	}
	// 正常情况，返回取到的粉丝数。
	return cnt, nil
}

// GetTotalFollowingCnt 给定当前用户id，查询该用户的关注总数。
func GetTotalFollowingCnt(userId int64) (int64, error) {
	var cnt int64
	// 查询出错，日志打印err msg，并return err
	err := DB.Model(Follow{}).
		Where("follower_id = ?", userId).
		Where("cancel = ?", 0).
		Count(&cnt).Error
	if err != nil { // 查询错误
		log.Println("查询用户关注总数时发生错误：", err.Error())
		return 0, err
	}

	// 查询成功，返回人数。
	return cnt, nil
}

// InsertFollowRelation 给定当前用户和目标对象id，插入其关注关系。  当前用户关注目标用户
func InsertFollowRelation(userId int64, curId int64) (bool, error) {
	// 生成需要插入的关系结构体。
	follow := Follow{
		UserId:     userId,
		FollowerId: curId,
		Cancel:     0,
	}
	// 插入失败，返回err.
	err := DB.Select("UserId", "FollowerId", "Cancel").Create(&follow).Error
	if err != nil {
		log.Println("插入关注记录时发生错误：", err.Error())
		return false, err
	}
	// 插入成功
	return true, nil
}

// GetFollowingsId 给定用户id，查询该用户所有关注者的id。
func GetFollowingsId(userId int64) ([]int64, error) {
	var follwingsId []int64
	err := DB.Model(Follow{}).
		Where("follower_id = ?", userId).
		Pluck("user_id", &follwingsId).Error
	if err != nil {
		log.Println("查询用户的关注用户id列表是发生错误：", err.Error())
		return nil, err
	}
	// 查询成功。
	return follwingsId, nil
}

// GetFollowersId 给定用户id，查询该用户所有的粉丝id
func GetFollowersId(userId int64) ([]int64, error) {
	var followersId []int64
	err := DB.Model(Follow{}).
		Where("user_id = ?", userId).
		Where("cancel = ?", 0).
		Pluck("follower_id", &followersId).Error

	if err != nil {
		//// 没有粉丝，但是不能算错。
		//if "record not found" == err.Error() {
		//	return nil, nil
		//}
		// 查询出错。
		log.Println("查询用户粉丝id时发生错误：", err.Error())
		return nil, err
	}

	return followersId, nil
}

// FindEverFollowing 给定当前用户和目标用户id，查看曾经是否有关注关系。
func FindEverFollowing(userId int64, curId int64) (*Follow, error) {
	// 用于存储查出来的关注关系。
	follow := Follow{}
	//当查询出现错误时，日志打印err msg，并return err.
	err := DB.Where("user_id = ?", userId).
		Where("follower_id = ?", curId).
		Where("cancel = ? or cancel = ?", 0, 1).
		Take(&follow).Error
	if err != nil {
		// 当没查到记录报错时，不当做错误处理。
		if "record not found" == err.Error() {
			return nil, nil
		}
		log.Println("查询曾经是否有关注关系时发生错误：", err.Error())
		return nil, err
	}
	//正常情况，返回取到的关系和空err.
	return &follow, nil
}

// UpdateFollowRelation 给定用户和目标用户的id，更新他们的关系为取消关注或再次关注。
func UpdateFollowRelation(userId int64, targetId int64, cancel int8) (bool, error) {
	// 更新失败，返回错误。
	err := DB.Model(Follow{}).
		Where("user_id = ?", userId).
		Where("follower_id = ?", targetId).
		Update("cancel", cancel).Error

	if err != nil {
		// 更新失败，打印错误日志。
		log.Println("更新关系为取消关注或者再次关注时发生错误： ", err.Error())
		return false, err
	}
	// 更新成功。
	return true, nil
}
