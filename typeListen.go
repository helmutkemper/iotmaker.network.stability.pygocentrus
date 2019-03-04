package pygocentrus

import (
	"fmt"
	"math/rand"
	"net"
	"time"
)

type Listen struct {
	In            Connection
	Out           Connection
	Pygocentrus   Pygocentrus
	error         error
	inConnection  net.Conn
	outConnection net.Conn
}

func (el *Listen) Listen() error {
	var listener net.Listener
	var err error

	listener, err = net.Listen(el.In.Protocol, el.In.Address)
	if err != nil {
		return err
	}

	for {
		fmt.Println("pygocentrus incoming data...")
		el.inConnection, err = listener.Accept()
		if err != nil {
			return err
		}

		go el.handle()
		if el.error != nil {
			return el.error
		}

	}
}

func (el *Listen) handle() {
	var err error
	var inChannel chan []byte
	var outChannel chan []byte
	var bytesBufferInChannel []byte
	var bytesBufferOutChannel []byte

	el.outConnection, err = net.Dial(el.Out.Protocol, el.Out.Address)
	if err != nil {
		el.error = err
		return
	}

	inChannel = el.makeChannelFromConnection(el.inConnection)
	outChannel = el.makeChannelFromConnection(el.outConnection)

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

func (el *Listen) makeChannelFromConnection(conn net.Conn) chan []byte {
	var connectionDataChannel chan []byte

	connectionDataChannel = make(chan []byte)

	go func() {

		var bytesBuffer = make([]byte, 1024)
		var bufferLength int
		var err error

		for {

			bufferLength, err = conn.Read(bytesBuffer)
			fmt.Printf("bytesBuffer: %v\n\n", bytesBuffer)
			if bufferLength > 0 {

				bytesBufferToChannel := make([]byte, bufferLength)
				copy(bytesBufferToChannel, bytesBuffer[:bufferLength])
				connectionDataChannel <- bytesBufferToChannel

			}
			if err != nil {
				connectionDataChannel <- nil
				el.error = err
				break
			}

		}

	}()

	return connectionDataChannel
}

func (el *Listen) pygocentrusDelay(inData []byte) (int, []byte) {
	//seelog.Debugf("%v%v were delayed by a pygocentrus attack: delay content", req.RemoteAddr, req.RequestURI)

	time.Sleep(time.Duration(el.inLineIntRange(el.Pygocentrus.Delay.Min, el.Pygocentrus.Delay.Max)) * time.Microsecond)

	return len(inData), inData

}

func (el *Listen) pygocentrusDontRespond(inData []byte) (int, []byte) {
	//seelog.Debugf("%v%v were eaten by a pygocentrus attack: dont respond", req.RemoteAddr, req.RequestURI)

	time.Sleep(time.Duration(el.inLineIntRange(el.Pygocentrus.Delay.Min, el.Pygocentrus.Delay.Max)) * time.Microsecond)
	return 0, nil

}

func (el *Listen) pygocentrusDeleteContent(inData []byte) (int, []byte) {
	//seelog.Debugf("%v%v were eaten by a pygocentrus attack: delete content", req.RemoteAddr, req.RequestURI)

	n := len(inData)
	inData = make([]byte, n)

	return n, inData

}

func (el *Listen) pygocentrusChangeContent(inData []byte) (int, []byte) {
	//seelog.Debugf("%v%v were eaten by a pygocentrus attack: change content", req.RemoteAddr, req.RequestURI)

	length := len(inData)
	forLength := el.Pygocentrus.ChangeContent.GetRandomByMaxMin(length)
	for i := 0; i != forLength; i += 1 {
		indexChange := el.Pygocentrus.ChangeContent.GetRandomByLength(length)
		inData = append(append(inData[:indexChange], byte(el.inLineRand().Intn(255))), inData[indexChange+1:]...)
	}

	return length, inData

}

func (el *Listen) inLineRand() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}

func (el *Listen) inLineIntRange(min, max int) int {
	return el.inLineRand().Intn(max-min) + min
}
