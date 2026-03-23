package middleware

import (
	"git.wh64.net/wserver/nanokuma/include/agent"
	"git.wh64.net/wserver/nanokuma/server/modules/repo"
	"github.com/gin-gonic/gin"
)

func CheckAgentIsAuthorized() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var err error
		var id string
		var rp repo.RepoModule
		var agent *agent.AgentData

		id = ctx.Request.Header.Get("Agent-ID")
		if id == "" {
			return
		}

		rp = *repo.Repo

		agent, err = rp.GetAgent(id)
		if err != nil {
			ctx.JSON(500, gin.H{
				"ok":      0,
				"message": "failed to get the agent information",
			})
			return
		}

		if !agent.Authorized {
			ctx.JSON(403, gin.H{
				"ok":      0,
				"message": "the agent is not authorized",
			})
			return
		}

		ctx.Next()
	}
}
