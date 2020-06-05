package pygocentrus

func (el *Listener) attackChangeContent(inData []byte) (int, []byte) {
	//seelog.Debugf("%v%v were eaten by a pygocentrus attack: change content", req.RemoteAddr, req.RequestURI)

	length := len(inData)
	forLength := el.Pygocentrus.ChangeContent.GetRandomByMaxMin(length)
	for i := 0; i != forLength; i += 1 {
		indexChange := el.Pygocentrus.ChangeContent.GetRandomByLength(length)
		inData = append(append(inData[:indexChange], byte(el.newRandGeneratorHeader().Intn(255))), inData[indexChange+1:]...)
	}

	return length, inData
}
