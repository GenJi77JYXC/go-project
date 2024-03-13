package middleware

import (
	"day1/model"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte("MetcB303")

type Claims struct {
	// 在payload里面放入用户id
	UserID uint
	jwt.RegisteredClaims
}

// 发token
func ReleaseToken(user model.User) (string, error) {
	// 过期时间一周
	expirationTime := time.Now().Add(24 * time.Hour * 7)
	claims := Claims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:  "Metc",
			Subject: "user token",
			//  之前版本是没有的，jwtv4之前这里有漏洞
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// 解token
func ParseToken(tokenString string) (*jwt.Token, *Claims, error) {
	claims := &Claims{}
	var token *jwt.Token
	var err error
	token, err = jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	return token, claims, err
}
