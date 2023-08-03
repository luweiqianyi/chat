package etcd

import (
	"chat/pkg/log"
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

// Person 随便定义一个数据结构, 模拟实际业务场景中的value值写入etcd数据库中
type Person struct {
	Name    string
	Age     int
	Address string
}

/* 对封装库的单元测试 */
func TestEtcdPutK1(t *testing.T) {
	cli := NewClient(NewConfig([]string{"localhost:2379"}, time.Second*30))
	defer cli.Close()

	k1 := "k1"
	p1 := Person{
		Name:    "zhangSan",
		Age:     23,
		Address: "HangZhou",
	}
	jsonBody, _ := json.Marshal(p1)
	v1 := string(jsonBody)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if err := cli.Put(ctx, k1, v1); err != nil {
		fmt.Printf("%v\n", err)
	}

	log.Infof("end")
}

func TestEtcdDelK1(t *testing.T) {
	cli := NewClient(NewConfig([]string{"localhost:2379"}, time.Second*30))
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	key := "k1"
	cli.Del(ctx, key)
}

func TestEtcdPutK2(t *testing.T) {
	cli := NewClient(NewConfig([]string{"localhost:2379"}, time.Second*30))
	defer cli.Close()

	k2 := "k2"
	p2 := Person{
		Name:    "LeeSi",
		Age:     28,
		Address: "ShangHai",
	}
	jsonBody, _ := json.Marshal(p2)
	v2 := string(jsonBody)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if err := cli.Put(ctx, k2, v2); err != nil {
		fmt.Printf("%v\n", err)
	}

	log.Infof("end")
}

func TestEtcdGetK1(t *testing.T) {
	cli := NewClient(NewConfig([]string{"localhost:2379"}, time.Second*30))
	defer cli.Close()

	k1 := "k1"

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	r1, _ := cli.Get(ctx, k1)
	fmt.Printf("%v\n", r1)

	pp1 := &Person{}
	json.Unmarshal([]byte(r1), pp1)
	fmt.Println(pp1)

	log.Infof("end")
}

func TestEtcdGetK2(t *testing.T) {
	cli := NewClient(NewConfig([]string{"localhost:2379"}, time.Second*30))
	defer cli.Close()

	k2 := "k2"

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	r2, _ := cli.Get(ctx, k2)
	fmt.Printf("%v\n", r2)

	pp2 := &Person{}
	json.Unmarshal([]byte(r2), pp2)
	fmt.Println(pp2)

	log.Infof("end")
}

// TestPutWithLease 测试一个key被授予一个lease时，当lease过期key自动被删除的情况
func TestPutWithLease(t *testing.T) {
	cli := NewClient(NewConfig([]string{"localhost:2379"}, time.Second*30))
	defer cli.Close()

	err := cli.PutWithLease("my-key", "Etcd Lease Auto Release", 10)
	if err != nil {
		fmt.Println(err)
	}
}

// TestEtcdKeepAlive 为key进行续约
func TestEtcdKeepAlive(t *testing.T) {
	endPoints := []string{"localhost:2379"}
	cfg := NewConfig(endPoints, time.Second*30)
	cli := NewClient(cfg)
	defer cli.Close()

	err := cli.StartKeepAlive("my-key", "keepalive", 10)
	if err != nil {
		fmt.Printf("TestEtcdKeepAlive: %v", err)
	}
}

func TestEtcdWatch(t *testing.T) {
	cli := NewClient(NewConfig([]string{"localhost:2379"}, time.Second*30))
	defer cli.Close()

	ctx := context.Background()
	err := cli.Watch(ctx, "k1", etcdWatchCallBack) // 观察k1的变化，通过交替运行TestEtcdPutK1、TestEtcdDelK1来观察程序执行现象
	if err != nil {
		fmt.Printf("TestEtcdWatch: %v\n", err)
	}
}

func etcdWatchCallBack(eventType string, info KV) {
	switch eventType {
	case Put:
		fmt.Printf("put event, info:%v\n", info)
	case Delete:
		fmt.Printf("delete event, info: %v\n", info)
	}
}

// 测试租约被撤销时，与其绑定的k-v键值对是否会被etcd自动删除,经过测试: 是自动删除的
func TestMyTest(t *testing.T) {
	cli := NewClient(NewConfig([]string{"localhost:2379"}, time.Second*30))
	defer cli.Close()

	cli.RevokeTest("Leebai", "loves drinking", 30)
}
