package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type FeedResponse struct {
	Response
	VideoList []Video `json:"video_list,omitempty"`
	NextTime  int64   `json:"next_time,omitempty"`
}

// Feed same demo video list for every request
func Feed(c *gin.Context) {
	//latest_time := c.DefaultQuery("latest_time", fmt.Sprintf("%d", (time.Now().Unix())))
	//latime, _ := strconv.Atoi(latest_time)
	//latetime := time.Unix(int64(latime), 0)
	//var videos []Video
	db.LogMode(true)
	fmt.Printf("%+v\n", Video{})
	fmt.Printf("%+v\n", User{})
	var videos []Video
	//DB.Debug().Model(&User{}).Preload("Profile").Find(&users)
	//db.Debug().Model(&Video{}).Preload("Author").Find(&videos)
	db.Debug().Model(&videos).Preload("Author").Find(&videos)
	//db.Where("publish_time < ?", latest_time).Find(&videos)

	for _, v := range videos {
		fmt.Println(v)
		fmt.Println("video favorite = ", v.FavoriteCount)
	}

	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0, StatusMsg: "succeed"},
		VideoList: videos,
		//VideoList: DemoVideos,
		NextTime: time.Now().Unix(),
	})

}
