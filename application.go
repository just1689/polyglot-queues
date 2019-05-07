package main

import (
	"flag"
	"fmt"
	"github.com/just1689/polyglot-queues/queues"
	"github.com/nsqio/go-nsq"
	"log"
	"sync"
	"time"
)

var Name = flag.String("name", "", "The name of the instance")
var AddrNSQD = flag.String("nsqd", "nsqd:4150", "The Address of nsq daemon")
var AddrNSQLookupD = flag.String("nsqld", "nsqlookupd:4161", "The Address of nsq lookup daemon")
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
	go subscribeNSQ()
	go publishNSQ()

	go func() {
		time.Sleep(1 * time.Minute)
		wg.Done()
	}()

	wg.Wait()
	stopSub <- true

}

func publishNSQ() {
	config := nsq.NewConfig()
	w, _ := nsq.NewProducer(*AddrNSQD, config)
	for {
		time.Sleep(1 * time.Second)
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
		F:             nsqReceive,
		RemoteStopper: stopSub,
	}
	queues.Subscribe(config)
}

func nsqReceive(c *queues.Config, b []byte) {
	fmt.Println(*Name, "has received:", string(b))
}
