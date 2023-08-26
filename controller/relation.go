package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserListResponse struct {
	Response
	UserList []User `json:"user_list"`
}

// RelationAction no practical effect, just check if token is valid
func RelationAction(c *gin.Context) {
	//token := c.Query("token")
	to_user_id := c.Query("to_user_id")
	action_type := c.Query("action_type")
	id, ok := c.Get("Id")
	if ok {
		id = id.(int64)
	}

	var user User
	db.Where("Id = ?", id).First(&user)
	var to_user User
	//db.Where("token = ?", token).First(&user)
	db.Where("id = ?", to_user_id).First(&to_user)

	if user.Id == 0 || to_user.Id == 0 {
		//用户不存在
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}

	if action_type == "1" {
		// 关注
		user.AddFollows(to_user.Id)
		to_user.AddFollowers(user.Id)
		c.JSON(http.StatusOK, Response{StatusCode: 0})
	} else if action_type == "2" {
		user.DeleteFollows(to_user.Id)
		to_user.DeleteFollowers(user.Id)
		c.JSON(http.StatusOK, Response{StatusCode: 0})
	} else {
		//invalid action_type
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "invalid action_type"})
		return
	}

}

// FollowList all users have same follow list
func FollowList(c *gin.Context) {
	user_id := c.Query("user_id")

	var user User

	//db.Where("token = ?", token).First(&user)
	db.Where("id = ?", user_id).First(&user)
	if user.Id == 0 {
		//用户不存在
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}

	var users []User

	if user.FollowCount > 0 {
		var ids []int64
		json.Unmarshal([]byte(user.Follows), &ids)
		db.Where("id IN (?)", ids).Find(&users)
		for _, u := range users {
			u.IsFollow = true
			db.Save(&u)
		}
		fmt.Println("------------relation----------update IsFollow = true: ", users)
	} else {
		// 数组为空时,直接返回空视频数组
		users = []User{}
	}
	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: users,
	})
}

// FollowerList all users have same follower list
func FollowerList(c *gin.Context) {
	user_id := c.Query("user_id")

	var user User

	//db.Where("token = ?", token).First(&user)
	db.Where("id = ?", user_id).First(&user)
	if user.Id == 0 {
		//用户不存在
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}

	var users []User

	if user.FollowerCount > 0 {
		var ids []int64
		json.Unmarshal([]byte(user.Followers), &ids)
		db.Where("id IN (?)", ids).Find(&users)
		var idss []int64
		json.Unmarshal([]byte(user.Follows), &idss)
		for _, u := range users {
			u.IsFollow = false
			db.Save(&u)
			fmt.Println("------------relation----------update IsFollow = false: ", users)
		}
		for _, u := range users {
			for _, i := range idss {
				if u.Id == i {
					u.IsFollow = true
					db.Save(&u)
					fmt.Println("------------relation----------update IsFollow = true: ", u)
					break
				}
			}

		}
	} else {
		// 数组为空时,直接返回空视频数组
		users = []User{}
	}
	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: users,
	})
}

// FriendList all users have same friend list
func FriendList(c *gin.Context) {
	user_id := c.Query("user_id")

	var user User

	//db.Where("token = ?", token).First(&user)
	db.Where("id = ?", user_id).First(&user)
	if user.Id == 0 {
		//用户不存在
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}

	var users []User

	if user.FollowerCount > 0 {
		var ids []int64
		json.Unmarshal([]byte(user.Followers), &ids)
		var idss []int64
		json.Unmarshal([]byte(user.Follows), &idss)
		//fmt.Println("---------friend-----ids, idss = ", ids, idss)
		db.Where("id IN (?) and id IN (?)", ids, idss).Find(&users)
		//fmt.Println("----------users:", users)

		for _, u := range users {

			u.IsFollow = true
			db.Save(&u)
			break

		}
	} else {
		// 数组为空时,直接返回空视频数组
		users = []User{}
	}
	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: users,
	})
	/*c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: []User{DemoUser},
	})*/
}

func (u *User) DeleteFollows(removeID int64) {
	// 1. 反序列化
	var ids []int64
	json.Unmarshal([]byte(u.Follows), &ids)

	// 2. 删除id
	for i, id := range ids {
		if id == removeID {
			ids = append(ids[:i], ids[i+1:]...)
			break
		}
	}
	// 3. 序列化并保存
	bytes, _ := json.Marshal(ids)
	u.Follows = string(bytes)

	u.FollowCount-- //关注数-1
	db.Save(&u)
}

func (u *User) AddFollows(addID int64) {
	var ids []int64
	if err := json.Unmarshal([]byte(u.Follows), &ids); err != nil {
		// 处理错误
	}
	// 追加新id
	ids = append(ids, addID)
	// 再次序列化为JSON字符串赋值给字段
	bytes, err := json.Marshal(ids)
	if err != nil {
		// 处理错误
	}
	u.Follows = string(bytes)

	u.FollowCount++ //关注数+1

	db.Save(&u) // 写入JSON字符串
}

func (u *User) DeleteFollowers(removeID int64) {
	// 1. 反序列化
	var ids []int64
	json.Unmarshal([]byte(u.Followers), &ids)

	// 2. 删除id
	for i, id := range ids {
		if id == removeID {
			ids = append(ids[:i], ids[i+1:]...)
			break
		}
	}
	// 3. 序列化并保存
	bytes, _ := json.Marshal(ids)
	u.Followers = string(bytes)

	u.FollowerCount-- //关注数-1
	db.Save(&u)
}

func (u *User) AddFollowers(addID int64) {
	var ids []int64
	if err := json.Unmarshal([]byte(u.Followers), &ids); err != nil {
		// 处理错误
	}
	// 追加新id
	ids = append(ids, addID)
	// 再次序列化为JSON字符串赋值给字段
	bytes, err := json.Marshal(ids)
	if err != nil {
		// 处理错误
	}
	u.Followers = string(bytes)

	u.FollowerCount++ //关注数+1

	db.Save(&u) // 写入JSON字符串
}
