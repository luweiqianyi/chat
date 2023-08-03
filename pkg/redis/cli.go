package redis

import (
	"chat/pkg/log"
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"sync"
	"time"
)

var gCli *Client
var once sync.Once

type Config struct {
	RedisAddr     string
	RedisPort     string
	RedisPassword string
	DBIndex       int
}

type Client struct {
	cli *redis.Client
}

func NewClient(cfg Config) *Client {
	once.Do(func() {
		gCli = &Client{
			cli: initRedis(cfg.RedisAddr, cfg.RedisPort, cfg.RedisPassword, cfg.DBIndex),
		}
	})
	return gCli
}

func (cli *Client) Cli() *redis.Client {
	return cli.cli
}

func initRedis(addr string, port string, password string, dbIndex int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", addr, port),
		Password: password,
		DB:       dbIndex,
	})
}

func (cli *Client) Set(key string, value interface{}, expiration time.Duration) error {
	if cli.cli == nil {
		return errors.New("redis cli required")
	}

	err := cli.cli.Set(context.Background(), key, value, expiration).Err()
	if err != nil {
		log.Errorf("set {key:%v,value:%v} failed, err: %v", key, value, err)
		return err
	}
	return nil
}

func (cli *Client) Get(key string) (string, error) {
	if cli.cli == nil {
		return "", errors.New("redis cli required")
	}

	value, err := cli.cli.Get(context.Background(), key).Result()
	if err != nil {
		log.Errorf("get key:%v failed, err: %v", key, err)
		return "", err
	}
	return value, nil
}

func (cli *Client) Del(key string) error {
	if cli.cli == nil {
		return errors.New("redis cli required")
	}
	err := cli.cli.Del(context.Background(), key).Err()
	if err != nil {
		log.Errorf("del key:%v failed, err: %v", key, err)
		return err
	}
	return nil
}

func (cli *Client) HSet(key string, field any, value interface{}) error {
	if cli.cli == nil {
		return errors.New("redis cli required")
	}

	ctx := context.Background()
	err := cli.cli.HSet(ctx, key, field, value).Err()
	if err != nil {
		log.Errorf("hSet {key:%v,value:%v} failed, err: %v", key, value, err)
		return err
	}
	return nil
}

func (cli *Client) HSetWithExpiration(key string, field any, value interface{}, expiration time.Duration) error {
	if cli.cli == nil {
		return errors.New("redis cli required")
	}

	ctx := context.Background()
	err := cli.cli.HSet(ctx, key, field, value).Err()
	cli.cli.Expire(ctx, key, expiration)
	if err != nil {
		log.Errorf("hSet {key:%v,value:%v} failed, err: %v", key, value, err)
		return err
	}
	return nil
}

func (cli *Client) HGet(key string, field string) (interface{}, error) {
	if cli.cli == nil {
		return nil, errors.New("redis cli required")
	}

	value, err := cli.cli.HGet(context.Background(), key, field).Result()
	if err != nil {
		log.Errorf("hGet {key:%v field:%v} failed, err: %v", key, field, err)
		return nil, err
	}
	return value, nil
}

func (cli *Client) HDel(key string, field string) error {
	if cli.cli == nil {
		return errors.New("redis cli required")
	}

	err := cli.cli.HDel(context.Background(), key, field).Err()
	if err != nil {
		log.Errorf("hDel {key:%v,field:%v} failed, err: %v", key, field, err)
		return err
	}
	return nil
}

func (cli *Client) HDelAllFields(key string) error {
	if cli.cli == nil {
		return errors.New("redis cli required")
	}

	fields, err := cli.HFields(key)
	if err != nil {
		return err
	}

	for _, field := range fields {
		cli.HDel(key, field)
	}
	return nil
}

func (cli *Client) HExists(key string, field string) bool {
	if cli.cli == nil {
		log.Errorf("redis cli required")
		return false
	}

	exist, err := cli.cli.HExists(context.Background(), key, field).Result()
	if err != nil {
		log.Errorf("hExist {key:%v,field:%v} failed, err: %v", key, field, err)
		return false
	}
	return exist
}

func (cli *Client) HFields(key string) ([]string, error) {
	if cli.cli == nil {
		log.Errorf("redis cli required")
		return nil, errors.New("redis cli required")
	}

	result, err := cli.cli.HKeys(context.Background(), key).Result()
	if err != nil {
		log.Errorf("HKeys get all fields of {key:%v} failed, err: %v", key, err)
		return nil, err
	}
	return result, nil
}

func (cli *Client) HGetAll(key string) (map[string]string, error) {
	if cli.cli == nil {
		log.Errorf("redis cli required")
		return nil, errors.New("redis cli required")
	}

	result, err := cli.cli.HGetAll(context.Background(), key).Result()
	if err != nil {
		log.Errorf("HGetAll field:value of {key:%v} failed, err: %v", key, err)
		return nil, err
	}
	return result, nil
}
