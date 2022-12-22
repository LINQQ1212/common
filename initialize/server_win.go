//go:build windows

package initialize

import (
	"net/http"
	"time"
)

func initServer(address string, router http.Handler) server {
	return &http.Server{
		Addr:           address,
		Handler:        http.AllowQuerySemicolons(router),
		ReadTimeout:    20 * time.Second,
		WriteTimeout:   20 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
}
