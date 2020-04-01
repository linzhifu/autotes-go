package user

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// Login user login autotest
func Login(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user User
		db.AutoMigrate(&user)
		data := make(map[string]interface{})
		resp := gin.H{
			"errcode": 0,
			"errmsg":  "ok",
		}
		type jsonData struct {
			LoginType string `json:"type"`
			Email     string `json:"email"`
			Captcha   string `json:"captcha"`
			Name      string `json:"username"`
			Psw       string `json:"password"`
		}
		var json jsonData
		c.ShouldBind(&json)
		switch loginType := json.LoginType; loginType {
		// 用户名密码登录
		case "username":
			name := json.Name
			psw := json.Psw
			db.Where(&User{Name: name}).First(&user)
			if user.Psw != psw {
				resp["errcode"] = 2
				resp["errmsg"] = "用户名密码不一致"
				c.JSON(http.StatusOK, resp)
				return
			}
		// 邮箱验证码登录
		case "email":
			email := json.Email
			captcha := json.Captcha
			// 对比缓存中的验证码是否一致
			cacheCaptcha, _ := cache.Get(email)
			if captcha != cacheCaptcha {
				resp["errcode"] = 2
				resp["errmsg"] = "邮箱验证码不一致"
				c.JSON(http.StatusOK, resp)
				return
			}
			var count int
			db.Where(&User{Email: email}).Find(&user).Count(&count)
			if count == 0 {
				user.Name = email
				user.Email = email
				user.Psw = "123"
				db.Create(&user)
			}
		default:
			resp["errcode"] = 1
			resp["errmsg"] = "invalid parameter"
			c.JSON(http.StatusOK, resp)
			return
		}
		// token
		token, err := GetToken(&user)
		if err != nil {
			resp["errcode"] = 3
			resp["errmsg"] = fmt.Sprintf("%s", err)
			c.JSON(http.StatusOK, resp)
			return
		}
		resp["token"] = token
		// id username email
		data["id"] = user.ID
		data["username"] = user.Name
		data["email"] = user.Email
		resp["data"] = data
		c.JSON(http.StatusOK, resp)
	}
}
