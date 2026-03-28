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
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"git.wh64.net/wserver/nanokuma/server/config"
	"github.com/gin-gonic/gin"
)

type WebServerModule struct {
	App       *gin.Engine
	webserver *http.Server
	errc      chan error
}

func (*WebServerModule) GetName() string {
	return "webserver"
}

func (m *WebServerModule) Load() error {
	var err error
	var conf = config.Get

	m.App = gin.Default()
	m.RouteAPI(m.App)

	m.webserver = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		Handler: m.App,
	}

	m.errc = make(chan error, 1)

	var proto = "http"
	if conf.SSL.Enable {
		proto = "https"
	}

	fmt.Printf("[webserver]: http webserver served at: %s://%s:%d\n", proto, config.Get.Host, config.Get.Port)

	go func() {
		if conf.SSL.Enable {
			m.errc <- m.webserver.ListenAndServeTLS(conf.SSL.CertFile, conf.SSL.KeyFile)
			return
		}

		m.errc <- m.webserver.ListenAndServe()
	}()

	select {
	case err = <-m.errc:
		return err
	case <-time.After(1500 * time.Millisecond):
		break
	}

	return nil
}

func (m *WebServerModule) Unload() error {
	var err error
	var ctx context.Context
	var cancel context.CancelFunc

	fmt.Printf("[webserver]: shutting down web server...\n")
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)

	err = m.webserver.Shutdown(ctx)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "[webserver]: FATAL: forced shutting down webserver...\n%v\n", err)
	}

	cancel()
	return nil
}

var WebServer *WebServerModule = &WebServerModule{}
