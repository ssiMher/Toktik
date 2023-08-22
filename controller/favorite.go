package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// FavoriteAction no practical effect, just check if token is valid
func FavoriteAction(c *gin.Context) {
	token := c.Query("token")
	video_id := c.Query("video_id")
	action_type := c.Query("action_type")

	var user User
	db.Where("token = ?", token).First(&user)
	if user.Id == 0 {
		//用户不存在
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}
	var video Video
	db.Where("Id = ?", video_id).First(&video)
	if video.Id == 0 {
		//视频不存在
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Video doesn't exist"})
		return
	}
	if action_type == "1" {
		video.FavoriteCount++
		// 添加视频ID到用户收藏列表
		//db.Model(&user).Update("FavoritedVideos", append(user.FavoritedVideos, video.Id))
		//user.FavoritedVideos = append(user.FavoritedVideos, video.FavoriteCount)
		user.AddFavoritedVideo(video.Id)
		user.AddTotalFavorited(1)
		// 视频收藏数+1
		db.Table("videos").Model(&video).Update("favorite_count", video.FavoriteCount)
		c.JSON(http.StatusOK, Response{StatusCode: 0})
	} else if action_type == "2" {
		video.FavoriteCount--
		// 从用户收藏列表中移除视频ID
		//db.Model(&user).Update("FavoritedVideos", removeVideoID(user.FavoritedVideos, video.Id))
		user.AddTotalFavorited(-1)
		user.DeleteFavoritedVideo(video.Id)

		db.Table("videos").Model(&video).Update("favorite_count", video.FavoriteCount)
		c.JSON(http.StatusOK, Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Invalid action_type"})
		return
	}

	/*if _, exist := usersLoginInfo[token]; exist {
		c.JSON(http.StatusOK, Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}*/
}

// FavoriteList all users have same favorite video list
func FavoriteList(c *gin.Context) {
	//latest_time := c.DefaultQuery("latest_time", fmt.Sprintf("%d", (time.Now().Unix())))
	user_id := c.Query("user_id")
	//TODO 鉴权
	var user User
	db.Where("id = ?", user_id).First(&user)
	var videos []Video
	if len(user.FavoritedVideos) > 0 {
		var ids []int64
		json.Unmarshal([]byte(user.FavoritedVideos), &ids)
		db.Where("id IN (?)", ids).Preload("Author").Find(&videos)
	} else {
		// 数组为空时,直接返回空视频数组
		videos = []Video{}
	}
	//db.Where("id IN ?", user.FavoritedVideos).Find(&videos)

	for _, v := range videos {
		v.IsFavorite = true
	}
	fmt.Println(videos)
	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0},
		VideoList: videos,
		//NextTime:  time.Now().Unix(),
	})
	//c.JSON(http.StatusOK, VideoListResponse{
	//	Response: Response{
	//		StatusCode: 0,
	//	},
	//	VideoList: DemoVideos,
	//})
}
