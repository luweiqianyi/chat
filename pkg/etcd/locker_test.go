package etcd

import (
	"fmt"
	"testing"
	"time"
)

func login() error {
	fmt.Printf("%v: login start...\n", time.Now().Format(time.RFC3339))
	time.Sleep(time.Second * 30)
	fmt.Printf("%v: login end.\n", time.Now().Format(time.RFC3339))
	return nil
}

func TestEtcdDistributedLock(t *testing.T) {
	cli := NewClient(NewConfig([]string{"localhost:2379"}, time.Second*10))
	if cli == nil {
		fmt.Printf("task executed failed: connected to etcd failed")
		return
	}
	defer cli.Close()

	// 业务login, 业务执行时间30s左右, 分布锁过期时间10s
	err := DoTaskWithDistributeLock(cli, "login", 10, login)
	if err != nil {
		fmt.Printf("%v \n", err)
	}
}
