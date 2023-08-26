package controller

import (
	"encoding/json"
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
	fmt.Println("------------------into Feed()---------------------")
	//db.LogMode(true)
	//fmt.Printf("%+v\n", Video{})
	//fmt.Printf("%+v\n", User{})
	var videos []Video
	//DB.Debug().Model(&User{}).Preload("Profile").Find(&users)
	//db.Debug().Model(&Video{}).Preload("Author").Find(&videos)
	db.Debug().Model(&videos).Preload("Author").Find(&videos)

	token := c.Query("token")
	var user User
	db.Where("token = ?", token).First(&user)

	for _, v := range videos {
		v.Author.IsFollow = false
		v.IsFavorite = false
		db.Save(&v)
	}
	fmt.Println("feed-------all videos = false--------------")

	if user.Id != 0 { //用户已登录
		fmt.Println("------------用户已登录")
		if user.FollowCount > 0 {
			var ids []int64
			json.Unmarshal([]byte(user.Follows), &ids)
			//db.Where("id IN (?)", ids).Preload("Follows").Find(&users)
			/*for _, v := range videos {

				v.Author.IsFollow = false
				fmt.Println("feed-------video:", v, ".Author.IsFollow = false")
				db.Save(&v)

			}*/
			for _, v := range videos {
				for _, i := range ids {
					if v.Author.Id == i {
						fmt.Println("feed-------video:", v, ".Author.IsFollow = true")
						v.Author.IsFollow = true
						db.Save(&v)
						break
					}
				}
			}
		} /*else {
			for _, v := range videos {
				v.Author.IsFollow = false
				db.Save(&v)
				fmt.Println("feed-------video:", v, ".Author.IsFollow = false")
			}
		}*/
		if user.FavoriteCount > 0 {
			var ids []int64
			json.Unmarshal([]byte(user.FavoritedVideos), &ids)
			//db.Where("id IN (?)", ids).Preload("Follows").Find(&users)
			/*for _, v := range videos {
				v.IsFavorite = false
				fmt.Println("feed-------video:", v, ".IsFavorite = false")
				db.Save(&v)
			}*/
			for _, v := range videos {
				for _, i := range ids {
					if v.Id == i {
						//fmt.Println("v.Aurhoe.Id = ", v.Author.Id, " i = ", i)
						v.IsFavorite = true
						db.Save(&v)
						fmt.Println("feed-------video:", v, ".IsFavorite = true")
						break
					}
				}
			}
		} /*else {
			for _, v := range videos {
				v.IsFavorite = false
				db.Save(&v)
				fmt.Println("feed-------video:", v, ".IsFavorite = false")
			}
		}*/
		//db.Where("publish_time < ?", latest_time).Find(&videos)
	} /*else { //游客状态
		for _, v := range videos {
			v.Author.IsFollow = false
			fmt.Println("游客-------feed-------video:", v, ".Author.IsFollow = false")
			db.Save(&v)

		}
	}*/

	/*for _, v := range videos {
		fmt.Println(v)
		fmt.Println("video favorite = ", v.FavoriteCount)
	}*/

	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0, StatusMsg: "succeed"},
		VideoList: videos,
		//VideoList: DemoVideos,
		NextTime: time.Now().Unix(),
	})

}
