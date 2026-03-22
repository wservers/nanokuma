package middleware

import (
	"git.wh64.net/wserver/nanokuma/server/modules/repo"
	"github.com/gin-gonic/gin"
)

func CheckRepoIsNil() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if repo.Repo == nil {
			ctx.JSON(500, gin.H{
				"ok":      0,
				"message": "\"repo\" service not served! please contact server administrator.",
			})
			return
		}

		ctx.Next()
	}
}
