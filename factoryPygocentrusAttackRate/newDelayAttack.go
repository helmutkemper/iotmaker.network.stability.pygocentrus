package factoryPygocentrusAttackRate

import (
	"github.com/helmutkemper/pygocentrus/attack"
)

func NewDelayAttack(attackRate float64, delayInMicrosecondMin, delayInMicrosecondMax int) (error, attack.Attack) {
	return nil, attack.Attack{
		Enabled: true,
		Delay: attack.RateMaxMin{
			Rate: attackRate,
			Min:  delayInMicrosecondMin,
			Max:  delayInMicrosecondMax,
		},
	}
}
