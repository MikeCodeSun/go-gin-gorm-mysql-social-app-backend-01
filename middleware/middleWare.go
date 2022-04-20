package middleWare

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"github.com/mikecodesun/backend-sql/utils"
)

func Auth() gin.HandlerFunc{
	return func(c *gin.Context) {
		var claims utils.Claims
		// get env secret
		if err := godotenv.Load(); err != nil {
			log.Fatal("Error, env load")
		}
		secret := os.Getenv("SECRET")
		// get cookie token
		token, err := c.Cookie("token")
		if err != nil {
			log.Fatal("Error, get cookie token")
		}
		// if token not exist
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "Cookie token not exist "})
			c.Abort()
			return
		}
		// parse token
		t, errJwt := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		if errJwt != nil {
			log.Fatal("Error, jwt parse")
		}
		// check jwt auth valid
		if t.Valid == false {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "Cookie token wrong "})
			c.Abort()
			return
		}
		// log.Println(claims)
		// set user to req
		c.Set("username", claims.Name)
		c.Set("usertype", claims.UserType)
		// next
		c.Next()
	}
}