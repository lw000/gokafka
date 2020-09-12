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

	// dataQueue := diskqueue.New("1111", "./", 1024*1024*100, 100000, 10000000, 100, time.Second*time.Duration(10), func(lvl diskqueue.LogLevel, f string, args ...interface{}) {
	// 	log.Println(args)
	// })
	// _ = dataQueue.Put([]byte("123123123123123123"))
	// _ = dataQueue.Close()

	ctx, cancel := context.WithCancel(context.Background())
	var Address = []string{"192.168.0.115:9092"}

	go kfk.AsyncProducer(ctx, "test", Address)

	go kfk.ClusterConsumer(ctx, Address, []string{"test"}, "test-consumer-group", 1)
	// go kfk.ClusterConsumer(ctx, Address, []string{"test"}, "test-consumer-group", 2)

	log.Println(<-signals)
	cancel()

	time.Sleep(time.Second)
}
