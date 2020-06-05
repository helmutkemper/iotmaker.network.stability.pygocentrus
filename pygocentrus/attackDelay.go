package pygocentrus

import "time"

func (el *Listener) attackDelay(inData []byte) (int, []byte) {
	//seelog.Debugf("%v%v were delayed by a pygocentrus attack: delay content", req.RemoteAddr, req.RequestURI)

	time.Sleep(time.Duration(el.randNumberBetweenRange(el.Pygocentrus.Delay.Min, el.Pygocentrus.Delay.Max)) * time.Microsecond)

	return len(inData), inData
}
