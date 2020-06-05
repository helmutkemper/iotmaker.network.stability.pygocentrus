package changeContent

import (
	"math/rand"
	"time"
)

func (el *ChangeContent) GetRandomByMaxMin(length int) int {
	r1 := rand.New(rand.NewSource(time.Now().UnixNano()))

	if el.ChangeRateMin != 0.0 || el.ChangeRateMax != 0.0 {
		var changeMin = int(float64(length) * el.ChangeRateMin)
		var changeMax = int(float64(length) * el.ChangeRateMax)

		return r1.Intn(changeMax-changeMin) + changeMin
	}

	return r1.Intn(el.ChangeBytesMax-el.ChangeBytesMin) + el.ChangeBytesMin
}
