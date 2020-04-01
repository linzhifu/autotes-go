package user

import (
	"crypto/md5"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/UncleBig/goCache"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

var (
	// 基于内存的key-value缓存，用于记录邮箱验证码
	cache = goCache.New(5*time.Minute, 10*time.Minute)
	// Sercret
	sercret = []byte("xiaoGuaGua")
	// token有效时间，单位s
	expireTime = 3600 * 10
	// error
	errExpired     error = errors.New("Token is expired")
	errNotValidYet error = errors.New("Token not active yet")
	errMalformed   error = errors.New("That's not even a token")
	errInvalid     error = errors.New("Couldn't handle this token")
)

//User info
type User struct {
	gorm.Model
	Name  string
	Email string
	Psw   string
}

// JWTClaims token infomation
type JWTClaims struct {
	jwt.StandardClaims
	ID uint `json:"userID"`
}

// BeforeSave encrypy psw by md5 before save
func (user *User) BeforeSave() {
	md5Psw := md5.Sum([]byte(user.Psw))
	user.Psw = fmt.Sprintf("%x", md5Psw)
}

// View user
func View(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user User
		data := make(map[string]interface{})
		resp := gin.H{
			"errcode": 0,
			"errmsg":  "ok",
		}
		id := c.Param("userID")
		db.Where("ID=?", id).First(&user)
		// Patch
		if c.Request.Method == http.MethodPatch {
			type Info struct {
				Name string `json:"username"`
				Psw  string `json:"password"`
			}
			var info Info
			err := c.ShouldBind(&info)
			if err != nil {
				resp["errcode"] = 1
				resp["errmsg"] = "invalid parameter"
				c.JSON(http.StatusOK, resp)
				return
			}
			// 保存前做MD5加密
			if info.Psw != "" {
				md5Psw := md5.Sum([]byte(info.Psw))
				info.Psw = fmt.Sprintf("%x", md5Psw)
			}
			db.Model(&user).Updates(&info)
		}
		// Delete
		if c.Request.Method == http.MethodDelete {
			if user.ID != 0 {
				db.Delete(&user)
			}
		}
		data["id"] = user.ID
		data["username"] = user.Name
		data["email"] = user.Email
		resp["data"] = data
		c.JSON(http.StatusOK, resp)
	}
}

// Views get all users
func Views(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var users []User
		resp := gin.H{
			"errcode": 0,
			"errmsg":  "ok",
		}
		db.Find(&users)
		datas := make([]map[string]interface{}, len(users))
		for i, user := range users {
			data := make(map[string]interface{})
			data["id"] = user.ID
			data["username"] = user.Name
			data["email"] = user.Email
			datas[i] = data
		}
		resp["data"] = datas
		c.JSON(http.StatusOK, resp)
	}
}
