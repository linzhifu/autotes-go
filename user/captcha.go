package user

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/UncleBig/goCache"
	"gopkg.in/gomail.v2"

	"github.com/gin-gonic/gin"
)

// Captcha get a captcha by email
func Captcha() gin.HandlerFunc {
	return func(c *gin.Context) {
		resp := gin.H{
			"errcode": 0,
			"errmsg":  "ok",
		}
		email, ok := c.GetQuery("email")
		if !ok {
			resp["errcode"] = 1
			resp["errmsg"] = "invalid parameter"
			c.JSON(http.StatusOK, resp)
			return
		}
		if _, ok := cache.Get(email); ok {
			resp["errcode"] = 2
			resp["errmsg"] = "验证码已发送，请稍后再获取"
			c.JSON(http.StatusOK, resp)
			return
		}
		code := GetRandomString(4)
		cache.Set(email, code, goCache.DefaultExpiration)
		err := SendCaptcha(email, code)
		if err != nil {
			resp["errcode"] = 3
			resp["errmsg"] = fmt.Sprintf("验证码发送失败：%s", err)
			c.JSON(http.StatusOK, resp)
			return
		}
		c.JSON(http.StatusOK, resp)
		return
	}
}

// GetRandomString new a captcha
func GetRandomString(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	rand.Seed(time.Now().Unix())
	for i := 0; i < length; i++ {
		result = append(result, bytes[rand.Intn(len(bytes))])
	}
	return string(result)
}

// SendCaptcha send captcha to user by email
func SendCaptcha(email string, code string) error {
	url := "http://autotest.longsys.com/"
	title := fmt.Sprintf("Captcha:%s", code)
	htmlContent := fmt.Sprintf(`<p>Hello %s：</p> 
	<p>This is longsys autotest system, your captcha is: %s</p>
	<a href='%s'>Please Login<a>`, email, code, url)
	m := gomail.NewMessage()
	m.SetHeader("From", "autotest<18129832245@163.com>")
	m.SetHeader("To", email)
	m.SetHeader("Subject", title)
	m.SetBody("text/html", htmlContent)

	d := gomail.NewDialer("smtp.163.com", 25, "18129832245@163.com", "Lin5535960")

	// Send the email to Bob, Cora and Dan.
	err := d.DialAndSend(m)
	return err
}
