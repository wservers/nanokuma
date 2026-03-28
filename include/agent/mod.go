// SPDX-License-Identifier: MIT
// Copyright (C) 2022-2026 WSERVER

package agent

import "time"

type AgentStatus string

var (
	Online  AgentStatus = "online"
	Offline AgentStatus = "offline"
)

type AgentData struct {
	Id           string      `json:"id"`
	IPAddr       string      `json:"ip_addr"`
	Port         int         `json:"port"`
	Hostname     string      `json:"hostname"`
	Authorized   bool        `json:"authorized"`
	Status       AgentStatus `json:"status"`
	LastActionAt time.Time   `json:"last_action_at"`
}

type AgentPayload struct {
	Id       string      `json:"id"`
	IPAddr   string      `json:"ip_addr"`
	Port     int         `json:"port"`
	Hostname string      `json:"hostname"`
	Status   AgentStatus `json:"status"`
}
