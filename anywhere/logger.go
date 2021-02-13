package anywhere

import (
	"net"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-hclog"
)

type logger struct {
}

func Open(opts *hclog.LoggerOptions) (hclog.Logger, error) {
	addr := os.Getenv("HCLOG_ANYWHERE_ADDR")
	if addr != "" {
		proto := "tcp"
		if addr[0] == '/' {
			proto = "unix"
		}

		c, err := net.Dial(proto, addr)
		if err == nil {
			return WrapConn(c, opts)
		}
	}

	path := filepath.Join(os.TempDir(), "hclog-anywhere.sock")
	if _, err := os.Stat(path); err == nil {
		c, err := net.Dial("unix", path)

		if err == nil {
			return WrapConn(c, opts)
		}
	}

	return hclog.L(), nil
}

func WrapConn(c net.Conn, opts *hclog.LoggerOptions) (hclog.Logger, error) {
	local := *opts

	local.JSONFormat = true
	local.Output = c

	return hclog.New(&local), nil
}
