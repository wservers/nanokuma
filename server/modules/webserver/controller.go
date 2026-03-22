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
	"git.wh64.net/wserver/nanokuma/server/modules/webserver/middleware"
	"github.com/gin-gonic/gin"
)

func (m *WebServerModule) RouteAPI(app *gin.Engine) {
	var api *gin.RouterGroup
	api = app.Group("/api", middleware.CheckRepoIsNil())

	// Agents
	api.POST("/agent/check", AgentCheck)
	api.GET("/agent", AgentGet)
	api.GET("/agents", AgentQuery)
	api.PUT("/agent/authorize", AgentAuthorize)
	api.DELETE("/agent/delete", AgentDelete)

	// Jobs
	api.POST("/job", middleware.CheckAgentIsAuthorized(), JobCreate)
	api.GET("/job", middleware.CheckAgentIsAuthorized(), JobRead)
	api.GET("/jobs", middleware.CheckAgentIsAuthorized(), JobQuery)
	api.PATCH("/job", middleware.CheckAgentIsAuthorized(), JobUpdateStatus)
	api.DELETE("/job", middleware.CheckAgentIsAuthorized(), JobDelete)
}
