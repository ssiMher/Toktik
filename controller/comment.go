package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type CommentListResponse struct {
	Response
	CommentList []Comment `json:"comment_list,omitempty"`
}

type CommentActionResponse struct {
	Response
	Comment Comment `json:"comment,omitempty"`
}

// CommentAction no practical effect, just check if token is valid
func CommentAction(c *gin.Context) {
	token := c.Query("token")
	actionType := c.Query("action_type")
	video_id := c.Query("video_id")
	var video Video
	db.Where("Id = ?", video_id).First(&video)
	if video.Id == 0 {
		//视频不存在
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Video doesn't exist"})
		return
	}
	var user User
	db.Where("token = ?", token).First(&user)
	if user.Id == 0 {
		//用户不存在
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}

	if actionType == "1" {
		text := c.Query("comment_text")
		comment := Comment{}
		comment.User = user
		comment.UserID = user.Id
		comment.Video = video
		comment.VideoID = video.Id
		comment.Content = text
		comment.CreateDate = time.Now().Format("01-02")
		db.Create(&comment)
		video.AddCommentCount(1)
		c.JSON(http.StatusOK, CommentActionResponse{Response: Response{StatusCode: 0},
			Comment: comment})
		/*Comment{
			Id:         1,
			User:       user,
			Content:    text,
			CreateDate: "05-01",
		}*/
		return
	} else {
		comment_id := c.Query("comment_id")
		var comment Comment
		db.Where("Id = ?", comment_id).First(&comment)
		if comment.Id == 0 {
			//评论不存在
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Comment doesn't exist"})
			return
		}
		video.AddCommentCount(-1)
		db.Delete(&comment)

		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) {
	//token := c.Query("token")
	//TODO() 鉴权
	video_id := c.Query("video_id")
	var video Video
	db.Where("Id = ?", video_id).First(&video)
	if video.Id == 0 {
		//视频不存在
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Video doesn't exist"})
		return
	}

	var comments []Comment

	db.Debug().Where("video_id = ?", video.Id).Preload("User").Preload("Video").Find(&comments)

	// 访问预加载的关联
	for _, cm := range comments {
		user := cm.User
		vide := cm.Video
		// 使用预加载的user和video
		fmt.Printf("User:%s Video:%s\n", user.Name, vide.Title)
	}
	c.JSON(http.StatusOK, CommentListResponse{
		Response:    Response{StatusCode: 0},
		CommentList: comments,
	})
}
