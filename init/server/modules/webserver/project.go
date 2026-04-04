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

	id, err = rp.CreateProject(payload)
	if err != nil {
		return
	}

	ctx.JSON(200, gin.H{
		"ok":      1,
		"message": "project created!",
		"id":      id,
	})
}

func ProjectRead(ctx *gin.Context) {
	var err error
	var id, repoURL string
	var rp repo.RepoModule
	var project *project.Project

	id = ctx.Query("project_id")
	repoURL = ctx.Query("repo_url")

	rp = *repo.Repo

	if id != "" {
		project, err = rp.GetProject(id)
	} else if repoURL != "" {
		project, err = rp.GetProjectByRepoURL(repoURL)
	} else {
		ctx.JSON(400, gin.H{
			"ok":      0,
			"message": "\"project_id\" or \"repo_url\" query must be contained",
		})
		return
	}

	if err != nil {
		ctx.JSON(500, gin.H{
			"ok":      0,
			"message": "failed to get the project information.",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"ok":      1,
		"message": "the project found",
		"data":    project,
	})
}

func ProjectUpdateRepoURL(ctx *gin.Context) {
	var err error
	var id string
	var rp repo.RepoModule
	var payload project.ProjectPayload

	id = ctx.Query("project_id")
	if id == "" {
		ctx.JSON(400, gin.H{
			"ok":      0,
			"message": "\"project_id\" query must be contained",
		})
		return
	}

	rp = *repo.Repo

	err = rp.UpdateProjectRepoURL(id, payload.RepoURL)
	if err != nil {
		ctx.JSON(500, gin.H{
			"ok":      0,
			"message": "failed to update repo url",
			"id":      id,
		})
	}
}

func ProjectDelete(ctx *gin.Context) {
	var err error
	var id string
	var rp repo.RepoModule

	id = ctx.Query("project_id")
	if id == "" {
		ctx.JSON(400, gin.H{
			"ok":      0,
			"message": "\"project_id\" query must be contained",
		})
		return
	}

	rp = *repo.Repo

	err = rp.DeleteProject(id)
	if err != nil {
		ctx.JSON(500, gin.H{
			"ok":      0,
			"message": "failed to delete project",
			"id":      id,
		})
		return
	}

	ctx.JSON(200, gin.H{
		"ok":      1,
		"message": "project deleted",
		"id":      id,
	})
}
