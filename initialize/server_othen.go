//go:build !windows

package initialize

import (
	"github.com/fvbock/endless"
	"net/http"
	"time"
)

func initServer(address string, router http.Handler) server {
	s := endless.NewServer(address, http.AllowQuerySemicolons(router))
	s.ReadHeaderTimeout = 20 * time.Second
	s.WriteTimeout = 20 * time.Second
	s.MaxHeaderBytes = 1 << 20
	s.ErrorLog = nil
	return s
}
