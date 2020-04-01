package script

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// Script a script info
type Script struct {
	gorm.Model
	UserID      uint
	ProjectID   uint
	Name        string
	Description string
	Type        string
	Result      bool `gorm:"DEFAULT:false"`
	Index       int  `gorm:"DEFAULT:1"`
}

// Views get all scripts or add script
func Views(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		resp := gin.H{
			"errcode": 0,
			"errmsg":  "ok",
		}
		switch c.Request.Method {
		case http.MethodGet:
			projectID := c.Query("project")
			srcType := c.Query("src_type")
			var scripts []Script
			db.Where("type = ? AND project_id = ?", srcType, projectID).Find(&scripts)
			datas := make([]map[string]interface{}, len(scripts))
			for i, script := range scripts {
				// 返回信息
				data := make(map[string]interface{})
				data["id"] = script.ID
				data["appname"] = script.Name
				data["appdes"] = script.Description
				data["index"] = script.Index
				data["user"] = script.UserID
				data["update_time"] = script.UpdatedAt
				data["result"] = script.Result
				datas[i] = data
			}
			resp["datas"] = datas
			c.JSON(http.StatusOK, resp)
		case http.MethodPost:
			type jsonData struct {
				Name        string `json:"appname"`
				Description string `json:"appdes"`
				Type        string `json:"src_type"`
				UserID      uint   `json:"user"`
				ProjectID   uint   `json:"project"`
			}
			var json jsonData
			err := c.ShouldBind(&json)
			if err != nil {
				resp["errcode"] = 1
				resp["errmsg"] = "invalid parameter"
				c.JSON(http.StatusOK, resp)
				return
			}
			// 创建
			var script Script
			script.Name = json.Name
			script.Description = json.Description
			script.Type = json.Type
			script.UserID = json.UserID
			script.ProjectID = json.ProjectID
			db.Create(&script)
			c.JSON(http.StatusOK, resp)
		}

	}
}

// View edit/delete a script
func View(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var script Script
		resp := gin.H{
			"errcode": 0,
			"errmsg":  "ok",
		}
		id := c.Param("scriptID")
		db.Where("ID=?", id).First(&script)
		switch c.Request.Method {
		// Patch
		case http.MethodPatch:
			type Info struct {
				Name        string `json:"appname"`
				Description string `json:"appdes"`
				Index       int    `json:"index"`
			}
			var info Info
			err := c.ShouldBind(&info)
			if err != nil {
				resp["errcode"] = 1
				resp["errmsg"] = "invalid parameter"
				c.JSON(http.StatusOK, resp)
				return
			}
			db.Model(&script).Updates(&info)
		// Delete
		case http.MethodDelete:
			if script.ID != 0 {
				db.Delete(&script)
			}
		}
		c.JSON(http.StatusOK, resp)
	}
}
