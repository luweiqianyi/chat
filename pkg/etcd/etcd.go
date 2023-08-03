package etcd

import (
	"chat/pkg/log"
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"sync"
	"time"
)

const (
	Put    = "Put"
	Delete = "Delete"
)

// KV 自定义etcd k-v键值对, Value: json数据
type KV struct {
	Key   string
	Value string
}

type WatchChangeCallback func(eventType string, info KV) // 传接口类型,增强程序可扩展性

type Client struct {
	cli *clientv3.Client
}

type Config struct {
	cfg clientv3.Config
}

func NewConfig(endPoints []string, dialTimeout time.Duration) Config {
	return Config{
		cfg: clientv3.Config{
			Endpoints:   endPoints,
			DialTimeout: dialTimeout,
		},
	}
}

func NewClient(cfg Config) *Client {
	cli, err := clientv3.New(cfg.cfg)
	if err != nil {
		log.Errorf("etcd connect failed,err:%v", err)
		return nil
	}
	return &Client{
		cli: cli,
	}
}

func (cli *Client) Close() error {
	if cli == nil || cli.cli == nil {
		return fmt.Errorf("ectd cli nil")

	}
	return cli.cli.Close()
}

func (cli *Client) Put(ctx context.Context, key string, value string, opts ...clientv3.OpOption) error {
	if cli.cli == nil {
		return fmt.Errorf("ectd cli nil")
	}
	_, err := cli.cli.Put(ctx, key, value, opts...)
	if err != nil {
		return fmt.Errorf("etcd put key:Owner{%v:%v} failed,err:%v", key, value, err)
	}
	return nil
}

func (cli *Client) Get(ctx context.Context, key string, opts ...clientv3.OpOption) (string, error) {
	if cli.cli == nil {
		return "", fmt.Errorf("ectd cli nil")
	}

	resp, err := cli.cli.Get(ctx, key, opts...)
	if err != nil {
		return "", err
	}

	for _, ev := range resp.Kvs {
		return fmt.Sprintf("%s", ev.Value), nil
	}
	return "", nil
}

func (cli *Client) Del(ctx context.Context, key string, opts ...clientv3.OpOption) error {
	if cli.cli == nil {
		return fmt.Errorf("ectd cli nil")
	}

	_, err := cli.cli.Delete(ctx, key, opts...)
	if err != nil {
		return err
	}
	return nil
}

func (cli *Client) CreateLease(ctx context.Context, ttl int64) (int64, error) {
	if cli.cli == nil {
		return 0, fmt.Errorf("ectd cli nil")
	}

	resp, err := cli.cli.Grant(ctx, ttl)
	if err != nil {
		return 0, err
	}
	return int64(resp.ID), nil
}

