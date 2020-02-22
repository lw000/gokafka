// goKafka project main.go
package main

import (
	"context"
	"gokafka/kfk"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())
	var Address = []string{"192.168.0.105:9092"}

	go kfk.AsyncProducer(ctx, Address)

	go kfk.ClusterConsumer(ctx, Address, []string{"test"}, "test-consumer-group")

	log.Println(<-signals)
	cancel()

	time.Sleep(time.Second)
}
