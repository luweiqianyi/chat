package etcd

import (
	"chat/pkg/log"
	"context"
	"fmt"
	uuid2 "github.com/google/uuid"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

const timeFormat = time.RFC3339

type DistributedLocker struct {
	cli         *clientv3.Client
	prefix      string
	leaseID     clientv3.LeaseID
	leaseCancel context.CancelFunc
	isLocked    bool
	KeyName     string // 锁名
	Owner       string // 锁拥有者, 只允许锁拥有者对该锁进行释放, 若锁的拥有者宕机, 锁会自动释放
}

func NewDistributedLocker(cli *Client, prefix string) (*DistributedLocker, error) {
	if cli == nil || cli.cli == nil {
		return nil, fmt.Errorf("create locker failed, etcd connected failed")
	}

	return &DistributedLocker{
		cli:    cli.cli,
		prefix: prefix,
	}, nil
}

func (locker *DistributedLocker) GetPrefix() string {
	return locker.prefix
}

func (locker *DistributedLocker) GetLeaseID() clientv3.LeaseID {
	return locker.leaseID
}

func (locker *DistributedLocker) SetLeaseID(id clientv3.LeaseID) {
	locker.leaseID = id
}

func (locker *DistributedLocker) IsLocked() bool {
	return locker.isLocked
}

func (locker *DistributedLocker) SetLocked(locked bool) {
	locker.isLocked = locked
}

func (locker *DistributedLocker) GetKeyName() string {
	return locker.KeyName
}

func (locker *DistributedLocker) SetKeyName(key string) {
	locker.KeyName = key
}

func (locker *DistributedLocker) GetOwner() string {
	return locker.Owner
}

func (locker *DistributedLocker) SetOwner(owner string) {
	locker.Owner = owner
}

// AcquireLock 获取分布式锁，内置锁获取次数(固定), 固定次数还没获取到锁返回失败
func (locker *DistributedLocker) AcquireLock(key string, ttl int64) error {
	if locker.IsLocked() {
		return fmt.Errorf("locker %v already locked", locker)
	}

	var ctx context.Context
	var cancel context.CancelFunc

	retryTimes := 3
	// 增加重试次数，防止因为etcd集群繁忙或者网络不佳等原因导致的请求失败而直接返回错误
	for i := 0; i < retryTimes; i++ {
		// 生成租约
		lease, err := locker.cli.Grant(context.Background(), ttl)
		if err != nil {
			continue
		}

		// 生成ctx对象，用来控制租约的生命周期
		ctx, cancel = context.WithCancel(context.Background())

		// 向etcd服务端发起心跳, 对租约进行续租
		leaseAckChan, err := locker.cli.KeepAlive(ctx, lease.ID)
		if err != nil {
			cancel()                                          // 取消续租, 会触发ctx.Done()
			locker.cli.Revoke(context.Background(), lease.ID) // 释放etcd上的租约资源
			return err
		}

		lc := locker.GetPrefix() + "-" + key
		uuid := uuid2.NewString()
		fmt.Printf("%v: try get locker{%v %v}...\n",
			time.Now().Format(timeFormat), lc, uuid)

		//// 等待该key被释放或过期,加入了backoff算法，限制了重试的间隔时间，避免了连续调用对etcd集群造成过大压力的问题。这段代码会无限阻塞当前协程
		// 造成协程资源无限占用，是不合适的
		//for {
		//	getResp, err := locker.cli.Get(context.Background(), KeyName)
		//	if err != nil {
		//		time.Sleep(300 * time.Millisecond)
		//		continue
		//	}
		//	if getResp.Count == 0 {
		//		break
		//	}
		//	time.Sleep(300 * time.Millisecond)
		//}

		// 利用etcd事务的原子性，往etcd中写入锁信息。写入前先判断锁有没有被其他占有
		resp, err := locker.cli.Txn(context.Background()).
			If(clientv3.Compare(clientv3.CreateRevision(lc), "=", 0)).
			Then(clientv3.OpPut(lc, uuid, clientv3.WithLease(lease.ID))).
			Commit()

		// 写入锁失败, 表示获取分布式锁失败
		if err != nil {
			cancel()
			locker.cli.Revoke(context.Background(), lease.ID) // 撤销上一次的租约
			fmt.Printf("try get locker {%v %v} %d times failed, err: %v\n", key, uuid, i+1, err)
			time.Sleep(300 * time.Millisecond) // 固定时间退避算法，避免对etcd造成比较大的访问压力
			continue                           // 不直接return，continue继续尝试获取
		}
		if !resp.Succeeded {
			cancel()
			locker.cli.Revoke(context.Background(), lease.ID)
			fmt.Printf("try get locker {%v %v} %d times failed, commitResp: %v\n", key, uuid, i+1, resp)
			time.Sleep(300 * time.Millisecond) // 固定时间退避算法，避免对etcd造成比较大的访问压力
			continue                           // 不直接return，continue继续尝试获取
		}

		log.Infof("distribute locker {%v %v} create success", lc, uuid)
		// 创建分布式锁成功
		// 关联租约
		locker.SetLeaseID(lease.ID)
		locker.leaseCancel = cancel
		locker.SetLocked(true)
		locker.SetKeyName(lc)
		locker.SetOwner(uuid)

		go func(ctx context.Context, leaseID clientv3.LeaseID, lc string, uuid string) {
			defer cancel() // 调用cancel取消对leaseID的lease的续租
			for {
				select {
				case <-ctx.Done(): // ctx被取消
					log.Infof("leaseID[%v] keepalive self-stopped!", leaseID)
					return
				case _, ok := <-leaseAckChan: // 租约到期, ok=false
					if !ok {
						log.Infof("leaseID[%v] keepalive detected stopped!", leaseID)
						return
					} else {
						fmt.Printf("%v: locker{%v %v} refresh success\n",
							time.Now().Format(timeFormat), lc, uuid)
					}
				}

				//// 确保锁与租约关联仍存在(TODO 这个查询是不是会加大对于etcd的访问压力??)
				//getResp, _ := locker.cli.Get(context.Background(), lc)
				//if getResp.Count == 0 || string(getResp.Kvs[0].Value) != uuid {
				//	return
				//}
				fmt.Printf("%v: run forloop\n", time.Now().Format(timeFormat))
			}
		}(ctx, lease.ID, lc, uuid)

		if locker.IsLocked() {
			fmt.Printf("%v: locker{%v %v} get success\n",
				time.Now().Format(timeFormat), lc, uuid)
			return nil
		} // else try next time
	}
	return fmt.Errorf("failed to acquire lock after retries, key: %s, locker: %v, retrytimes: %v",
		key, locker, retryTimes)
}

// ReleaseLock 通过删除locker.prefix+key键来达到释放分布式锁的目的
func (locker *DistributedLocker) ReleaseLock() error {
	if !locker.IsLocked() {
		return fmt.Errorf("locker %v already released", locker)
	}

	// 如果当前客户端宕机也不会影响锁的释放, 因为当前客户端宕机后就无法向etcd服务端发送租约心跳, 这必然会导致锁在租约过期后被主动释放

	curLockerKey := locker.GetKeyName()
	// 释放前先检查锁是否已经释放或者被其他客户端占有，这两种情况不允许释放锁
	resp, err := clientv3.NewKV(locker.cli).Get(context.Background(), curLockerKey)
	if err != nil {
		fmt.Printf("locker{%v} release failed,not locked by self,err: %v\n", curLockerKey, err)
		return fmt.Errorf("locker{%v} release failed,not locked by self,err: %v", curLockerKey, err)
	}
	if len(resp.Kvs) == 0 {
		fmt.Printf("locker{%v} release failed,already released\n", curLockerKey)
		return fmt.Errorf("locker{%v} release failed,already released", curLockerKey)
	}
	owner := string(resp.Kvs[0].Value)  // curLockerKey可能被多个单位持有，需要获取目前持有者的名字
	curLockerOwner := locker.GetOwner() // this持有者
	if owner != curLockerOwner {
		fmt.Printf("locker{%v} release failed: lock held by another client", curLockerKey)
		return fmt.Errorf("locker{%v} release failed: lock held by another client", curLockerKey)
	}

	// 取消锁的租约，达到释放锁的目的
	locker.SetLocked(false)
	locker.leaseCancel()                                                  // 取消租约，锁会被etcd自动释放
	_, err = locker.cli.Revoke(context.Background(), locker.GetLeaseID()) // 撤销etcd中租约,避免资源泄露
	locker.SetLeaseID(0)

	if err == nil {
		fmt.Printf("locker{%v %v} release success\n", curLockerKey, curLockerOwner)
	} else {
		fmt.Printf("locker{%v %v} release failed,err: %v\n", curLockerKey, curLockerOwner, err)
	}

	return err
}

func DoTaskWithDistributeLock(cli *Client, key string, ttl int64, task func() error) error {
	if cli == nil {
		return fmt.Errorf("task executed failed: cli ni")
	}
	locker, err := NewDistributedLocker(cli, "task")
	if err != nil {
		return fmt.Errorf("task executed failed: locker create failed, err: %v", err)
	}

	err = locker.AcquireLock(key, ttl)
	if err != nil {
		return fmt.Errorf("task executed failed: locker create failed, err: %v", err)
	}

	defer func() {
		err := locker.ReleaseLock()
		if err != nil {
			fmt.Printf("task executed failed: locker release failed, err: %v\n", err)
		}
	}()

	err = task()
	if err != nil {
		fmt.Printf("task executed failed,err: %v", err)
	}

	return nil
}
