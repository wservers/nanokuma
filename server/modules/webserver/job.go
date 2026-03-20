// SPDX-License-Identifier: GPL-2.0-or-later
/*
 * nanokuma
 * Copyright (C) 2022-2026 WSERVER
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 2 of the License, or
 * (at your option) any later version.
 */

package webserver

import (
	"git.wh64.net/wserver/nanokuma/server/modules/repo"
	"git.wh64.net/wserver/nanokuma/shared/agent"
	"git.wh64.net/wserver/nanokuma/shared/job"
	"github.com/gin-gonic/gin"
)

func JobCreate(ctx *gin.Context) {
	var err error
	var id, agentID string
	var rp repo.RepoModule
	var agent *agent.AgentData
	var payload job.JobPayload

	agentID = ctx.Query("agent_id")
	if agentID == "" {
		ctx.JSON(400, gin.H{
			"ok":      0,
			"message": "\"agent_id\" query must be contained",
		})
		return
	}

	if repo.Repo == nil {
		ctx.JSON(500, gin.H{
			"ok":      0,
			"message": "\"repo\" service not served! please contact server administrator.",
		})
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

	if err = ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(400, gin.H{
			"ok":      0,
			"message": "payload is not json! please send payload for json.",
		})
		return
	}

	if id, err = rp.CreateJob(&payload); err != nil {
		return
	}

	ctx.JSON(201, gin.H{"ok": 1, "message": "job created!", "id": id})
}

func JobRead(ctx *gin.Context) {
	var err error
	var job *job.Job
	var id, agentID string
	var rp repo.RepoModule
	var agent *agent.AgentData

	id = ctx.Query("job_id")
	if id == "" {
		ctx.JSON(400, gin.H{
			"ok":      0,
			"message": "\"job_id\" query must be contained",
		})
		return
	}

	agentID = ctx.Query("agent_id")
	if agentID == "" {
		ctx.JSON(400, gin.H{
			"ok":      0,
			"message": "\"agent_id\" query must be contained",
		})
		return
	}

	if repo.Repo == nil {
		ctx.JSON(500, gin.H{
			"ok":      0,
			"message": "\"repo\" service not served! please contact server administrator.",
		})
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

	job, err = rp.GetJob(id)
	if err != nil {
		ctx.JSON(500, gin.H{
			"ok":      0,
			"message": "failed to get the job information",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"ok":      1,
		"message": "the job found",
		"data":    job,
	})
}

func JobQuery(ctx *gin.Context) {}

func JobUpdateStatus(ctx *gin.Context) {}

func JobDelete(ctx *gin.Context) {
	var err error
	var id, agentID string
	var rp repo.RepoModule
	var agent *agent.AgentData

	id = ctx.Query("job_id")
	if id == "" {
		ctx.JSON(400, gin.H{
			"ok":      0,
			"message": "\"job_id\" query must be contained",
		})
		return
	}

	agentID = ctx.Query("agent_id")
	if agentID == "" {
		ctx.JSON(400, gin.H{
			"ok":      0,
			"message": "\"agent_id\" query must be contained",
		})
		return
	}

	if repo.Repo == nil {
		ctx.JSON(500, gin.H{
			"ok":      0,
			"message": "\"repo\" service not served! please contact server administrator.",
		})
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

	err = rp.DeleteJob(id)
	if err != nil {
		ctx.JSON(500, gin.H{
			"ok":      0,
			"message": "failed to delete job " + id,
		})
	}

	ctx.JSON(200, gin.H{
		"ok":      1,
		"message": "job deleted",
		"id":      id,
	})
}
