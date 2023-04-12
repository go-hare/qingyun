package token

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

var JwtKey = []byte("yingxiaozhu") // 加密key

type MyClaims struct {
	UserNmae string `json:"username"`
	UserId   int64  `json:"user_id"`
	jwt.StandardClaims
}

func CreateJwtToken(user_id int64) (jwtStr string, err error) {
	// 加密一个token
	claims := MyClaims{
		UserNmae: "fishing",
		UserId:   user_id,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix() - 60,      // 一分钟之前开始生效
			ExpiresAt: time.Now().Unix() + 60*60*2, // 两个小时后失效
			Issuer:    "签发人",                       // 签发人
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtKey)
}

func ParseJwtToken(token_string string) (myClaims *MyClaims, err error) {
	// 解密token
	parseToken, err := jwt.ParseWithClaims(token_string, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})
	if err != nil {
		return nil, err
	}
	return parseToken.Claims.(*MyClaims), nil
}
