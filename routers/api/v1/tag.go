package v1

import (
	"net/http"

	"github.com/unknwon/com"
	"github.com/gin-gonic/gin"
	"github.com/astaxie/beego/validation"

	"github.com/seaung/gin-blog/pkg/e"
	"github.com/seaung/gin-blog/pkg/util"
	"github.com/seaung/gin-blog/pkg/setting"
	"github.com/seaung/gin-blog/models"
)

func GetTags(c *gin.Context) {
	name := c.Query("name")

	maps := make(map[string]interface{})
	data := make(map[string]interface{})

	if name != "" {
		maps["name"] = name
	}

	var state int = -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		maps["state"] = state
	}

	code := e.SUCCESS

	data["list"] = models.GetTags(util.GetPage(c), setting.PageSize, maps)
	data["total"] = models.GetTagTotal(maps)

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg": e.GetMsg(code),
		"data": data,
	})
}

func AddTag(c *gin.Context) {
	name := c.Query("name")
	state := com.StrTo(c.DefaultQuery("state", "0")).MustInt()
	createdBy := c.Query("created_by")

	valid := validation.Validation{}
	valid.Required(name, "name").Message("name is required")
	valid.MaxSize(name, 100, "name").Message("name max length 100")
	valid.Required(createdBy, "created_by").Message("created_by is Required")
	valid.MaxSize(createdBy, 100, "created_by").Message("created_by max length 100")
	valid.Range(state, 0, 1, "state").Message("state 0 or 1")

	code := e.INVALID_PARAMS

	if !valid.HasErrors() {
		code = e.SUCCESS
		models.AddTag(name, state, createdBy)
	} else {
		code = e.ERROR_EXIST_TAG
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg": e.GetMsg(code),
		"data": make(map[string]string),
	})
}

func EditTag(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()
	name := c.Query("name")
	modifiedBy := c.Query("modified_by")
	valid := validation.Validation{}
	var state int = -1

	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		valid.Range(state, 0, 1, "state").Message("state 0 or 1")
	}

	valid.Required(id, "id").Message("ID is required")
	valid.Required(modifiedBy, "modified_by").Message("modified_by is require")
	valid.MaxSize(modifiedBy, 100, "modified_by").Message("modified_by length 100")
	valid.MaxSize(name, 100, "name").Message("name max length 100")

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		code = e.SUCCESS
		if models.ExistTagByID(id) {
			data := make(map[string]interface{})
			data["modified_by"] = modifiedBy
			if name != "" {
				data["name"] = name
			}

			if state != 1 {
				data["state"] = state
			}
			models.EditTag(id, data)
		} else {
			code = e.ERROR_NOT_EXIST_TAG
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg": e.GetMsg(code),
		"data": make(map[string]string),
	})
}

func DeleteTag(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID must gt 0")

	code := e.INVALID_PARAMS

	if !valid.HasErrors() {
		code = e.SUCCESS
		if models.ExistTagByID(id) {
			models.DeleteTag(id)
		} else {
			code = e.ERROR_NOT_EXIST_TAG
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg": e.GetMsg(code),
		"data": make(map[string]string),
	})
}
