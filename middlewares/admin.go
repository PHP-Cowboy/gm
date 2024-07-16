package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// 超级管理员校验
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, ok := c.Get("claims")
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "claims获取失败",
			})
			c.Abort()
			return
		}

		userInfo := claims.(*CustomClaims)

		if userInfo == nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"msg": "无权限",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
