package listener

import "net"

func (el *Listener) makeChannelFromConnection(conn net.Conn) chan []byte {
	var connectionDataChannel chan []byte

	connectionDataChannel = make(chan []byte)

	go func() {

		var bytesBuffer = make([]byte, 1024)
		var bufferLength int
		var err error

		for {

			bufferLength, err = conn.Read(bytesBuffer)

			if el.attack != nil {
				bufferLength, bytesBuffer = el.attack(bytesBuffer)
			}

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
