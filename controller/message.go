package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

var tempChat = map[string][]Message{}

var messageIdSequence = int64(1)

type ChatResponse struct {
	Response
	MessageList []Message `json:"message_list"`
}

// MessageAction no practical effect, just check if token is valid
func MessageAction(c *gin.Context) {
	//token := c.Query("token")
	to_user_id := c.Query("to_user_id")
	action_type := c.Query("action_type")
	content := c.Query("content")

	var to_user User
	id, ok := c.Get("Id")
	if ok {
		id = id.(int64)
	}

	var user User
	db.Where("Id = ?", id).First(&user)
	//db.Where("token = ?", token).First(&user)
	db.Where("id = ?", to_user_id).First(&to_user)
	if user.Id == 0 || to_user.Id == 0 {
		//用户不存在
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}

	if action_type == "1" {
		// 发送消息
		var msg Message
		msg.FromUserId = user.Id
		msg.ToUserId = to_user.Id
		msg.Content = content
		//msg.CreateTime = time.Now().Format("2006-01-02 15:04:05")
		msg.CreateTime = time.Now().Unix()
		db.Create(&msg)

		c.JSON(http.StatusOK, Response{StatusCode: 0})
	} else {
		//invalid action_type
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "invalid action_type"})
		return
	}
}

// MessageChat all users have same follow list
func MessageChat(c *gin.Context) {
	fmt.Println("----------MessageChat----------")
	//token := c.Query("token")
	to_user_id := c.Query("to_user_id")
	pre_msg_time := c.Query("pre_msg_time")

	var to_user User
	id, ok := c.Get("Id")
	if ok {
		id = id.(int64)
	}

	var user User
	db.Where("Id = ?", id).First(&user)

	//db.Where("token = ?", token).First(&user)
	db.Where("id = ?", to_user_id).First(&to_user)

	fmt.Println("user_id:", user.Id, "to_id:", to_user.Id)

	if user.Id == 0 || to_user.Id == 0 {
		//用户不存在
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}
	var msgs []Message
	db.Where("(from_user_id = ? and to_user_id = ?  and create_time > ?) or (from_user_id = ? and to_user_id = ?  and create_time > ?)", user.Id, to_user.Id, pre_msg_time, to_user.Id, user.Id, pre_msg_time).Find(&msgs)
	//db.Where("create_time > ?", pre_msg_time).Find(&msgs)
	//fmt.Println("msgs: ", msgs)
	c.JSON(http.StatusOK, ChatResponse{Response: Response{StatusCode: 0, StatusMsg: "get msgs"}, MessageList: msgs})
}

func genChatKey(userIdA int64, userIdB int64) string {
	if userIdA > userIdB {
		return fmt.Sprintf("%d_%d", userIdB, userIdA)
	}
	return fmt.Sprintf("%d_%d", userIdA, userIdB)
}
