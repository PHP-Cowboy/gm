package api

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gm/daos"
	"gm/request"
	"gm/utils/ecode"
	"gm/utils/rsp"
)

func DictTypeList(c *gin.Context) {
	var req request.DictTypeList

	if err := c.ShouldBind(&req); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	result, err := daos.DictTypeList(req)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.SucJson(c, result)
	return
}

// 字典列表
func DictList(c *gin.Context) {
	var req request.DictList

	if err := c.ShouldBind(&req); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	list, err := daos.DictList(req)
	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.SucJson(c, list)
	return
}

func GetOneDict(c *gin.Context) {
	var req request.GetOneDict

	if err := c.ShouldBind(&req); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	dict, err := daos.GetOneDict(req)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.SucJson(c, dict)
	return
}

func SaveDict(c *gin.Context) {
	req := request.EditDict{}

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&req, bindingBody); err != nil {
		rsp.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	err := daos.SaveDict(req)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}
	rsp.Success(c)
}

func ChangeValues(c *gin.Context) {

	// 创建一个缓冲区来存储请求体的数据
	var buf bytes.Buffer

	// 将请求体的数据复制到缓冲区中
	_, err := buf.ReadFrom(c.Request.Body)
	if err != nil {
		rsp.ErrorJSON(c, ecode.ReadBodyError)
		return
	}

	// 重置请求体，以便后续的处理程序可以继续使用它（如果需要）
	//c.Request.Body = io.NopCloser(&buf)

	values := make(map[string]string, 0)

	err = json.Unmarshal(buf.Bytes(), &values)
	if err != nil {
		rsp.ErrorJSON(c, ecode.DataParsingFailure)
		return
	}

	err = daos.ChangeValues(values)

	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}
	rsp.Success(c)
}
