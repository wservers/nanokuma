# SPDX-License-Identifier: GPL-2.0-or-later
# Copyright (C) 2022-2026 WSERVER

SERVER_DIR = init/server/
AGENT_DIR = init/agent/
TARGET = nanokuma nanokuma-agent

all: server agent

server: $(SERVER_DIR)
	make -C $(SERVER_DIR)

agent: $(AGENT_DIR)
	make -C $(AGENT_DIR)

clean:
	@rm -f $(TARGET)

.PHONY: all server agent clean
