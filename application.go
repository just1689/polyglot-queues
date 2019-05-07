package main

import (
	"flag"
	"fmt"
	"github.com/just1689/polyglot-queues/queues"
	"github.com/nats-io/go-nats"
	"github.com/nsqio/go-nsq"
	"log"
	"sync"
	"time"
)

var Name = flag.String("name", "", "The name of the instance")
var AddrNSQD = flag.String("nsqd", "nsqd:4150", "The Address of nsq daemon")
var AddrNSQLookupD = flag.String("nsqld", "nsqlookupd:4161", "The Address of nsq lookup daemon")
var AddrNats = flag.String("nats", "nats://nats-cluster-node-2:4222", "The address of the NATS")
var stopSub chan bool

const Global = "global"
const Version = "2.0"

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	stopSub := make(chan bool)

	flag.Parse()

	if Name == nil || *Name == "" {
		log.Fatal("No name provided")
	}

	fmt.Println(*Name, "@ Version", Version)

	//NSQ
	go subscribeNSQ()
	go publishNSQ()

	//NATS
	go subscribeNATS()
	go publishNATS()

	go func() {
		time.Sleep(30 * time.Second)
		wg.Done()
	}()

	wg.Wait()
	stopSub <- true

}

func publishNSQ() {
	config := nsq.NewConfig()
	w, _ := nsq.NewProducer(*AddrNSQD, config)
	for {
		time.Sleep(500 * time.Millisecond)
		msg := fmt.Sprint("Hello from ", *Name)
		err := w.Publish(Global, []byte(msg))
		if err != nil {
			log.Panic("Could not connect")
		}
	}
	w.Stop()

}

func subscribeNSQ() {
	config := queues.Config{
		Address:       *AddrNSQD,
		LookupAddress: *AddrNSQLookupD,
		Topic:         Global,
		Channel:       *Name,
		F:             handlerNSQ,
		RemoteStopper: stopSub,
	}
	queues.Subscribe(config)
}

func handlerNSQ(c *queues.Config, b []byte) {
	fmt.Println(*Name, "has received:", string(b), "(NSQ)")
}

///
///		NATS
///

func publishNATS() {
	nc, _ := nats.Connect(*AddrNats)
	for {
		time.Sleep(100 * time.Millisecond)
		nc.Publish(Global, []byte(fmt.Sprint("Hello from ", *Name)))
	}
	nc.Close()
}

func handleNATS(m *nats.Msg) {
	fmt.Println(*Name, "has received:", string(m.Data), "(NATS)")
}

func subscribeNATS() {
	nc, _ := nats.Connect(*AddrNats)
	_, err := nc.Subscribe(Global, handleNATS)
	if err != nil {
		fmt.Println(err)
	}
}
