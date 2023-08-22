package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"net/http"
)

// usersLoginInfo use map to store user info, and key is username+password for demo
// user data will be cleared every time the server starts
// test data: username=zhanglei, password=douyin
var usersLoginInfo = map[string]User{
	"zhangleidouyin": {
		Id:            1,
		Name:          "zhanglei",
		FollowCount:   10,
		FollowerCount: 5,
		IsFollow:      true,
	},
}

var userIdSequence = int64(1)

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User User `json:"user"`
}

func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	//token := username + password

	//db, err := gorm.Open("mysql", "root:123456@(127.0.0.1:3306)/db1?charset=utf8mb4&parseTime=True&loc=Local")
	//if err != nil {
	//	panic(err)
	//}
	//defer db.Close()

	// 自动迁移
	db.AutoMigrate(&User{})

	var user User
	db.Where("name = ?", username).First(&user)
	if user.Id != 0 {
		// 用户已存在
		c.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "User already exists",
		})
		return
	}

	// 创建新用户
	//atomic.AddInt64(&userIdSequence, 1)

	user.Name = username
	user.Password = password
	user.Token = username + password
	//user.Id = userIdSequence
	db.Create(&user)

	c.JSON(http.StatusOK, UserLoginResponse{
		Response: Response{StatusCode: 0},
		UserId:   user.Id,
		Token:    username + password,
	})

	//if _, exist := usersLoginInfo[token]; exist {
	//	c.JSON(http.StatusOK, UserLoginResponse{
	//		Response: Response{StatusCode: 1, StatusMsg: "User already exist"},
	//	})
	//} else {
	//	atomic.AddInt64(&userIdSequence, 1)
	//	newUser := User{
	//		Id:   userIdSequence,
	//		Name: username,
	//	}
	//	usersLoginInfo[token] = newUser
	//	c.JSON(http.StatusOK, UserLoginResponse{
	//		Response: Response{StatusCode: 0},
	//		UserId:   userIdSequence,
	//		Token:    username + password,
	//	})
	//}
}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	token := username + password

	//db, err := gorm.Open("mysql", "root:123456@(127.0.0.1:3306)/db1?charset=utf8mb4&parseTime=True&loc=Local")
	//if err != nil {
	//	panic(err)
	//}
	//defer db.Close()

	var user User
	fmt.Println("here1")
	db.Where("name = ?", username).First(&user)
	fmt.Println("here2")
	if user.Id == 0 {
		// 用户不存在
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
		return
	} else {
		//用户存在
		if password == user.Password {
			//密码正确
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: Response{StatusCode: 0},
				UserId:   user.Id,
				Token:    token,
			})
		} else {
			//密码错误
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: Response{StatusCode: 1, StatusMsg: "Password wrong"},
			})
			return
		}

	}

	//if user, exist := usersLoginInfo[token]; exist {
	//	c.JSON(http.StatusOK, UserLoginResponse{
	//		Response: Response{StatusCode: 0},
	//		UserId:   user.Id,
	//		Token:    token,
	//	})
	//} else {
	//	c.JSON(http.StatusOK, UserLoginResponse{
	//		Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
	//	})
	//}
}

func UserInfo(c *gin.Context) {
	token := c.Query("token")
	var user User
	db.Where("token = ?", token).First(&user)
	if user.Id != 0 {
		// 用户存在
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 0},
			User:     user,
		})
		return
	} else {
		//用户不存在
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	}

	//if user, exist := usersLoginInfo[token]; exist {
	//	c.JSON(http.StatusOK, UserResponse{
	//		Response: Response{StatusCode: 0},
	//		User:     user,
	//	})
	//} else {
	//	c.JSON(http.StatusOK, UserResponse{
	//		Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
	//	})
	//}
}

func (u *User) DeleteFavoritedVideo(removeID int64) {
	// 1. 反序列化
	var ids []int64
	json.Unmarshal([]byte(u.FavoritedVideos), &ids)

	// 2. 删除id
	for i, id := range ids {
		if id == removeID {
			ids = append(ids[:i], ids[i+1:]...)
			break
		}
	}
	// 3. 序列化并保存
	bytes, _ := json.Marshal(ids)
	u.FavoritedVideos = string(bytes)

	db.Save(&u)
}

func (u *User) AddFavoritedVideo(addID int64) {
	var ids []int64
	if err := json.Unmarshal([]byte(u.FavoritedVideos), &ids); err != nil {
		// 处理错误
	}
	fmt.Println("ids=", ids)
	// 追加新id
	ids = append(ids, addID)
	fmt.Println("ids=", ids)
	// 再次序列化为JSON字符串赋值给字段
	bytes, err := json.Marshal(ids)
	if err != nil {
		// 处理错误
	}
	u.FavoritedVideos = string(bytes)

	db.Save(&u) // 写入JSON字符串
}

func (u *User) AddTotalFavorited(n int64) {
	//user.FavoriteCount += n(n may be 1 or -1)
	u.FavoriteCount += n
	db.Save(&u) // 写入JSON字符串
}

func (u *User) AddWorkCount(n int64) {
	//user.FavoriteCount += n(n may be 1 or -1)
	u.WorkCount += n
	db.Save(&u) // 写入JSON字符串
}
