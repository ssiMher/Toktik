package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type VideoListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {
	//token := c.PostForm("token")
	id, ok := c.Get("Id")
	if ok {
		id = id.(int64)
	}

	var user User
	db.Where("Id = ?", id).First(&user)
	fmt.Println("id = ", id, user)
	if user.Id == 0 {
		//用户不存在
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}

	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	//filename := filepath.Base(data.Filename)
	video := &Video{}

	filename := c.PostForm("title")
	video.Title = filename
	//user := usersLoginInfo[token]
	finalName := fmt.Sprintf("%d_%s", user.Id, filename)
	finalName += ".mp4"
	//saveFile := filepath.Join("./public/", finalName)
	saveFile := "public/" + finalName
	user.AddWorkCount(1)
	if err = c.SaveUploadedFile(data, saveFile); err != nil {
		fmt.Println("上传视频失败")
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	video.PlayUrl = NgrokHost + "/video/" + finalName
	//video.PlayUrl = saveFile
	video.Author = user

	fmt.Println("publish: user: ", video.Author)

	video.UserID = user.Id
	video.PublishTime = fmt.Sprintf("%d", (time.Now().Unix()))
	//video.CoverUrl = "https://7606-111-49-156-134.ngrok-free.app/code/Go/simple-demo-main/simple-demo-main/public/image.jpg"
	video.CoverUrl = NgrokHost + "/images/image.jpg"

	db.Create(video)
	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  finalName + " uploaded successfully",
	})
}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {
	//token:=c.Query("token")
	id := c.Query("user_id")
	// 查询该用户发布的所有视频
	var videos []Video
	db.Where("user_id = ?", id).Find(&videos)

	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: videos,
	})
}
