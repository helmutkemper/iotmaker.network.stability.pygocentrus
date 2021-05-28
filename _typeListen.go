package pygocentrus

import (
  "errors"
  "fmt"
  "github.com/helmutkemper/util"
  "log"
  "math/rand"
  "net"
  "sync"
  "syscall"
  "time"
)

type pygocentrusListenFunc func(inData []byte, length int) (int, []byte)

type data struct {
	channel chan bool
	buffer  [][]byte
	length  []int
}

func (el *data) init() {
	el.buffer = make([][]byte, 0)
	el.channel = make(chan bool, 1)
	el.length = make([]int, 0)
}

type Listen struct {
	In            Connection
	Out           Connection
	Pygocentrus   Pygocentrus
	error         error
	inConnection  *net.TCPConn
	inData        data
	outData       data
	outConnection *net.TCPConn
	attack        pygocentrusListenFunc
	mutex         sync.Mutex
	ticker        *time.Ticker
}

func (el *Listen) Listen() (err error) {

  el.inData.init()
  el.outData.init()

  var listener *net.TCPListener

  listener, err = net.ListenTCP("tcp", el.In.A)
  if err != nil {
    log.Printf("err: %v", err.Error())
    util.TraceToLog()
    return
  }

  go el.handle()

  for {
    el.inConnection, err = listener.AcceptTCP()
    if err != nil {
      log.Printf("err: %v", err.Error())
      util.TraceToLog()
      return
    }

    //el.ticker = time.NewTicker(11000 * time.Millisecond)

    el.outConnection, err = net.DialTCP("tcp", nil, el.Out.A)
    if err != nil {
      el.error = err
      log.Printf("err: %v", err.Error())
      util.TraceToLog()
      return
    }

    el.makeChannelFromInDataConnection(el.inConnection, &el.inData, "in")
    el.makeChannelFromInDataConnection(el.outConnection, &el.outData, "out")

    if el.error != nil {
      log.Printf("err: %v", el.error.Error())
      util.TraceToLog()
      return el.error
    }
  }

  return
}

func (el *Listen) handle() {
  log.Print("handle()")
  defer log.Print("handle().end")

	var err error

	for {
		select {
		case <-el.inData.channel:
			el.mutex.Lock()
			for {
				if len(el.inData.buffer) == 0 {
					el.mutex.Unlock()
					break
				}

				_, err = el.outConnection.Write(el.inData.buffer[0])
				if err != nil {
          if !errors.Is(err, syscall.EPIPE) {
            el.error = err
            log.Printf("err: %v", err.Error())
            util.TraceToLog()
            return
          }
				}

				el.inData.buffer = el.inData.buffer[1:]
			}

		case <-el.outData.channel:
			el.mutex.Lock()
			for {
				if len(el.outData.buffer) == 0 {
					el.mutex.Unlock()
					break
				}

				_, err = el.inConnection.Write(el.outData.buffer[0])
				if err != nil {
          if !errors.Is(err, syscall.EPIPE) {
            el.error = err
            log.Printf("err: %v", err.Error())
            util.TraceToLog()
            return
          }
				}

				el.outData.buffer = el.outData.buffer[1:]
			}
		}
	}
}

func (el *Listen) makeChannelFromInDataConnection(conn *net.TCPConn, data *data, direction string) {
  log.Print("makeChannelFromInDataConnection()")
  defer log.Print("makeChannelFromInDataConnection().end")
	go func() {

		var bufferLength int
		var err error

		err = conn.SetKeepAlive(true)
		if err != nil {
      log.Printf("err: %v", err.Error())
      util.TraceToLog()
      panic(err)
		}

		for {
			var buffer = make([]byte, 2048)
			bufferLength, err = conn.Read(buffer)
      if !errors.Is(err, syscall.EPIPE) {
        if err != nil && err.Error() != "EOF" {
          el.error = err
          log.Printf("err: %v", err.Error())
          util.TraceToLog()
          break
        }

        if err != nil && err.Error() == "EOF" {
          break
        }
      }

			//fmt.Printf("%v\n%v\n", direction, hex.Dump(buffer[:bufferLength]))

			if cap(data.buffer) == 0 {
				data.buffer = make([][]byte, 0)
			}
			data.buffer = append(data.buffer, buffer[:bufferLength])

			if cap(data.length) == 0 {
				data.length = make([]int, 0)
			}
			data.length = append(data.length, bufferLength)

			if len(data.channel) == 0 {
				data.channel <- true
			}
		}
	}()
}

func (el *Listen) pygocentrusDelay(inData []byte, length int) (int, []byte) {
	//seelog.Debugf("%v%v were delayed by a pygocentrus attack: delay content", req.RemoteAddr, req.RequestURI)

	time.Sleep(time.Duration(el.inLineIntRange(el.Pygocentrus.Delay.Min, el.Pygocentrus.Delay.Max)))

	return length, inData

}

func (el *Listen) pygocentrusDontRespond(inData []byte, length int) (int, []byte) {
	fmt.Printf("pygocentrus attack: dont respond\n")

	time.Sleep(time.Duration(el.inLineIntRange(el.Pygocentrus.Delay.Min, el.Pygocentrus.Delay.Max)) * time.Microsecond)
	return 0, nil

}

func (el *Listen) pygocentrusDeleteContent(inData []byte, length int) (int, []byte) {
	//seelog.Debugf("%v%v were eaten by a pygocentrus attack: delete content", req.RemoteAddr, req.RequestURI)

	inData = make([]byte, length)

	return length, inData

}

func (el *Listen) pygocentrusChangeContent(inData []byte, length int) (int, []byte) {
	//seelog.Debugf("%v%v were eaten by a pygocentrus attack: change content", req.RemoteAddr, req.RequestURI)

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

func (el *Listen) SelectAttack() {

	var randAttack int
	var list = make([]pygocentrusListenFunc, 0)

	if el.Pygocentrus.Delay.Rate != 0.0 {

		if el.Pygocentrus.Delay.Rate >= el.inLineRand().Float64() {
			list = append(list, el.pygocentrusDelay)
		}

	}

	if el.Pygocentrus.DontRespond.Rate != 0.0 {

		if el.Pygocentrus.DontRespond.Rate >= el.inLineRand().Float64() {
			list = append(list, el.pygocentrusDontRespond)
		}

	}

	if el.Pygocentrus.DeleteContent != 0.0 {

		if el.Pygocentrus.DeleteContent >= el.inLineRand().Float64() {
			list = append(list, el.pygocentrusDeleteContent)
		}

	}

	if el.Pygocentrus.ChangeContent.Rate != 0.0 {

		if el.Pygocentrus.ChangeContent.Rate >= el.inLineRand().Float64() {
			list = append(list, el.pygocentrusChangeContent)
		}

	}

	listLength := len(list)
	if listLength != 0 {
		el.Pygocentrus.SetAttack()
		randAttack = el.inLineRand().Intn(len(list))
		el.attack = list[randAttack]
	}
}
