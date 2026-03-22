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
	"git.wh64.net/wserver/nanokuma/include/job"
	"git.wh64.net/wserver/nanokuma/server/modules/repo"
	"github.com/gin-gonic/gin"
)

func JobCreate(ctx *gin.Context) {
	var err error
	var id string
	var rp repo.RepoModule
	var payload job.JobPayload

	rp = *repo.Repo

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
	var id string
	var job *job.Job
	var rp repo.RepoModule

	id = ctx.Query("job_id")
	if id == "" {
		ctx.JSON(400, gin.H{
			"ok":      0,
			"message": "\"job_id\" query must be contained",
		})
		return
	}

	rp = *repo.Repo

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

func JobQuery(ctx *gin.Context) {
	var err error
	var jobs []job.Job
	var projectID string
	var rp repo.RepoModule

	projectID = ctx.Query("agent_id")
	if projectID == "" {
		ctx.JSON(400, gin.H{
			"ok":      0,
			"message": "\"job_id\" query must be contained",
		})
		return
	}

	rp = *repo.Repo

	jobs, err = rp.GetJobs(projectID)
	if err != nil {
		ctx.JSON(500, gin.H{
			"ok":      0,
			"message": "failed to get the job information",
		})
		return
	}

	if len(jobs) == 0 {
		ctx.JSON(404, gin.H{
			"ok":      0,
			"message": "jobs not found by project id " + projectID,
		})
		return
	}

	ctx.JSON(200, gin.H{
		"ok":      1,
		"message": "the job found",
		"data":    jobs,
	})
}

func JobUpdateStatus(ctx *gin.Context) {
	var err error
	var id string
	var rp repo.RepoModule
	var payload struct {
		State job.JobState `json:"state" binding:"required"`
	}

	id = ctx.Query("job_id")
	if id == "" {
		ctx.JSON(400, gin.H{
			"ok":      0,
			"message": "\"job_id\" query must be contained",
		})
		return
	}

	rp = *repo.Repo

	if err = ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(400, gin.H{
			"ok":      0,
			"message": "payload is not json! please send payload for json.",
		})
		return
	}

	if err = rp.UpdateJobState(id, payload.State); err != nil {
		ctx.JSON(500, gin.H{
			"ok":      0,
			"message": "failed to update job state",
		})
	}

	ctx.JSON(200, gin.H{
		"ok":      1,
		"message": "success to update job state!",
		"id":      id,
	})
}

func JobDelete(ctx *gin.Context) {
	var err error
	var id string
	var rp repo.RepoModule

	id = ctx.Query("job_id")
	if id == "" {
		ctx.JSON(400, gin.H{
			"ok":      0,
			"message": "\"job_id\" query must be contained",
		})
		return
	}

	rp = *repo.Repo

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
