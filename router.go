package main

import (
	"github.com/RaymondCode/simple-demo/controller"
	"github.com/gin-gonic/gin"
	"io"
	"os"
)

func initRouter(r *gin.Engine) {
	// public directory is used to serve static resources
	r.Static("/static", "./public")

	apiRouter := r.Group("/douyin")

	// basic apis
	apiRouter.GET("/feed/", controller.Feed)
	apiRouter.GET("/user/", controller.UserInfo)
	apiRouter.POST("/user/register/", controller.Register)
	apiRouter.POST("/user/login/", controller.Login)
	apiRouter.POST("/publish/action/", controller.Publish)
	apiRouter.GET("/publish/list/", controller.PublishList)

	r.GET("/video/:name", func(c *gin.Context) {

		name := c.Param("name")

		c.File("E:/code/Go/simple-demo-main/simple-demo-main/public/" + name)

		file, _ := os.Open("E:/code/Go/simple-demo-main/simple-demo-main/public/" + name)

		c.Header("Content-Type", "video/mp4")
		c.Stream(func(w io.Writer) bool {
			io.Copy(w, file)
			return false
		})

	})

	r.GET("/images/:name", func(c *gin.Context) {
		name := c.Param("name")
		//file, _ := os.Open("./images/" + name)
		filepath := "./public/" + name
		//c.Header("Content-Type", "image/png") // 按图片类型设置
		c.File(filepath)
		// 处理请求

	})

	// extra apis - I
	apiRouter.POST("/favorite/action/", controller.FavoriteAction)
	apiRouter.GET("/favorite/list/", controller.FavoriteList)
	apiRouter.POST("/comment/action/", controller.CommentAction)
	apiRouter.GET("/comment/list/", controller.CommentList)

	// extra apis - II
	apiRouter.POST("/relation/action/", controller.RelationAction)
	apiRouter.GET("/relation/follow/list/", controller.FollowList)
	apiRouter.GET("/relation/follower/list/", controller.FollowerList)
	apiRouter.GET("/relation/friend/list/", controller.FriendList)
	apiRouter.GET("/message/chat/", controller.MessageChat)
	apiRouter.POST("/message/action/", controller.MessageAction)
}
