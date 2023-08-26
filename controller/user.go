package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"net/http"
	"time"
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

	//user.Id = userIdSequence
	db.Create(&user)
	tokenString, _ := GenToken(user)
	user.Token = tokenString

	c.JSON(http.StatusOK, UserLoginResponse{
		Response: Response{StatusCode: 0},
		UserId:   user.Id,
		Token:    tokenString,
	})
}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	//token := username + password

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
			// 生成Token
			tokenString, _ := GenToken(user)
			/*c.JSON(http.StatusOK, gin.H{
				"code": 2000,
				"msg":  "success",
				"data": gin.H{"token": tokenString},
			})*/
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: Response{StatusCode: 0},
				UserId:   user.Id,
				Token:    tokenString,
			})
		} else {
			//密码错误
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: Response{StatusCode: 1, StatusMsg: "Password wrong"},
			})
			return
		}

	}
}

func UserInfo(c *gin.Context) {
	//token := c.Query("token")
	user_id := c.Query("user_id")
	var user User
	//var user2 User
	db.Where("id = ?", user_id).First(&user)
	//db.Where("id = ?", user_id).First(&user2)
	if user.FollowCount > 0 {
		var ids []int64
		json.Unmarshal([]byte(user.Follows), &ids)
		//db.Where("id IN (?)", ids).Preload("Follows").Find(&users)
		user.IsFollow = false
		db.Save(&user)

	}
	//fmt.Println("user: ", user)
	//fmt.Println("user2: ", user2)
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

}

// CustomClaims 自定义声明类型 并内嵌jwt.RegisteredClaims
// jwt包自带的jwt.RegisteredClaims只包含了官方字段
// 假设我们这里需要额外记录一个username字段，所以要自定义结构体
// 如果想要保存更多信息，都可以添加到这个结构体中
type CustomClaims struct {
	// 可根据需要自行添加字段
	//Username             string `json:"username"`
	Name                 string `json:"name,omitempty"`
	Id                   int64  `json:"id,omitempty" gorm:"primary_key;auto_increment"`
	jwt.RegisteredClaims        // 内嵌标准的声明
}

const TokenExpireDuration = time.Hour * 24

// CustomSecret 用于加盐的字符串
var CustomSecret = []byte("夏天夏天悄悄过去")

// GenToken 生成JWT
func GenToken(user User) (string, error) {
	// 创建一个我们自己的声明
	claims := CustomClaims{
		user.Name,
		user.Id,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpireDuration)),
			Issuer:    "my-project", // 签发人
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(CustomSecret)
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (*CustomClaims, error) {
	// 解析token
	// 如果是自定义Claim结构体则需要使用 ParseWithClaims 方法
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		// 直接使用标准的Claim则可以直接使用Parse方法
		//token, err := jwt.Parse(tokenString, func(token *jwt.Token) (i interface{}, err error) {
		return CustomSecret, nil
	})
	if err != nil {
		return nil, err
	}
	// 对token对象中的Claim进行类型断言
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid { // 校验token
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

// JWTAuthMiddleware 基于JWT的认证中间件
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URI
		// 这里假设Token放在Header的Authorization中，并使用Bearer开头
		// 这里的具体实现方式要依据你的实际业务情况决定
		token := c.Query("token")
		if token == "" {
			token = c.PostForm("token")
		}
		if token == "" {
			// 无token也允许通过
			c.Next()
			return
		}

		// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
		mc, err := ParseToken(token)
		if err != nil {
			fmt.Println("无效的Token")
			c.JSON(http.StatusOK, gin.H{
				"code": 2005,
				"msg":  "无效的Token",
			})
			c.Abort()
			return
		}
		// 将当前请求的username信息保存到请求的上下文c上
		c.Set("Name", mc.Name)
		c.Set("Id", mc.Id)
		fmt.Println("----------------通过鉴权------------------")
		fmt.Println("name: ", mc.Name, " id = ", mc.Id)

		c.Next() // 后续的处理函数可以用过c.Get("username")来获取当前请求的用户信息
	}
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
