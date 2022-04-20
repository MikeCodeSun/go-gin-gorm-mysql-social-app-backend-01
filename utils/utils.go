package utils

import (
	"errors"
	"log"
	

	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
)

// Claims struct
type Claims struct {
	Name string
	UserType string
	jwt.StandardClaims
}
// generate jwt token
func GenerateJwt(name string, userType string) string{
	// get env secret varible
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error env load")
	}
	secret := os.Getenv("SECRET")
// claims 
	var claims = Claims{
		Name: name,
		UserType: userType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * 24).Unix(),
		},
	}
// generate new token with claims
token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
if err != nil {
	log.Fatal("Error, Sign jwt token")
}
// return token
return token
}

// check user type
func CheckUser(c *gin.Context) error{
	userType := c.GetString("usertype")
	  // log.Printf("user type: %s ",  userType)
	  var err error
	  err = nil 
	if userType != "admin" {
		err = errors.New("Un auth user type")
		return err
	}
	
	return err
}