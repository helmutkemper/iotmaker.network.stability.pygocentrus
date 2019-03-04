package main

import (
	"fmt"
	"net"
)

type Listen struct {
	InAddress     string
	InProtocol    string
	OutAddress    string
	OutProtocol   string
	error         error
	inConnection  net.Conn
	outConnection net.Conn
}

func (el *Listen) Listen() error {
	var listener net.Listener
	var err error

	listener, err = net.Listen(el.InProtocol, el.InAddress)
	if err != nil {
		return err
	}

	for {
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

	el.outConnection, err = net.Dial(el.OutProtocol, el.OutAddress)
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
			fmt.Printf("%v\n\n", bytesBuffer)
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

func main() {

}
