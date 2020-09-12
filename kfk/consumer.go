package kfk

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Shopify/sarama"
	"github.com/bsm/sarama-cluster" // support automatic consumer-group rebalancing and offset tracking
)

// 支持brokers cluster的消费者
func ClusterConsumer(ctx context.Context, brokers, topics []string, groupId string, id int) {
	cfg := cluster.NewConfig()
	cfg.Consumer.Return.Errors = true
	cfg.Group.Return.Notifications = true
	cfg.Consumer.Offsets.Initial = sarama.OffsetNewest
	cfg.Consumer.Offsets.CommitInterval = time.Second
	// cfg.Consumer.Offsets.AutoCommit.Enable = true
	// cfg.Consumer.Offsets.AutoCommit.Interval = time.Second

	consumer, err := cluster.NewConsumer(brokers, groupId, topics, cfg)
	if err != nil {
		log.Printf("%s: sarama.NewSyncProducer err, message=%s \n", groupId, err.Error())
		return
	}
	defer func() {
		_ = consumer.Close()
	}()

	// consume errors
	go func() {
		for err := range consumer.Errors() {
			log.Printf("%s:Error: %s\n", groupId, err.Error())
		}
	}()

	// consume notifications
	go func() {
		for ntf := range consumer.Notifications() {
			log.Printf("%s:Rebalanced: %+v \n", groupId, ntf)
		}
	}()

	// consume messages, watch signals
	var success int
_EXIT:
	for {
		select {
		case msg, ok := <-consumer.Messages():
			if ok {
				_, _ = fmt.Fprintf(os.Stdout, "%d >> %s:%s/%d/%d\t%s\t%s\n", id, groupId, msg.Topic, msg.Partition, msg.Offset, msg.Key, msg.Value)
				consumer.MarkOffset(msg, "") // mark message as processed
				success++
			}
		case <-ctx.Done():
			break _EXIT
		}
	}
	_, _ = fmt.Fprintf(os.Stdout, "%s consume %d messages \n", groupId, success)
}
