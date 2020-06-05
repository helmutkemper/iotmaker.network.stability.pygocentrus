package factoryPygocentrusConnection

import (
	"github.com/helmutkemper/pygocentrus/connection"
	"github.com/helmutkemper/pygocentrus/listener"
)

/*
en: Prepare a new TCP connection
pt_br: Prepara uma nova conexão TCP
*/
func NewConnectionTCP(address string) connection.Connection {
	return connection.Connection{
		Address:  address,
		Protocol: listener.KProtocolTCP,
	}
}
