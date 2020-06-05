package pygocentrus

import "net"

func (el *Listener) handle() {
	var err error
	var inChannel chan []byte
	var outChannel chan []byte
	var bytesBufferInChannel []byte
	var bytesBufferOutChannel []byte

	el.outConnection, err = net.Dial(el.Out.Protocol.String(), el.Out.Address)
	if err != nil {
		el.error = err
		return
	}

	inChannel = el.makeChannelFromConnection(el.inConnection)
	outChannel = el.makeChannelFromConnection(el.outConnection)

	if el.Pygocentrus.Enabled == true {

		var randAttack int

		var list = make([]pygocentrusListenFunc, 0)

		if el.Pygocentrus.Delay.Rate != 0.0 {

			if el.Pygocentrus.Delay.Rate >= el.newRandGeneratorHeader().Float64() {
				list = append(list, el.attackDelay)
			}

		}

		if el.Pygocentrus.DontRespond.Rate != 0.0 {

			if el.Pygocentrus.DontRespond.Rate >= el.newRandGeneratorHeader().Float64() {
				list = append(list, el.attackDontRespond)
			}

		}

		if el.Pygocentrus.DeleteContent != 0.0 {

			if el.Pygocentrus.DeleteContent >= el.newRandGeneratorHeader().Float64() {
				list = append(list, el.attackDeleteContent)
			}

		}

		if el.Pygocentrus.ChangeContent.Rate != 0.0 {

			if el.Pygocentrus.ChangeContent.Rate >= el.newRandGeneratorHeader().Float64() {
				list = append(list, el.attackChangeContent)
			}

		}

		listLength := len(list)
		if listLength != 0 {
			//el.Pygocentrus.SetAttack()
			randAttack = el.newRandGeneratorHeader().Intn(len(list))
			el.attack = list[randAttack]
		}
	}

	for {
		select {
		case bytesBufferInChannel = <-inChannel:
			if bytesBufferInChannel == nil {
				return
			} else {
				_, err = el.outConnection.Write(bytesBufferInChannel)
				if err != nil {
					el.error = err
					return
				}
			}
		case bytesBufferOutChannel = <-outChannel:
			if bytesBufferOutChannel == nil {
				return
			} else {
				_, err = el.inConnection.Write(bytesBufferOutChannel)
				if err != nil {
					el.error = err
					return
				}
			}
		}
	}
}
