package connection

import "github.com/helmutkemper/pygocentrus/pygocentrus"

type Connection struct {
	Address  string
	Protocol pygocentrus.Protocol
}
