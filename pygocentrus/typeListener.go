package pygocentrus

import (
	"github.com/helmutkemper/pygocentrus/attack"
	"github.com/helmutkemper/pygocentrus/connection"
	"net"
)

type Listener struct {
	In            connection.Connection
	Out           connection.Connection
	Pygocentrus   attack.Attack
	error         error
	inConnection  net.Conn
	outConnection net.Conn
	attack        pygocentrusListenFunc
}