// PutWithLease 创建一个过期时间为leaseTTL的租约，和etcd中的{key,Owner}条项进行绑定，租约到期，该条项会被自动删除。leaseTTL单位: 秒
func (cli *Client) PutWithLease(key string, value string, leaseTTL int64) error {
	if cli.cli == nil {
		return fmt.Errorf("ectd cli nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	leaseID, err := cli.CreateLease(ctx, leaseTTL)
	if err != nil {
		return err
	}

	_, err = cli.cli.Put(ctx, key, value, clientv3.WithLease(clientv3.LeaseID(leaseID)))
	if err != nil {
		return fmt.Errorf("put[%v]{%v} with lease[%v] failed,err: %v", key, value, leaseID, err)
	}
	return nil
}

func (cli *Client) StartKeepAlive(key string, value string, leaseTTL int64) error {
	if cli.cli == nil {
		return fmt.Errorf("ectd cli nil")
	}

	// 测试： 创建一个context对象，将其传给其创建的子协程可以控制该子协程的声明周期
	//ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	//defer cancel()
	// be caution: 这里肯定不能用一个timeout的context，在Revoke时会发生context deadline exceeded错误，导致资源不能被正确释放

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	resp, err := cli.cli.Grant(ctx, leaseTTL)
	if err != nil {
		return fmt.Errorf("StartKeepAlive 1: %v", err)
	}
	leaseID := resp.ID

	_, err = cli.cli.Put(ctx, key, value, clientv3.WithLease(leaseID))
	if err != nil {
		return fmt.Errorf("put[%v]{%v} with lease[%v] failed,err: %v", key, value, leaseID, err)
	}

	if err != nil {
		return fmt.Errorf("StartKeepAlive 2: %v", err)
	}

	// 向etcd服务端发起发送心跳，进行lease的租约续期
	leaseKeepAliveResponseCh, err := cli.cli.KeepAlive(ctx, leaseID)
	if err != nil {
		return fmt.Errorf("StartKeepAlive 3: %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func(ctx context.Context) {
		defer wg.Done()
		defer func() {
			if p := recover(); p != nil {
				log.Panicf("%v", p)
			}
		}()

		existFlag := false
		for !existFlag {
			select {
			case <-ctx.Done(): // cancel()函数会触发该case
				existFlag = true
				fmt.Printf("ctx.Done(), exist \n")
			case leaseKeepAliveResponse, _ := <-leaseKeepAliveResponseCh:
				fmt.Printf("%v: %v\n", time.Now().Format("2006-01-02 15:04:05.000"), leaseKeepAliveResponse)
				if leaseKeepAliveResponse == nil {
					fmt.Printf("leaseKeepAliveResponse nil, exist \n")
					existFlag = true
				} else {
					log.Infof("leaseID:%v, ttl:%v", leaseID, leaseKeepAliveResponse.TTL)
					fmt.Printf("leaseID:%v, ttl:%v\n", leaseID, leaseKeepAliveResponse.TTL)
				}
			}
		}
	}(ctx)
	wg.Wait() // 阻塞当前协程，让其等待子协程运行结束

	// 撤销租约，避免资源浪费
	_, err = cli.cli.Revoke(ctx, leaseID)
	fmt.Printf("revoke")
	if err != nil {
		return fmt.Errorf("StartKeepAlive 5: %v", err)
	}
	return nil
}

// Watch 阻塞当前协程一直监听key上面的变化，指定回调函数来进行处理
func (cli *Client) Watch(ctx context.Context, key string, callback WatchChangeCallback, opts ...clientv3.OpOption) error {
	if cli.cli == nil {
		return fmt.Errorf("ectd cli nil")
	}

	fmt.Printf("watch key[%s]'s change\n", key)
	log.Infof("watch key[%s]'s change", key)
	watchChan := cli.cli.Watch(ctx, key, opts...)
	for watchResponse := range watchChan {
		fmt.Printf("received key[%s]'s change response:%v\n", key, watchResponse)
		log.Infof("received key[%s]'s change response:%v", key, watchResponse)
		for _, ev := range watchResponse.Events {
			fmt.Printf("Key %q Value %q\n", ev.Kv.Key, ev.Kv.Value)
			log.Infof("key[%s]'s change response,k-v: %v %v", key, ev.Kv.Key, ev.Kv.Value)

			info := KV{
				Key:   fmt.Sprintf("%s", ev.Kv.Key),
				Value: fmt.Sprintf("%s", ev.Kv.Value),
			}
			switch ev.Type {
			case clientv3.EventTypePut:
				callback(Put, info)
			case clientv3.EventTypeDelete:
				callback(Delete, info)
			}
		}
	}
	fmt.Printf("end watch[%s]\n", key)
	log.Infof("end watch[%s]", key)
	return nil
}

// RevokeTest leaseTTL传入一个大于30s的值 模拟在程序运行的第10s后，对租约进行Revoke时，观察上面在etcd中创建的键值对是否会被自动删除
func (cli *Client) RevokeTest(key string, value string, leaseTTL int64) error {
	if cli.cli == nil {
		return fmt.Errorf("ectd cli nil")
	}

	ch := make(chan int64)
	go func(key string, value string, leaseTTL int64) error {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		resp, err := cli.cli.Grant(ctx, leaseTTL)
		if err != nil {
			return fmt.Errorf("StartKeepAlive 1: %v", err)
		}
		leaseID := resp.ID

		_, err = cli.cli.Put(ctx, key, value, clientv3.WithLease(leaseID))
		if err != nil {
			return fmt.Errorf("put[%v]{%v} with lease[%v] failed,err: %v", key, value, leaseID, err)
		}

		if err != nil {
			return fmt.Errorf("StartKeepAlive 2: %v", err)
		}

		// 向etcd服务端发起发送心跳，进行lease的租约续期
		leaseKeepAliveResponseCh, err := cli.cli.KeepAlive(ctx, leaseID)
		if err != nil {
			return fmt.Errorf("StartKeepAlive 3: %v", err)
		}

		var wg sync.WaitGroup
		wg.Add(1)
		go func(ctx context.Context) {
			defer wg.Done()
			defer func() {
				if p := recover(); p != nil {
					log.Panicf("%v", p)
				}
			}()

			existFlag := false
			for !existFlag {
				select {
				case <-ctx.Done(): // cancel()函数会触发该case
					existFlag = true
					fmt.Printf("ctx.Done(), exist \n")
				case leaseKeepAliveResponse, _ := <-leaseKeepAliveResponseCh:
					fmt.Printf("%v: %v\n", time.Now().Format("2006-01-02 15:04:05.000"), leaseKeepAliveResponse)
					if leaseKeepAliveResponse == nil {
						fmt.Printf("leaseKeepAliveResponse nil, exist \n")
						existFlag = true
					} else {
						log.Infof("leaseID:%v, ttl:%v", leaseID, leaseKeepAliveResponse.TTL)
						fmt.Printf("leaseID:%v, ttl:%v\n", leaseID, leaseKeepAliveResponse.TTL)
					}
				}
			}
		}(ctx)

		ch <- int64(leaseID) // 告诉另外一个协程租约ID
		wg.Wait()            // 阻塞当前协程，让其等待子协程运行结束

		// 撤销租约，避免资源浪费
		_, err = cli.cli.Revoke(ctx, leaseID)
		fmt.Printf("revoke")
		if err != nil {
			return fmt.Errorf("StartKeepAlive 5: %v", err)
		}

		return nil
	}(key, value, leaseTTL)

	// 以上代码是对某个键值对进行不断续租的操作
	// 以下代码是在模拟在程序运行的第10s后，对租约进行Revoke时，观察上面在etcd中创建的
	// 键值对是否会被自动删除
	time.Sleep(time.Second * 10)
	leaseID := <-ch // 阻塞直到获取租约ID
	_, err := cli.cli.Revoke(context.Background(), clientv3.LeaseID(leaseID))
	fmt.Printf("revoke\n")
	if err != nil {
		return fmt.Errorf("StartKeepAlive 5: %v", err)
	}

	return nil
}
