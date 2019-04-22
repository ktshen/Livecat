package crawler

import (
	"net"
)

func getFreePort() (port int, err error) {
	ln, err := net.Listen("tcp", "[::]:0")
	handleError(err, "[getPort] listen err")
	port = ln.Addr().(*net.TCPAddr).Port
	err = ln.Close()
	return
}
