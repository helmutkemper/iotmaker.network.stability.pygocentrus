package pygocentrus

import "net"

func (el *Listener) AttackListener() error {
	var listener net.Listener
	var err error

	listener, err = net.Listen(el.In.Protocol.String(), el.In.Address)
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
