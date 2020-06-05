package pygocentrus

func (el *Listener) attackDeleteContent(inData []byte) (int, []byte) {
	//seelog.Debugf("%v%v were eaten by a pygocentrus attack: delete content", req.RemoteAddr, req.RequestURI)

	n := len(inData)
	inData = make([]byte, n)

	return n, inData
}
