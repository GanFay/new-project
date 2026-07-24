package main

import (
	"log"
	"time"

	"github.com/ganfay/split-notify/internal/client"
	"github.com/ganfay/split-notify/internal/consumer"
	"github.com/ganfay/split-notify/internal/processor"
)

func main() {
	time.Sleep(time.Second * 10)
	coreClient, err := client.NewCoreClient("app:50001")
	if err != nil {
		log.Fatalln("Main: Error init coreClient")
		return
	}
	defer coreClient.Close()

	proc := processor.NewProcessor(coreClient)

	cons, err := consumer.NewConsumer("amqp://guest:guest@rabbitmq:5672/", "test")
	if err != nil {
		log.Fatalf("Main: Consumer error. Error details: %v", err)
		return
	}
	defer cons.Close()
	cons.Start(proc)
}
