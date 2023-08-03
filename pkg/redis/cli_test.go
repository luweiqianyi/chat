package redis

import (
	"context"
	"fmt"
	"sync"
	"testing"
)

var once2 sync.Once
var gRedisCli *Client

func Instance() *Client {
	once2.Do(func() {
		gRedisCli = NewClient(Config{
			RedisAddr:     "localhost",
			RedisPort:     "6379",
			RedisPassword: "",
			DBIndex:       0,
		})
	})

	return gRedisCli
}

// test "rooms roomID1 not exist", which result will get; redis server is not start, which result will get
// both of above situations will return nil and an error object, so only by err cannot judge whether
// "rooms roomID1" is existing or not
func TestClient_HGet(t *testing.T) {
	value, err := Instance().HGet("rooms", "roomID1")
	fmt.Printf("%v %v", value, err)
}

func TestConsumer(t *testing.T) {
	pubSub := Instance().Cli().Subscribe(context.Background(), "redisTopic")

	ch := pubSub.Channel()

	for msg := range ch {
		fmt.Printf("received: %v\n", msg)
	}
}

func TestProducer(t *testing.T) {
	err := Instance().Cli().Publish(context.Background(), "redisTopic", "Hello world").Err()
	if err != nil {
		fmt.Printf("%v\n", err)
	}

	err = Instance().Cli().Close()
	if err != nil {
		fmt.Printf("%v\n ", err)
	}
}
