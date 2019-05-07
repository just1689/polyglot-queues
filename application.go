package main

import (
	"flag"
	"fmt"
	"github.com/just1689/polyglot-queues/queues"
	"github.com/nsqio/go-nsq"
	"github.com/sirupsen/logrus"
	"log"
	"sync"
	"time"
)

var Name = flag.String("name", "", "The name of the instance")
var AddrNSQD = flag.String("nsqd", "nsqd:4150", "The Address of nsq daemon")
var AddrNSQLookupD = flag.String("nsqld", "nsqlookupd:4161", "The Address of nsq lookup daemon")
var PubCount = flag.Int("pubs", 5, "The number of times to publish a message")
var stopSub chan bool

const Global = "global"
const Version = "1.0"

func main() {
	logrus.Println(*Name, " v", Version)
	var wg sync.WaitGroup
	wg.Add(1)
	stopSub := make(chan bool)

	flag.Parse()
	if Name == nil || *Name == "" {
		log.Fatal("No name provided")
	}

	logrus.Infoln("I am ", *Name)
	subscribeNSQ()
	publishNSQ()

	go func() {
		time.Sleep(1 * time.Minute)
		wg.Done()
	}()

	wg.Wait()
	stopSub <- true

}

func publishNSQ() {
	go func() {
		for i := 0; i < *PubCount; i++ {
			time.Sleep(2 * time.Second)
			config := nsq.NewConfig()
			w, _ := nsq.NewProducer(*AddrNSQD, config)
			msg := fmt.Sprint("Hello from ", *Name)
			err := w.Publish(Global, []byte(msg))
			if err != nil {
				log.Panic("Could not connect")
			} else {
				log.Println("Sent message")
			}
			w.Stop()
		}
	}()
}

func subscribeNSQ() {
	config := queues.Config{
		Address:       *AddrNSQD,
		LookupAddress: *AddrNSQLookupD,
		Topic:         Global,
		Channel:       Global,
		F:             nsqReceive,
		RemoteStopper: stopSub,
	}
	queues.Subscribe(config)
}

func nsqReceive(c *queues.Config, b []byte) {
	logrus.Println(*Name, " has received: ", string(b))
}
