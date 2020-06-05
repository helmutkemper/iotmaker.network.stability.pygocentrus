package factoryPygocentrus

import (
	"github.com/helmutkemper/pygocentrus/attack"
	"github.com/helmutkemper/pygocentrus/connection"
	"github.com/helmutkemper/pygocentrus/pygocentrus"
)

//en: Factory from pygocentrus to make a attack between two TCP connections.
//    Example: connection between code and database
//
//pt_br: Fabrica para montar um ataque do pygocentrus entre duas conexões TCP.
//    Exemplo: conexão entre o código e o banco de dados
//
//    Example:
//
//      import (
//          "github.com/helmutkemper/pygocentrus/attack"
//          "github.com/helmutkemper/pygocentrus/connection"
//          "github.com/helmutkemper/pygocentrus/factoryPygocentrusAttackRate"
//          "github.com/helmutkemper/pygocentrus/factoryPygocentrusConnection"
//          "github.com/helmutkemper/pygocentrus/listener"
//      )
//
//      func main() {
//        var err error
//        var attck attack.Attack
//        in := factoryPygocentrusConnection.NewConnectionTCP("10.0.0.1")
//        out := factoryPygocentrusConnection.NewConnectionTCP("10.0.0.2")
//        err, attck = factoryPygocentrusAttackRate.NewDelayAttack(
//          0.1,
//          200,
//          1000000,
//        )
//        if err != nil {
//          panic(err)
//        }
//
//        attackListener := NewPygocentrus(in, out, attck)
//        err = attackListener.AttackListener()
//        if err != nil {
//          panic(err)
//        }
//      }
func NewPygocentrus(inConnectionData, outConnectionData connection.Connection, attackData attack.Attack) pygocentrus.Listener {
	return pygocentrus.Listener{
		In:          inConnectionData,
		Out:         outConnectionData,
		Pygocentrus: attackData,
	}
}
