package testcase

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// TestCase a case info
type TestCase struct {
	gorm.Model
	UserID      uint
	ProjectID   uint
	Name        string
	Description string
	Type        string
	Result      bool `gorm:"DEFAULT:false"`
	Index       int  `gorm:"DEFAULT:1"`
}

// Views get all cases or add a case
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
			var testCases []TestCase
			db.Where("type = ? AND project_id = ?", srcType, projectID).Find(&testCases)
			datas := make([]map[string]interface{}, len(testCases))
			for i, testCase := range testCases {
				// 返回信息
				data := make(map[string]interface{})
				data["id"] = testCase.ID
				data["appname"] = testCase.Name
				data["appdes"] = testCase.Description
				data["index"] = testCase.Index
				data["user"] = testCase.UserID
				data["update_time"] = testCase.UpdatedAt
				data["src_type"] = testCase.Type
				data["result"] = testCase.Result
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
			var testCase TestCase
			testCase.Name = json.Name
			testCase.Description = json.Description
			testCase.Type = json.Type
			testCase.UserID = json.UserID
			testCase.ProjectID = json.ProjectID
			db.Create(&testCase)
			c.JSON(http.StatusOK, resp)
		}

	}
}

// View edit/delete a testCase
func View(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var testCase TestCase
		resp := gin.H{
			"errcode": 0,
			"errmsg":  "ok",
		}
		id := c.Param("testCaseID")
		db.Where("ID=?", id).First(&testCase)
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
			db.Model(&testCase).Updates(&info)
		// Delete
		case http.MethodDelete:
			if testCase.ID != 0 {
				db.Delete(&testCase)
			}
		}
		c.JSON(http.StatusOK, resp)
	}
}
