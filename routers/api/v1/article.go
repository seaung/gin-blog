package v1

import (
	"net/http"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"github.com/astaxie/beego/validation"

	"github.com/seaung/gin-blog/models"
	"github.com/seaung/gin-blog/pkg/e"
	"github.com/seaung/gin-blog/pkg/setting"
	"github.com/seaung/gin-blog/pkg/util"
)

func GetArticle(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()
	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID Must gt 0")

	code := e.INVALID_PARAMS
	var data interface{}

	if !valid.HasErrors() {
		if models.ExistArticleByID(id) {
			data = models.GetArticle(id)
			code = e.SUCCESS
		} else {
			code = e.ERROR_NOT_EXIST_ARTICLE
		}
	} else {
		for _, err := range valid.Errors {
			log.Printf("err.key: %s err.message: %s", err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg": e.GetMsg(code),
		"data": data,
	})
}

func GetArticles(c *gin.Context) {
	data := make(map[string]interface{})
	maps := make(map[string]interface{})
	valid := validation.Validation{}

	var state int = -1

	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		maps["state"] = state
		valid.Range(state, 0, 1, "state").Message("state 0 or 1")
	}

	var tagId int = -1
	if arg := c.Query("tag_id"); arg != "" {
		tagId = com.StrTo(arg).MustInt()
		maps["tag_id"] = tagId
		valid.Min(tagId, 1, "tag_id").Message("tagid must gt 0")
	}

	code := e.INVALID_PARAMS

	if !valid.HasErrors() {
		code = e.SUCCESS
		data["lists"] = models.GetArticles(util.GetPage(c), setting.PageSize, maps)
		data["total"] = models.GetArticleTotal(maps)
	} else {
		for _, err := range valid.Errors {
			log.Printf("err.key: %s, err.message : %s", err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg": e.GetMsg(code),
		"data": data,
	})
}

func AddArticle(c *gin.Context) {
	tagId := com.StrTo(c.Query("tag_id")).MustInt()
	title := c.Query("title")
	desc := c.Query("desc")
	content := c.Query("content")
	createdBy := c.Query("created_by")
	state := com.StrTo(c.DefaultQuery("state", "0")).MustInt()

	valid := validation.Validation{}
	valid.Min(tagId, 1, "tag_id").Message("tag gt 0")
	valid.Required(title, "title").Message("title is Required")
	valid.Required(desc, "desc").Message("desc is Required")
	valid.Required(content, "content").Message("content is Required")
	valid.Required(createdBy, "created_by").Message("created_by is Required")
	valid.Range(state, 0, 1, "state").Message("state 0 or 1")

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		if models.ExistTagByID(tagId) {
			data := make(map[string]interface{})
			data["tag_id"] = tagId
			data["title"] = title
			data["desc"] = desc
			data["content"] = content
			data["created_by"] = createdBy
			data["state"] = state

			models.AddArticle(data)
			code = e.SUCCESS
		} else {
			code = e.ERROR_EXIST_TAG
		}
	} else {
		for _, err := range valid.Errors {
			log.Printf("err.key : %s, err.message : %s", err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg": e.GetMsg(code),
		"data": make(map[string]interface{}),
	})
}

func EditArticle(c *gin.Context) {
	valid := validation.Validation{}

	id := com.StrTo(c.Param("id")).MustInt()
	tagId := com.StrTo(c.Query("tag_di")).MustInt()
	title := c.Query("desc")
	content := c.Query("content")
	desc := c.Query("desc")
	modifiedBy := c.Query("modified_by")

	var state int = -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		valid.Range(state, 0, 1, "state").Message("state 0 or 1")
	}
	valid.Min(id, 1, "id").Message("id gt 0")
	valid.MaxSize(title, 100, "title").Message("title is Required")
	valid.MaxSize(desc, 255, "desc").Message("title is Required")
	valid.MaxSize(content, 65535, "content").Message("content is Required")
	valid.Required(modifiedBy, "modified_by").Message("modifiedBy is Required")
	valid.MaxSize(modifiedBy, 100, "modified_by").Message("modifiedBy gt 100")

	code := e.INVALID_PARAMS

	if !valid.HasErrors() {
		if models.ExistArticleByID(id) {
			if models.ExistTagByID(tagId) {
				data := make(map[string]interface{})
				if tagId > 0 {
					data["tag_id"] = tagId
				}

				if title != "" {
					data["title"] = title
				}

				if desc != "" {
					data["desc"] = desc
				}

				if content != "" {
					data["content"] = content
				}

				data["modified_by"] = modifiedBy
				models.EditArticle(id, data)
				code = e.SUCCESS
			} else {
				code = e.ERROR_NOT_EXIST_ARTICLE
			}
		} else {
			for _, err := range valid.Errors {
				log.Printf("err.key : %s, err.message : %s", err.Key, err.Message)
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg": e.GetMsg(code),
		"data": make(map[string]string),
	})
}

func DeleteArticle(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()
	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID must gt 0")
	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		if models.ExistArticleByID(id) {
			models.DeleteArticle(id)
			code = e.SUCCESS
		} else {
			code = e.ERROR_NOT_EXIST_ARTICLE
		}
	} else {
		for _, err := range valid.Errors {
			log.Printf("err.key : %s, err.message : %s", err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg": e.GetMsg(code),
		"data": make(map[string]string),
	})
}
