package main

import (
	"encoding/hex"
	"fmt"
	dockerBuilder "github.com/helmutkemper/iotmaker.docker.builder"
	pygocentrus "github.com/helmutkemper/iotmaker.network.stability.pygocentrus"
	"github.com/helmutkemper/util"
	"github.com/nats-io/nats.go"
	"log"
	"sync"
	"time"
)

type ParserFunc struct{}

func (e ParserFunc) Parser(data []byte, direction string) (dataSize int, err error) {
	fmt.Printf("direction: %v\n%v\n", direction, hex.Dump(data))
	dataSize = len(data)
	return
}

func main() {
	installNats()
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
	proxy.SetBufferSize(32 * 1024)
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
	nc, err := nats.Connect("nats://127.0.0.1:4222")
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

	// This callback will process each message slowly
	sub, err := nc.Subscribe(subject, func(m *nats.Msg) {
		if string(m.Data) == msgAfterDrain {
			errCh <- fmt.Errorf("Should not have received this message")
			return
		}

		log.Printf("message: %s", m.Data)

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
