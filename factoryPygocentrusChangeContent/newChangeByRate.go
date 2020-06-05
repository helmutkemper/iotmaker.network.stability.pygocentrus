package factoryPygocentrusChangeContent

import "github.com/helmutkemper/pygocentrus/changeContent"

func NewChangeByRate(attackRate, lengthMin, lengthMax float64) (error, changeContent.ChangeContent) {
	var err error
	ret := changeContent.ChangeContent{
		ChangeRateMin: lengthMin,
		ChangeRateMax: lengthMax,
		Rate:          attackRate,
	}

	err = ret.Verify()

	return err, ret
}
