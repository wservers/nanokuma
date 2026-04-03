package webserver

import (
	"git.wh64.net/wserver/nanokuma/include/project"
	"git.wh64.net/wserver/nanokuma/server/modules/repo"
	"github.com/gin-gonic/gin"
)

func ProjectCreate(ctx *gin.Context) {
	var err error
	var id string
	var rp repo.RepoModule
	var payload project.ProjectPayload

	err = ctx.ShouldBindJSON(&payload)
	if err != nil {
		ctx.JSON(400, gin.H{
			"ok":      0,
			"message": "payload is not a json! please send payload for json.",
		})
		return
	}

	rp = *repo.Repo

	id, err = rp.ProjectCreate(payload)
	if err != nil {
		return
	}

	ctx.JSON(200, gin.H{"ok": 1, "message": "project created!", "id": id})
}

func ProjectRead(ctx *gin.Context) {}

func ProjectUpdateRepoURL(ctx *gin.Context) {}

func ProjectDelete(ctx *gin.Context) {}
