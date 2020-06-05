package connection

import "github.com/helmutkemper/pygocentrus/listener"

type Connection struct {
	Address  string
	Protocol listener.Protocol
}
