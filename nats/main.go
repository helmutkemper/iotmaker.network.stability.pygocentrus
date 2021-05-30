package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	dockerBuilder "github.com/helmutkemper/iotmaker.docker.builder"
	pygocentrus "github.com/helmutkemper/iotmaker.network.stability.pygocentrus"
	"github.com/helmutkemper/util"
	"github.com/nats-io/nats.go"
	"log"
	"strconv"
	"sync"
	"time"
)

type NatsInfo struct {
	ServerId   string `json:"server_id"`
	ServerName string `json:"server_name"`
	Version    string `json:"version"`
	Proto      int    `json:"proto"`
	GitCommit  string `json:"git_commit"`
	GoVersion  string `json:"go"`
	Host       string `json:"host"`
	Port       int    `json:"port"`
	Headers    bool   `json:"headers"`
	PayloadMax int64  `json:"max_payload"`
	ClientId   int64  `json:"client_id"`
	ClientIp   string `json:"client_ip"`
	Cluster    string `json:"cluster"`
}

type NatsConnection struct {
	Verbose      bool   `json:"verbose"`
	Pedantic     bool   `json:"pedantic"`
	TlsRequired  bool   `json:"tls_required"`
	Name         string `json:"name"`
	Lang         string `json:"lang"`
	Version      string `json:"version"`
	Protocol     int64  `json:"protocol"`
	Echo         bool   `json:"echo"`
	Headers      bool   `json:"headers"`
	NoResponders bool   `json:"no_responders"`
}

type Filter struct{}

type ParserFunc struct {
	Filter
}

func (e Filter) filterUpToTheByteSlice(data, end []byte) (found bool, filtered, leftover []byte) {
	var lenData = len(data)
	var dataCopy = make([]byte, lenData)
	copy(dataCopy, data)

	var lenEnd = len(end)
	var cursorEnd, cursorLeftover int

	var indexEnd = bytes.Index(dataCopy, end)
	if indexEnd == -1 {
		return
	}

	cursorEnd = indexEnd
	cursorLeftover = cursorEnd + lenEnd

	filtered = make([]byte, cursorEnd)
	copy(filtered, dataCopy[0:cursorEnd])

	leftover = make([]byte, lenData-cursorLeftover)
	copy(leftover, dataCopy[cursorLeftover:])

	found = true
	return
}

func (e Filter) filterLength(data []byte, length int) (found bool, filtered, leftover []byte) {
	var lenData = len(data)
	var dataCopy = make([]byte, lenData)
	copy(dataCopy, data)

	var cursorEnd, cursorLeftover int

	cursorEnd = length
	cursorLeftover = cursorEnd

	filtered = make([]byte, cursorEnd)
	copy(filtered, dataCopy[0:cursorEnd])

	leftover = make([]byte, lenData-cursorLeftover)
	copy(leftover, dataCopy[cursorLeftover:])

	found = true
	return
}

func (e Filter) filter(data, start, end []byte) (found bool, filtered, leftover []byte) {
	var lenData = len(data)
	var dataCopy = make([]byte, lenData)
	copy(dataCopy, data)

	var lenStart = len(start)
	var lenEnd = len(end)
	var cursorStart, cursorEnd, cursorLeftover int

	var indexStart = bytes.Index(dataCopy, start)
	if indexStart == -1 {
		return
	}

	cursorStart = indexStart + lenStart

	var indexEnd = bytes.Index(dataCopy[cursorStart:], end)
	if indexEnd == -1 {
		return
	}

	cursorEnd = cursorStart + indexEnd
	cursorLeftover = cursorEnd + lenEnd

	filtered = make([]byte, cursorEnd-cursorStart)
	copy(filtered, dataCopy[cursorStart:cursorEnd])

	//                      [:indexStart]+(data[cursorLeftover:])
	leftover = make([]byte, indexStart+(lenData-cursorLeftover))
	copy(leftover, append(dataCopy[:indexStart], dataCopy[cursorLeftover:]...))

	found = true
	return
}

func (e ParserFunc) FilterInfo(data []byte) (found bool, infoJson, leftover []byte, err error) {
	var lenData = len(data)
	var dataCopy = make([]byte, lenData)
	copy(dataCopy, data)

	found, infoJson, leftover = e.filter(dataCopy, []byte("INFO "), []byte(" \r\n"))
	return
}

func (e ParserFunc) FilterConnect(data []byte) (found bool, connectionJson, leftover []byte, err error) {
	var lenData = len(data)
	var dataCopy = make([]byte, lenData)
	copy(dataCopy, data)

	found, connectionJson, leftover = e.filter(dataCopy, []byte("CONNECT "), []byte("\r\n"))
	return
}

func (e ParserFunc) FilterPing(data []byte) (found bool, ping, leftover []byte, err error) {
	var lenData = len(data)
	var dataCopy = make([]byte, lenData)
	copy(dataCopy, data)

	found, ping, leftover = e.filter(dataCopy, []byte("PING"), []byte("\r\n"))
	ping = []byte("PING")
	return
}

