package pygocentrus

import "net"

type Connection struct {
	Address  string
	Protocol string
	A        *net.TCPAddr
}
