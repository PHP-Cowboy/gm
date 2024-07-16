package api

import (
	"github.com/gin-gonic/gin"
	"gm/daos"
	"gm/utils/rsp"
)

func Online(c *gin.Context) {

	online, err := daos.Online()
	if err != nil {
		rsp.ErrorJSON(c, err)
		return
	}

	rsp.SucJson(c, online)
}

func GetVersion(c *gin.Context) {
	rsp.SucJson(c, gin.H{"v": "1.0.0"})
}
