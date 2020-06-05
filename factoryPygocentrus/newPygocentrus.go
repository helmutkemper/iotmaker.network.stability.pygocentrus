package factoryPygocentrus

import (
	"github.com/helmutkemper/pygocentrus"
	"github.com/helmutkemper/pygocentrus/connection"
	"github.com/helmutkemper/pygocentrus/factoryPygocentrusAttackRate"
	"github.com/helmutkemper/pygocentrus/factoryPygocentrusConnection"
	"github.com/helmutkemper/pygocentrus/listener"
)

func NewPygocentrus(inConnectionData, outConnectionData connection.Connection, attackData pygocentrus.Pygocentrus) listener.Listener {
	return listener.Listener{
		In:          inConnectionData,
		Out:         outConnectionData,
		Pygocentrus: attackData,
	}
}

func teste() {
	var err error
	var attack pygocentrus.Pygocentrus
	in := factoryPygocentrusConnection.NewConnectionTCP("10.0.0.1")
	out := factoryPygocentrusConnection.NewConnectionTCP("10.0.0.2")
	err, attack = factoryPygocentrusAttackRate.NewDelayAttack()
}
