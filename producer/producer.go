package main

import (
	"context"
	"fmt"
	"github.com/Shopify/sarama"
	"log"
	"os"
	"time"
)

func AsyncProducer(ctx context.Context, topic string, address []string) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	// config.Producer.Partitioner = 默认为message的hash
	p, err := sarama.NewAsyncProducer(address, config)
	if err != nil {
		log.Printf("sarama.NewSyncProducer err, message=%s \n", err)
		return
	}

	var enqueued, successes, errors int

	// 发送成功message计数
	go func() {
		for range p.Successes() {
			successes++
		}
	}()

	// 发送失败计数
	go func() {
		for err := range p.Errors() {
			log.Printf("%+v 发送失败，err：%s\n", err.Msg, err.Err)
			errors++
		}
	}()

	var messageIndex = 0
_EXIT:
	for {
		messageIndex++
		value := fmt.Sprintf("this is a message. index=%d", messageIndex)
		msg := &sarama.ProducerMessage{
			Topic: topic,
			Key:   sarama.ByteEncoder([]byte("login")),
			Value: sarama.ByteEncoder(value),
		}

		select {
		case p.Input() <- msg: // 发送消息
			enqueued++
		case <-ctx.Done():
			break _EXIT
		}
		time.Sleep(time.Millisecond * time.Duration(10))
	}

	_, _ = fmt.Fprintf(os.Stdout, "发送数=%d，发送成功数=%d，发送失败数=%d \n", enqueued, successes, errors)
}