func (e ParserFunc) FilterPong(data []byte) (found bool, pong, leftover []byte, err error) {
	var lenData = len(data)
	var dataCopy = make([]byte, lenData)
	copy(dataCopy, data)

	found, pong, leftover = e.filter(dataCopy, []byte("PONG"), []byte("\r\n"))
	pong = []byte("PONG")
	return
}

func (e ParserFunc) FilterSubscribe(data []byte) (found bool, id int64, subject []byte, leftover []byte, err error) {
	var lenData = len(data)
	var dataCopy = make([]byte, lenData)
	copy(dataCopy, data)

	var idAsByte []byte
	found, subject, leftover = e.filter(dataCopy, []byte("SUB "), []byte(" "))
	if found == false {
		return
	}

	_, idAsByte, leftover = e.filter(leftover, []byte(" "), []byte("\r\n"))

	id, err = strconv.ParseInt(string(idAsByte), 10, 64)
	return
}

func (e ParserFunc) FilterPublish(data []byte) (found bool, subject, message, leftover []byte, err error) {
	var lenData = len(data)
	var dataCopy = make([]byte, lenData)
	copy(dataCopy, data)

	var lengthAsByte []byte
	var length int64
	found, subject, leftover = e.filter(dataCopy, []byte("PUB "), []byte(" "))
	if found == false {
		return
	}

	_, lengthAsByte, leftover = e.filterUpToTheByteSlice(leftover, []byte("\r\n"))

	length, err = strconv.ParseInt(string(lengthAsByte), 10, 64)

	_, message, leftover = e.filterLength(leftover, int(length))

	_, _, leftover = e.filterUpToTheByteSlice(leftover, []byte("\r\n"))
	return
}

func (e ParserFunc) FilterUnsubscribe(data []byte) (found bool, id int64, leftover []byte, err error) {
	var lenData = len(data)
	var dataCopy = make([]byte, lenData)
	copy(dataCopy, data)

	var idAsBute []byte
	found, idAsBute, leftover = e.filter(dataCopy, []byte("UNSUB "), []byte(" \r\n"))
	if found == false {
		return
	}

	id, err = strconv.ParseInt(string(idAsBute), 10, 64)
	return
}

func (e ParserFunc) FilterMessage(data []byte) (found bool, id int64, subject []byte, message []byte, leftover []byte, err error) {
	var lenData = len(data)
	var dataCopy = make([]byte, lenData)
	copy(dataCopy, data)

	var idAsByte []byte
	var lengthAsByte []byte
	var length int64
	found, subject, leftover = e.filter(dataCopy, []byte("MSG "), []byte(" "))
	if found == false {
		return
	}

	found, idAsByte, leftover = e.filterUpToTheByteSlice(leftover, []byte(" "))
	if found == false {
		return
	}

	id, err = strconv.ParseInt(string(idAsByte), 10, 64)
	if err != nil {
		return
	}

	_, lengthAsByte, leftover = e.filterUpToTheByteSlice(leftover, []byte("\r\n"))

	length, err = strconv.ParseInt(string(lengthAsByte), 10, 64)
	if err != nil {
		return
	}

	_, message, leftover = e.filterLength(leftover, int(length))

	_, _, leftover = e.filterUpToTheByteSlice(leftover, []byte("\r\n"))
	return
}

func (e ParserFunc) Parser(data []byte, direction string) (dataSize int, err error) {
	dataSize = len(data)
	var dataParser = make([]byte, dataSize)
	copy(dataParser, data)
	fmt.Println("")
	//fmt.Println("")

	var infoJson []byte
	var connectJson []byte
	var info NatsInfo
	var connection NatsConnection

	// \r 0x0D
	// \n 0x0A
	fmt.Printf("direction: %v\n%v\n", direction, hex.Dump(data))

	for {
		var l = len(dataParser)

		if l == 0 {
			break
		}

		if bytes.HasPrefix(dataParser, []byte("INFO ")) == true {
			_, infoJson, dataParser, err = e.FilterInfo(dataParser)
			if err != nil {
				return
			}

			err = json.Unmarshal(infoJson, &info)
			if err != nil {
				return
			}
		}

		if bytes.HasPrefix(dataParser, []byte("CONNECT ")) == true {
			_, connectJson, dataParser, err = e.FilterConnect(dataParser)
			if err != nil {
				return
			}

			log.Printf("%s", connectJson)
			err = json.Unmarshal(connectJson, &connection)
			if err != nil {
				return
			}
		}

		if bytes.HasPrefix(dataParser, []byte("PING\r\n")) == true {
			_, _, dataParser, err = e.FilterPing(dataParser)
			if err != nil {
				return
			}
		}

		if bytes.HasPrefix(dataParser, []byte("PONG\r\n")) == true {
			_, _, dataParser, err = e.FilterPong(dataParser)
			if err != nil {
				return
			}
		}

		if bytes.HasPrefix(dataParser, []byte("SUB ")) == true {
			_, _, _, dataParser, err = e.FilterSubscribe(dataParser)
			if err != nil {
				return
			}
		}

		if bytes.HasPrefix(dataParser, []byte("PUB ")) == true {
			_, _, _, dataParser, err = e.FilterPublish(dataParser)
			if err != nil {
				return
			}
		}

		if bytes.HasPrefix(dataParser, []byte("UNSUB ")) == true {
			_, _, dataParser, err = e.FilterUnsubscribe(dataParser)
			if err != nil {
				return
			}
		}

		if bytes.HasPrefix(dataParser, []byte("MSG ")) == true {
			_, _, _, _, dataParser, err = e.FilterMessage(dataParser)
			if err != nil {
				return
			}
		}

		if l == len(dataParser) {
			log.Print("falhou!!!!!!!!!!!!!!!!!!")
			break
		}
	}

	return
}

