package config

// jwt令牌的Secret
var Secret = "WanderingEarth"

// 局域网ip
var IpUrl = "http://192.168.2.4:3000"

// redis地址
var RedisUrl = "127.0.0.1"

// Feed流每次获得的视频数量
var VideoCount = 5

// FmtUser相关
var DefaultAvatar = IpUrl + "/static/images/IronMan.jpg"
var DefaultBGI = IpUrl + "/static/images/background.jpg"
var DefaultSign = "活着就是为了改变世界"
