package listener

import (
	"github.com/helmutkemper/pygocentrus"
	"github.com/helmutkemper/pygocentrus/connection"
	"net"
)

type Listener struct {
	In            connection.Connection
	Out           connection.Connection
	Pygocentrus   pygocentrus.Pygocentrus
	error         error
	inConnection  net.Conn
	outConnection net.Conn
	attack        pygocentrusListenFunc
}
