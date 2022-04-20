package controllers

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mikecodesun/backend-sql/database"
	"github.com/mikecodesun/backend-sql/models"

	"github.com/mikecodesun/backend-sql/utils"
	"golang.org/x/crypto/bcrypt"
)

// get db
var db = database.ConnectDB()
// init db
func init() {
	// db.AutoMigrate(&models.User{})
	db.Migrator().CreateTable(&models.Post{})
	db.Migrator().CreateTable(&models.User{})
}

var validate = validator.New()

func Register() gin.HandlerFunc{
	return func(c *gin.Context) {
		var user models.User
		var foundUser models.User
		// get req.body user (name, password, email, )
		if err := c.BindJSON(&user); err != nil {
			log.Fatal("Error Json Bind User")
		}
		// validate req.body input
	if err :=	validate.Struct(user); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	// db check user exist?
  db.Where("name = ?", user.Name).Find(&foundUser)
	if foundUser.Name != "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "User name already exist"})
		return
	}
	// hash password
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		log.Fatal("Error hash password")
	}
	user.Password = string(hashPassword)
	// create user
	db.Create(&user)
	// res
  c.JSON(http.StatusOK, user)
	}
}
// 4-19 login
func Login() gin.HandlerFunc{
	return func(c *gin.Context) {
		var user models.User
		var foundUser models.User
		// get req user
		if err := c.BindJSON(&user); err != nil {
			log.Fatal("Error bind json user")
		}
		// find user exist?
		db.Where("name = ?", user.Name).Find(&foundUser)
		if foundUser.Name == "" {
			c.JSON(http.StatusNotFound, gin.H{"msg": "User not exist"})
			return
		}
		// check user input password match hashpassword
		if err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(user.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "password not right"})
			return
		}
		
		// generate jwt token
		token := utils.GenerateJwt(foundUser.Name, foundUser.User_type)
		// set cookie
		c.SetCookie("token", token, 60 * 60 * 24, "/", "localhost", false, true)
		// res token
		c.JSON(http.StatusOK, gin.H{"msg": "User login Successfully"})
	}
}
// log out user , clear cookie
func Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.SetCookie("token", "", 0, "", "", false, true)
		c.JSON(http.StatusOK, gin.H{"msg": "User log out Successfully!"})
	}
}

func GetAllUsers() gin.HandlerFunc{
	return func(c *gin.Context) {
		var users []models.User
		// check user type 
		err := utils.CheckUser(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": err.Error()})
			return
		}
		// find all 
		res := db.Find(&users)
		log.Println(res.RowsAffected)
		if res.RowsAffected < 1 {
			
			c.JSON(http.StatusNotFound, gin.H{"msg": "no users"})
			return
		}
		// send res 
		c.JSON(http.StatusOK, users)
	}
}
func GetUserById() gin.HandlerFunc{
	return func(c *gin.Context) {
		var user models.User
		// check user type
		err := utils.CheckUser(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": err.Error()})
			return
		}
		// get param
		id := c.Param("id")
		// find one
		res := db.Where("id = ?", id).Find(&user)
		if res.RowsAffected < 1 {
			log.Println(res.RowsAffected)
			c.JSON(http.StatusNotFound, gin.H{"msg": "no user exist"})
			return
		}
		// send res
		c.JSON(http.StatusOK, user)
	}
}

func CreatePost() gin.HandlerFunc{
	return func(c *gin.Context) {
		var post models.Post
		// get req json post body
		if err := c.BindJSON(&post); err != nil {
			log.Fatal("Error, post json bind")
		}
		// check update post input body not empty
		if strings.TrimSpace(post.Body) == "" {
			c.JSON(http.StatusNotFound, gin.H{"msg": "Post must not be empty"})
			return
		}
			// get username 
			username := c.GetString("username")
			post.UserName = username
		// validate post
		if err := validate.Struct(post); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
			return
		}
		// db create new post
	  db.Create(&post)
		// res
		c.JSON(http.StatusOK, post)
	}
}

func GetAllPosts() gin.HandlerFunc{
	return func(c *gin.Context) {
		var posts []models.Post
		// find all posts
		res := db.Find(&posts)
		// check posts not null
		if res.RowsAffected < 1 {
			c.JSON(http.StatusNotFound, gin.H{"msg": "No Posts"})
			return
		}
		// res
		c.JSON(http.StatusOK, posts)
	}
}

func DeletePost() gin.HandlerFunc{
	return func(c *gin.Context) {
		var post models.Post
		// get route param id
		id := c.Param("id")
		// find the post need to delete
		res := db.Where("id = ?", id).Find(&post)
		// check post exist?
		if res.RowsAffected < 1 {
			c.JSON(http.StatusNotFound, gin.H{"msg": "No Post"})
			return
		}
		// check user have the auth to delete
		username := c.GetString("username")
		if username != post.UserName {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "Not have auth to delete"})
			return
		}
		// delete post from db
		db.Where("id = ?", id).Delete(&post)
		// res
		c.JSON(http.StatusOK, gin.H{"msg": "psot delete"})

	}
}

func GetPostById() gin.HandlerFunc{
	return func(c *gin.Context) {
		var post models.Post
		// get route param id
		id := c.Param("id")
		// find the post by id
		res := db.Where("id = ?", id).Find(&post)
		// check post exist?
		if res.RowsAffected < 1 {
			c.JSON(http.StatusNotFound, gin.H{"msg": "No Post"})
			return
		}
		c.JSON(http.StatusOK, post)
	}
}
func UpdatePost() gin.HandlerFunc{
	return func(c *gin.Context) {
		var post models.Post
		var updatePost models.Post
		// get req.body post
		if err := c.BindJSON(&updatePost); err != nil {
			log.Fatal("Error, update bind json")
		}
		// check update post input body not empty
		if strings.TrimSpace(updatePost.Body) == "" {
			c.JSON(http.StatusNotFound, gin.H{"msg": "Post must not be empty"})
			return
		}
		// get route param id
		id := c.Param("id") 
		// find the post need to updata
		res := db.Where("id = ?", id).Find(&post)
		// check post exist?
		if res.RowsAffected < 1 {
			c.JSON(http.StatusNotFound, gin.H{"msg": "No Post"})
			return
		}
		// check user get auth to update
		username := c.GetString("username")
		if post.UserName != username {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "have no auth to update"})
			return
		}
		// update post 
		db.Model(&post).Update("body", updatePost.Body)
		// res
		c.JSON(http.StatusOK, post)
	}
}