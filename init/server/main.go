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

package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"git.wh64.net/wserver/config"
	"git.wh64.net/wserver/nanokuma/core"
	cnf "git.wh64.net/wserver/nanokuma/server/config"
	"git.wh64.net/wserver/nanokuma/server/modules/database"
	"git.wh64.net/wserver/nanokuma/server/modules/webserver"
)

func main() {
	var err error
	var kuma *core.NanoKuma
	var quit chan os.Signal

	err = config.Load(cnf.CONFIG_PATH, &cnf.Get, cnf.DefaultConfig)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}

	kuma = core.NewNanoKuma("NanoKuma")

	kuma.AddModule(database.Database)
	kuma.AddModule(webserver.WebServer)

	err = kuma.Init()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}

	quit = make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	err = kuma.Destroy()
	if err != nil {
		return
	}

	err = config.Unload(&cnf.Get)
	if err != nil {
		return
	}
}
