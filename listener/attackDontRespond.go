package listener

import "time"

func (el *Listener) attackDontRespond(inData []byte) (int, []byte) {
	//seelog.Debugf("%v%v were eaten by a pygocentrus attack: dont respond", req.RemoteAddr, req.RequestURI)

	time.Sleep(time.Duration(el.randNumberBetweenRange(el.Pygocentrus.Delay.Min, el.Pygocentrus.Delay.Max)) * time.Microsecond)
	return 0, nil
}
