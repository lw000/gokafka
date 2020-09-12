// goKafka project main.go
package main

import (
	"context"
	"os"
	"os/signal"
	"time"
)

func main() {
	defer func() {
		time.Sleep(time.Second)
	}()
	// dataQueue := diskqueue.New("1111", "./", 1024*1024*100, 100000, 10000000, 100, time.Second*time.Duration(10), func(lvl diskqueue.LogLevel, f string, args ...interface{}) {
	// 	log.Println(args)
	// })
	// _ = dataQueue.Put([]byte("123123123123123123"))
	// _ = dataQueue.Close()

	ctx, cancel := context.WithCancel(context.Background())
	var Address = []string{"192.168.0.115:9092", "192.168.0.115:9093", "192.168.0.115:9094"}

	go AsyncProducer(ctx, "test", Address)

	go AsyncProducer(ctx, "test-1", Address)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	<-signals

	cancel()
}
