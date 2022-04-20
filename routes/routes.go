package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mikecodesun/backend-sql/controllers"
	middleWare "github.com/mikecodesun/backend-sql/middleware"
)

func UserRoute(r *gin.Engine) {
	r.POST("/users/register", controllers.Register())
	r.POST("/users/login", controllers.Login())
	r.GET("/users", middleWare.Auth(), controllers.GetAllUsers())
	r.GET("/users/:id",middleWare.Auth(), controllers.GetUserById())
	r.GET("/users/logout", controllers.Logout())
}

func PostRoute(r *gin.Engine) {
	r.POST("/posts",middleWare.Auth(), controllers.CreatePost())
	r.GET("/posts", controllers.GetAllPosts())
	r.GET("/posts/:id", controllers.GetPostById())
	r.DELETE("/posts/:id",middleWare.Auth(), controllers.DeletePost())
	r.PATCH("/posts/:id",middleWare.Auth(), controllers.UpdatePost())
}