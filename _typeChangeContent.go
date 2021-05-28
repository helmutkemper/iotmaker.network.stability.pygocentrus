package pygocentrus

import (
	"errors"
	"math/rand"
	"time"
)

type ChangeContent struct {
	ChangeRateMin  float64
	ChangeRateMax  float64
	ChangeBytesMin int
	ChangeBytesMax int
	Rate           float64
}

func (el *ChangeContent) prepare() error {
	if el.Rate == 0.0 {
		return nil
	}

	if el.ChangeRateMin == el.ChangeRateMax && el.ChangeBytesMin == el.ChangeBytesMax && el.ChangeRateMin == 0.0 {
		el.Rate = 0.0
		return errors.New("pygocentrus attack > changeContent > rate max = 0 and rate min = 0")
	}

	if el.ChangeRateMin > el.ChangeRateMax {
		return errors.New("pygocentrus attack > changeContent > rate > the minimum value is greater than the maximum value")
	}

	if el.ChangeBytesMin > el.ChangeBytesMax {
		return errors.New("pygocentrus attack > changeContent > bytes > the minimum value is greater than the maximum value")
	}

	if (el.ChangeRateMin != 0.0 || el.ChangeRateMax != 0.0) && (el.ChangeBytesMin != 0.0 || el.ChangeBytesMax != 0.0) {
		return errors.New("pygocentrus attack > changeContent > you must choose option rate change or option bytes change")
	}

	return nil
}

func (el *ChangeContent) GetRandomByMaxMin(length int) int {
	r1 := rand.New(rand.NewSource(time.Now().UnixNano()))

	if el.ChangeRateMin != 0.0 || el.ChangeRateMax != 0.0 {
		var changeMin = int(float64(length) * el.ChangeRateMin)
		var changeMax = int(float64(length) * el.ChangeRateMax)

		return r1.Intn(changeMax-changeMin) + changeMin
	}

	return r1.Intn(el.ChangeBytesMax-el.ChangeBytesMin) + el.ChangeBytesMin
}

func (el *ChangeContent) GetRandomByLength(length int) int {
	return rand.New(rand.NewSource(time.Now().UnixNano())).Intn(length)
}