func main() {
	//installNats()
	proxy()
	natsTest()
}

func installNats() {

	dockerBuilder.GarbageCollector()

	var err error
	// Prepara a instalação do container para a imagem nats:latest
	var natsDocker = dockerBuilder.ContainerBuilder{}
	natsDocker.SetPrintBuildOnStrOut()
	// Aponta o gerenciador de rede [opcional]
	// Como o gateway é 10.0.0.1, o primeiro container gerado fica no endereço 10.0.0.2
	//natsDocker.SetNetworkDocker(netDocker)
	// Determina o nome da imagem a ser usada
	natsDocker.SetImageName("nats:latest")
	// Determina o nome do container a ser criado
	natsDocker.SetContainerName("nats_delete_after_test")

	// Você pode expor a porta 4222 para o fora da rede
	//natsDocker.AddPortToOpen("4222")

	// Você pode trocar uma porta 4222 para 4200 e a expor para fora da rede
	natsDocker.AddPortToChange("4222", "4223")

	// Espera pelo texto abaixo no log do container antes de proceguir
	natsDocker.SetWaitStringWithTimeout(
		"Listening for route connections on 0.0.0.0:6222",
		40*time.Second,
	)

	// Inicializa o objeto depois de todas as configurações feitas
	err = natsDocker.Init()
	if err != nil {
		util.TraceToLog()
		log.Printf("Error: %v", err.Error())
		return
	}

	// Baixa a imagem caso a mesma não exista e deve ser usado apenas em caso de imagens públicas
	err = natsDocker.ImagePull()
	if err != nil {
		util.TraceToLog()
		log.Printf("Error: %v", err.Error())
		return
	}

	// Monta o container a partir da imagem baixada por ImagePull() e definida em SetImageName()
	err = natsDocker.ContainerBuildFromImage()
	if err != nil {
		util.TraceToLog()
		log.Printf("Error: %v", err.Error())
		return
	}
}

func proxy() {
	var p pygocentrus.ParserInterface = &ParserFunc{}

	var proxy pygocentrus.Proxy
	proxy.SetBufferSize(1024)
	//proxy.SetDelayMillesecond(5999, 6000)
	proxy.SetParserFunction(p)

	go func() {
		var err error
		err = proxy.Proxy("localhost:4222", "localhost:4223")
		if err != nil {
			panic(err)
		}
	}()
	time.Sleep(time.Second)
}

func natsTest() {
	nc, err := nats.Connect("nats://127.0.0.1:4222", nats.Timeout(time.Second*20))
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	done := sync.WaitGroup{}
	done.Add(1)

	count := 0
	errCh := make(chan error, 1)

	msgAfterDrain := "not this one"

	// Just to not collide using the demo server with other users.
	subject := nats.NewInbox()
	//subject = "1234567890"

	// This callback will process each message slowly
	sub, err := nc.Subscribe(subject, func(m *nats.Msg) {
		if string(m.Data) == msgAfterDrain {
			errCh <- fmt.Errorf("Should not have received this message")
			return
		}

		//log.Printf("message from sample code: %s", m.Data)

		time.Sleep(100 * time.Millisecond)
		count++
		if count == 2 {
			done.Done()
		}
	})

	// Send 2 messages
	for i := 0; i < 2; i++ {
		nc.Publish(subject, []byte("hello"))
	}

	// Call Drain on the subscription. It unsubscribes but
	// wait for all pending messages to be processed.
	if err := sub.Drain(); err != nil {
		log.Fatal(err)
	}

	// Send one more message, this message should not be received
	nc.Publish(subject, []byte(msgAfterDrain))

	// Wait for the subscription to have processed the 2 messages.
	done.Wait()

	// Now check that the 3rd message was not received
	select {
	case e := <-errCh:
		log.Fatal(e)
	case <-time.After(200 * time.Millisecond):
		// OK!
	}
}
