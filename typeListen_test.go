package pygocentrus

import (
	"fmt"
	"log"
	"net"
	"time"
)

func ExampleListen_Listen() {

	listen := Listen{
		In: Connection{
			Address:  ":3001",
			Protocol: "tcp",
		},
		Out: Connection{
			Address:  ":3000",
			Protocol: "tcp",
		},
		Pygocentrus: Pygocentrus{
			Enabled: false,
		},
	}

	for {
		err := listen.Listen()
		if err != nil {
			log.Printf(err.Error())
		}
	}

	time.Sleep(time.Second * 1)

	conn, err := net.Dial("tcp", "127.0.0.1:50000")
	if err != nil {
		log.Fatal("could not connect to server: ", err)
	}
	defer conn.Close()

	conn.Write([]byte("ola mundo!"))

	time.Sleep(time.Second * 1)

	fmt.Println("data send ok")

	// Output:
	// data send ok
}
