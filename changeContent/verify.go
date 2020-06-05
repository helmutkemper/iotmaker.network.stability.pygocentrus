package changeContent

import "errors"

func (el *ChangeContent) Verify() error {
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
