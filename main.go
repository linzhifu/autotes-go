package main

import (
	"autotest/project"
	"autotest/user"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func dbInit() *gorm.DB {
	db, err := gorm.Open("mysql", "leo:123456@(127.0.0.1)/autotest_go?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect autotest_go database")
	}
	// user
	var user user.User
	db.AutoMigrate(&user)
	// product
	var project project.Project
	db.AutoMigrate(&project)
	return db
}

func main() {
	db := dbInit()
	defer db.Close()
	r := gin.Default()
	r.Use(cors.Default())
	r.GET("/api/v1/captcha", user.Captcha())
	r.POST("/api/v1/login", user.Login(db))
	r.GET("/api/v1/user", user.JWTAuth(), user.Views(db))
	r.Any("/api/v1/user/:userID", user.JWTAuth(), user.View(db))
	r.Any("/api/v1/project", user.JWTAuth(), project.Views(db))
	r.Any("/api/v1/project/:projectID", user.JWTAuth(), project.View(db))
	r.Run(":8000")
}

// Cors CORS
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		origin := c.Request.Header.Get("Origin")
		var headerKeys []string
		for k := range c.Request.Header {
			headerKeys = append(headerKeys, k)
		}
		headerStr := strings.Join(headerKeys, ", ")
		if headerStr != "" {
			headerStr = fmt.Sprintf("access-control-allow-origin, access-control-allow-headers, %s", headerStr)
		} else {
			headerStr = "access-control-allow-origin, access-control-allow-headers"
		}
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Headers", "*")
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Set("content-type", "application/json")
		}

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
		}

		c.Next()
	}
}
