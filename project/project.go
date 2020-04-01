package project

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// Project a project info
type Project struct {
	gorm.Model
	UserID      uint
	Name        string
	Description string
	Webresult   bool `gorm:"DEFAULT:false"`
	Apiresult   bool `gorm:"DEFAULT:false"`
	Appresult   bool `gorm:"DEFAULT:false"`
	Result      bool `gorm:"DEFAULT:false"`
}

// Views get all projects or add project
func Views(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		resp := gin.H{
			"errcode": 0,
			"errmsg":  "ok",
		}
		switch c.Request.Method {
		case http.MethodGet:
			var projects []Project
			db.Find(&projects)
			datas := make([]map[string]interface{}, len(projects))
			for i, project := range projects {
				// 返回信息
				data := make(map[string]interface{})
				data["id"] = project.ID
				data["proname"] = project.Name
				data["prodes"] = project.Description
				data["webresult"] = project.Webresult
				data["apiresult"] = project.Apiresult
				data["appresult"] = project.Appresult
				data["result"] = project.Result
				data["user"] = project.UserID
				data["update_time"] = project.UpdatedAt
				datas[i] = data
			}
			resp["datas"] = datas
			c.JSON(http.StatusOK, resp)
		case http.MethodPost:
			type jsonData struct {
				Name        string `json:"proname"`
				Description string `json:"prodes"`
				UserID      string `json:"user"`
			}
			var json jsonData
			err := c.ShouldBind(&json)
			if err != nil {
				resp["errcode"] = 1
				resp["errmsg"] = "invalid parameter"
				c.JSON(http.StatusOK, resp)
				return
			}
			userID, _ := strconv.Atoi(json.UserID)
			// 创建
			var project Project
			project.Name = json.Name
			project.Description = json.Description
			project.UserID = uint(userID)
			db.Create(&project)
			// 返回信息
			data := make(map[string]interface{})
			data["id"] = project.ID
			data["proname"] = project.Name
			data["prodes"] = project.Description
			data["webresult"] = project.Webresult
			data["apiresult"] = project.Apiresult
			data["appresult"] = project.Appresult
			data["result"] = project.Result
			data["user"] = project.UserID
			data["update_time"] = project.UpdatedAt
			resp["data"] = data
			c.JSON(http.StatusOK, resp)
		}

	}
}

// View edit/delete a project
func View(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var project Project
		resp := gin.H{
			"errcode": 0,
			"errmsg":  "ok",
		}
		id := c.Param("projectID")
		db.Where("ID=?", id).First(&project)
		switch c.Request.Method {
		// Patch
		case http.MethodPatch:
			type Info struct {
				Name        string `json:"proname"`
				Description string `json:"prodes"`
			}
			var info Info
			err := c.ShouldBind(&info)
			if err != nil {
				resp["errcode"] = 1
				resp["errmsg"] = "invalid parameter"
				c.JSON(http.StatusOK, resp)
				return
			}
			db.Model(&project).Updates(&info)
		// Delete
		case http.MethodDelete:
			if project.ID != 0 {
				db.Delete(&project)
			}
		}
		c.JSON(http.StatusOK, resp)
	}
}
