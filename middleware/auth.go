package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/songcser/gingo/config"
	"net/http"
)

type CustomClaims struct {
	UserID     int    `json:"user_id"`
	CustomData string `json:"custom_data"`
	jwt.StandardClaims
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			// 验证密钥

			myKey := config.GVA_CONFIG.User.AuthKey

			return []byte(myKey), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
			userId := claims.UserID
			customData := claims.CustomData

			println("获取当前 uuid : ", userId)

			// 在这里您可以访问解密后的声明内容，进行相应的处理

			c.Set("user_id", userId)
			c.Set("custom_data", customData)

			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
	}
}
