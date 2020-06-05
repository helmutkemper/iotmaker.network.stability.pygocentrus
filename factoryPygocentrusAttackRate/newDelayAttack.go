package factoryPygocentrusAttackRate

import (
	"github.com/helmutkemper/pygocentrus"
)

func NewDelayAttack(attackRate float64, delayInMicrosecondMin, delayInMicrosecondMax int) (error, pygocentrus.Pygocentrus) {
	return nil, pygocentrus.Pygocentrus{
		Enabled: true,
		Delay: pygocentrus.RateMaxMin{
			Rate: attackRate,
			Min:  delayInMicrosecondMin,
			Max:  delayInMicrosecondMax,
		},
	}
}
