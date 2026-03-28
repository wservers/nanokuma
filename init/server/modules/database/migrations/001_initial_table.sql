-- SPDX-License-Identifier: GPL-2.0-or-later
-- Copyright (C) 2022-2026 WSERVER

-- create agents table
CREATE TABLE IF NOT EXISTS <prefix>agents(
	id VARCHAR(36) NOT NULL,
	ip_addr VARCHAR(45) NOT NULL,
	port INTEGER NOT NULL,
	hostname VARCHAR(50) NOT NULL,
	authorized TINYINT(1) NOT NULL DEFAULT 0,
	`status` ENUM('online', 'offline') DEFAULT 'offline',
	last_action_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	PRIMARY KEY(id)
);

-- create projects table
CREATE TABLE IF NOT EXISTS <prefix>projects(
	id VARCHAR(36) NOT NULL,
	repo_url VARCHAR(255) NOT NULL,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	PRIMARY KEY(id)
);

-- create jobs table
CREATE TABLE IF NOT EXISTS <prefix>jobs(
	id VARCHAR(36) NOT NULL,
	agent_id VARCHAR(36),
	project_id VARCHAR(36) NOT NULL,
	branch VARCHAR(255) NOT NULL,
	command VARCHAR(2048) NOT NULL,
	`state` ENUM('queued', 'running', 'success', 'failed') DEFAULT 'queued',
	created_at	DATETIME DEFAULT CURRENT_TIMESTAMP,
	started_at	DATETIME,
	finished_at DATETIME,
	lease_until DATETIME,
	PRIMARY KEY(id),
	FOREIGN KEY(project_id) REFERENCES <prefix>projects(id)
		ON UPDATE CASCADE ON DELETE CASCADE,
	FOREIGN KEY(agent_id) REFERENCES <prefix>agents(id)
		ON UPDATE CASCADE ON DELETE SET NULL
);
