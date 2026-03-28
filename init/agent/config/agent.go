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

package config

type ServerConn struct {
	Host string
	Port int
}

type AgentConfig struct {
	Server ServerConn
}

const CONFIG_PATH = "agent-config.toml"

var (
	Get           AgentConfig = AgentConfig{}
	DefaultConfig AgentConfig = AgentConfig{
		Server: ServerConn{
			Host: "0.0.0.0",
			Port: 8080,
		},
	}
)
