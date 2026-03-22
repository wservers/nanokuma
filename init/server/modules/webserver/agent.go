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
	"fmt"

	"git.wh64.net/wserver/nanokuma/include/agent"
	"git.wh64.net/wserver/nanokuma/server/modules/repo"
	"github.com/gin-gonic/gin"
)

func AgentCheck(ctx *gin.Context) {
	var err error
	var rp repo.RepoModule
	var payload agent.AgentPayload
	var data agent.AgentData

	rp = *repo.Repo

	err = ctx.ShouldBindJSON(&payload)
	if err != nil {
		ctx.JSON(400, gin.H{
			"ok":      0,
			"message": "payload is not json! please send payload for json.",
		})
		return
	}

	data = agent.AgentData{
		Id:       payload.Id,
		IPAddr:   payload.IPAddr,
		Port:     payload.Port,
		Hostname: payload.Hostname,
		Status:   payload.Status,
	}

	err = rp.UpsertAgent(&data)
	if err != nil {
		ctx.JSON(500, gin.H{
			"ok":      0,
			"message": "error occurred when checking agent status.",
		})
		return
	}

	ctx.JSON(200, gin.H{"ok": 1, "message": "agent upsert detected", "agent_id": payload.Id})
}

func AgentGet(ctx *gin.Context) {
	var err error
	var id string
	var rp repo.RepoModule
	var data *agent.AgentData

	id = ctx.Query("agent_id")
	if id == "" {
		ctx.JSON(400, gin.H{
			"ok":      0,
			"message": "\"id\" query must be contained",
		})
		return
	}

	rp = *repo.Repo

	data, err = rp.GetAgent(id)
	if err != nil {
		return
	}

	ctx.JSON(200, gin.H{
		"ok":      1,
		"message": "agent data loaded",
		"data":    *data,
	})
}

func AgentQuery(ctx *gin.Context) {
	var err error
	var rp repo.RepoModule
	var data []agent.AgentData = make([]agent.AgentData, 0)

	rp = *repo.Repo

	data, err = rp.GetAgents()
	if err != nil {
		return
	}

	ctx.JSON(200, gin.H{
		"ok":      1,
		"message": "all agent data loaded",
		"data":    data,
	})
}

func AgentAuthorize(ctx *gin.Context) {
	var err error
	var id string
	var rp repo.RepoModule

	id = ctx.Query("agent_id")
	if id == "" {
		ctx.JSON(400, gin.H{
			"ok":      0,
			"message": "\"id\" query must be contained",
		})
		return
	}

	rp = *repo.Repo

	_, err = rp.GetAgent(id)
	if err != nil {
		ctx.JSON(404, gin.H{
			"ok":      0,
			"message": "agent data is not found",
		})
		return
	}

	err = rp.AuthorizeAgent(id)
	if err != nil {
		ctx.JSON(500, gin.H{
			"ok":      0,
			"message": "failed authorize agent!",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"ok":      1,
		"message": "agent authorized!",
		"id":      id,
	})
}

func AgentDelete(ctx *gin.Context) {
	var err error
	var id string
	var rp repo.RepoModule

	id = ctx.Query("agent_id")
	if id == "" {
		ctx.JSON(400, gin.H{
			"ok":      0,
			"message": "\"id\" query must be contained",
		})
		return
	}

	rp = *repo.Repo

	err = rp.DeleteAgent(id)
	if err != nil {
		ctx.JSON(500, gin.H{
			"ok":      0,
			"message": fmt.Sprintf("failed to delete agent id for %s", id),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"ok":      1,
		"message": "agent deleted",
		"id":      id,
	})
}
