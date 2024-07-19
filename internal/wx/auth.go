package wx

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/songcser/gingo/config"
	"github.com/songcser/gingo/middleware"
)

func GenerateToken(userId int, customData string) (string, error) {
	claims := middleware.CustomClaims{
		UserID:     userId,
		CustomData: customData,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 设置密钥并签名
	myKey := config.GVA_CONFIG.User.AuthKey
	signedToken, err := token.SignedString([]byte(myKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
