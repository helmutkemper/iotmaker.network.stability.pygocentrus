package main

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	iotmakerdocker "github.com/helmutkemper/iotmaker.docker/v1.0.0"
	pygocentrus "github.com/helmutkemper/iotmaker.network.stability.pygocentrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"io"
	"net"
	"runtime/debug"
	"time"
)

func main() {
	var err error
	var pullStatusChannel = iotmakerdocker.NewImagePullStatusChannel()

	go func(c chan iotmakerdocker.ContainerPullStatusSendToChannel) {

		for {
			select {
			case status := <-c:
				//fmt.Printf("image pull status: %+v\n", status)

				if status.Closed == true {
					//fmt.Println("image pull complete!")
				}
			}
		}

	}(*pullStatusChannel)

	// stop and remove containers and garbage collector
	//err = toolsgarbagecollector.RemoveAllByNameContains("delete")
	//if err != nil {
	//  panic(string(debug.Stack()))
	//}
	//
	//p, _ := nat.NewPort("tcp", "27017")
	//_, err = factorycontainermongodb.NewSingleEphemeralInstanceMongoWithPort(
	//  "container_delete_before_test",
	//  p,
	//  factorycontainermongodb.KMongoDBVersionTag_latest,
	//  pullStatusChannel,
	//)
	//if err != nil {
	//  panic(string(debug.Stack()))
	//}

	var mongoClient *mongo.Client
	var ctx context.Context
	mongoClient, err = mongo.NewClient(options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
	if err != nil {
		panic(string(debug.Stack()))
	}

	ctx, _ = context.WithTimeout(context.Background(), 120*time.Second)
	err = mongoClient.Connect(ctx)
	if err != nil {
		panic(string(debug.Stack()))
	}

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 120*time.Second)
	err = mongoClient.Ping(ctx, readpref.Primary())
	if err != nil {
		panic(string(debug.Stack()))
	}
	defer cancel()

	err = mongoClient.Disconnect(ctx)
	if err != nil {
		panic(string(debug.Stack()))
	}

	l := pygocentrus.Listen{
		In: pygocentrus.Connection{
			Address:  ":27016",
			Protocol: "tcp",
		},
		Out: pygocentrus.Connection{
			Address:  ":27017",
			Protocol: "tcp",
		},
		Pygocentrus: pygocentrus.Pygocentrus{
			Enabled: true,
			Delay: pygocentrus.RateMaxMin{
				Rate: 1.0,
				Min:  int(time.Millisecond * 100),
				Max:  int(time.Millisecond * 500),
			},
			DontRespond: pygocentrus.RateMaxMin{
				Rate: 0,
				Min:  0,
				Max:  0,
			},
			ChangeLength: 0,
			ChangeContent: pygocentrus.ChangeContent{
				ChangeRateMin:  0,
				ChangeRateMax:  0,
				ChangeBytesMin: 0,
				ChangeBytesMax: 0,
				Rate:           0,
			},
			DeleteContent: 0,
		},
	}

	go func() {
		err = l.Listen()
		if err != nil {
			panic(err)
		}
	}()

	//go func() {
	//  var err error
	//  err = dial("tcp4", "127.0.0.1:27017", "127.0.0.1:27016")
	//  if err != nil {
	//    panic(err)
	//  }
	//}()

	//go proxy(27016, 27017)

	//var conn net.Conn
	//conn, err = net.Dial("tcp", ":27016")
	//if err != nil {
	//  panic(err)
	//}
	//var data = make([]byte, 0)
	//for i := 0; i != 1024; i += 1 {
	//  data = append(data, byte(i))
	//}
	//data = append(data, byte(0x00))
	//
	//_, err = conn.Write(data)
	//if err != nil {
	//  panic(err)
	//}
	//conn.Close()
	//
	//time.Sleep(time.Second*500)
	//os.Exit(0)

	fmt.Printf("conexão\n")

	mongoClient, err = mongo.NewClient(options.Client().ApplyURI("mongodb://127.0.0.1:27016"))
	if err != nil {
		panic(string(debug.Stack()))
	}

	err = mongoClient.Connect(ctx)
	if err != nil {
		panic(string(debug.Stack()))
	}

	ctx, cancel = context.WithTimeout(context.Background(), 120*time.Second)
	err = mongoClient.Ping(ctx, readpref.Primary())
	if err != nil {
		panic(string(debug.Stack()))
	}

	type Trainer struct {
		Name string
		Age  int
		City string
	}
	collection := mongoClient.Database("test").Collection("trainers")
	ash := Trainer{"Ash", 10, "Pallet Town"}
	_, err = collection.InsertOne(context.TODO(), ash)
	if err != nil {
		panic(err)
	}
	fmt.Printf("fim\n")
}

func proxy(inPort, outPort int) {
	var err error

	var network = "tcp"
	var inAddress *net.TCPAddr = &net.TCPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: inPort,
	}
	var outAddress *net.TCPAddr = &net.TCPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: outPort,
	}

	var inTcpListener net.Listener
	var inTcpConn net.Conn
	var outTcpConn net.Conn

	var inCh = make(chan error, 1)
	var outCh = make(chan error, 1)

	fmt.Printf("in address: ':%v'\n", inAddress.Port)
	fmt.Printf("out address: ':%v'\n", outAddress.Port)

	inAddress, err = net.ResolveTCPAddr(network, fmt.Sprintf("%v", inAddress))
	if err != nil {
		debug.PrintStack()
		panic(err)
	}

	inTcpListener, err = net.Listen(network, fmt.Sprintf("%v", inAddress))
	if err != nil {
		debug.PrintStack()
		panic(err)
	}

	inTcpConn, err = inTcpListener.Accept()
	if err != nil {
		debug.PrintStack()
		panic(err)
	}

	outTcpConn, err = net.Dial(network, fmt.Sprintf("127.0.0.1:%v", outAddress.Port))
	if err != nil {
		debug.PrintStack()
		panic(err)
	}

	inCh = makeChannelFromConnection(inTcpConn, outTcpConn)
	if err != nil {
		debug.PrintStack()
		panic(err)
	}

	outCh = makeChannelFromConnection(outTcpConn, inTcpConn)

	for {
		select {
		case err := <-inCh:
			if err != nil {
				debug.PrintStack()
				panic(err)
			}

		case err := <-outCh:
			if err != nil {
				debug.PrintStack()
				panic(err)
			}
		}
	}

}

func makeChannelFromConnection(inConn, outConn net.Conn) (ch chan error) {

	ch = make(chan error, 1)
	go func() {

		var bytesBuffer = make([]byte, 2048)
		var bufferLength int
		var err error

		for {

			bufferLength, err = inConn.Read(bytesBuffer)
			if err != nil && errors.Is(err, io.EOF) {
				fmt.Printf("EOF\n")

				ch <- err
				return
			}

			if err != nil {
				ch <- err
				return
			}

			buffer := make([]byte, bufferLength, bufferLength)
			copy(buffer, bytesBuffer)
			fmt.Printf("buffer len: %v\n", bufferLength)
			fmt.Printf("%v\n\n", hex.Dump(bytesBuffer[:bufferLength]))
			_, err = outConn.Write(buffer)
			if err != nil {
				ch <- err
				return
			}
		}

	}()

	return
}
