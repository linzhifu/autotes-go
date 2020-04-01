package user

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/dgrijalva/jwt-go"
)

// GetToken create a token
func GetToken(user *User) (string, error) {
	claims := &JWTClaims{
		ID: user.ID,
	}
	// 设置token生效时间
	claims.IssuedAt = time.Now().Unix()
	claims.ExpiresAt = time.Now().Add(time.Second * time.Duration(expireTime)).Unix()
	// 签发人
	claims.Issuer = "LeoLin"
	// // 使用指定的签名方法HS256创建签名对象token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(sercret)
}

// JWTAuth JWT auth
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("token")
		userID := c.Query("userID")
		if token == "" {
			c.JSON(http.StatusOK, gin.H{
				"errcode": 1,
				"errmsg":  "invalid parameter",
			})
			c.Abort()
			return
		}
		id, _ := strconv.Atoi(userID)
		var claims *JWTClaims
		claims, err := ParseToken(token)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"errcode": 3,
				"errmsg":  fmt.Sprintf("%s", err),
			})
			c.Abort()
			return
		}
		// 判断权限 get-无需权限，其他需要ID=1或者自己操作
		if c.Request.Method != http.MethodGet {
			if (claims.ID != 1) && (claims.ID != uint(id)) {
				c.JSON(http.StatusOK, gin.H{
					"errcode": 3,
					"errmsg":  "没有权限",
				})
				c.Abort()
				return
			}
		}
		c.Next()
		return
	}
}

// ParseToken parse JWT
func ParseToken(strToken string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(strToken, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return sercret, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, errMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, errExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, errNotValidYet
			} else {
				return nil, errInvalid
			}
		}
	}
	if err := token.Claims.Valid(); err != nil {
		return nil, err
	}
	return claims, nil
}
