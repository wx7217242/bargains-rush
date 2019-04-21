package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/garyburd/redigo/redis"
	"github.com/wx7217242/bargains-rush/proxy/config"
	"github.com/wx7217242/bargains-rush/proxy/service"
	"github.com/wx7217242/go-common"
	etcd_client "go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"strings"
)

var (
	appConfig = config.AppConfig{
		BargainsRushConfMap: make(map[int]*config.BargainsRushConf, 1024),
	}
	redisPool  *redis.Pool
	etcdClient *etcd_client.Client
)

func initConf() error {

	// 初始化redis配置
	var redisConf common.RedisConf
	redisConf.Addr = beego.AppConfig.String("redis_addr")
	if len(redisConf.Addr) == 0 {
		return fmt.Errorf("redis addr can't be empty");
	}

	if redis_max_idle, err := beego.AppConfig.Int("redis_max_idle"); err != nil {
		redisConf.MaxIdle = 0
	} else {
		redisConf.MaxIdle = redis_max_idle
	}

	if redis_max_active, err := beego.AppConfig.Int("redis_max_active"); err != nil {
		redisConf.MaxActive = 0
	} else {
		redisConf.MaxActive = redis_max_active
	}

	if redis_idle_timeout, err := beego.AppConfig.Int("redis_idle_timeout"); err != nil {
		redisConf.IdleTimeout = 0
	} else {
		redisConf.IdleTimeout = redis_idle_timeout
	}
	appConfig.RedisConf = redisConf

	// 初始化etcd
	key_prefix := beego.AppConfig.String("etcd_bargains_rush_key_prefix")
	product_key := beego.AppConfig.String("etcd_bargains_rush_product_key")

	if len(key_prefix) == 0 || len(product_key) == 0 {
		return fmt.Errorf("etcd_bargains_rush_key_prefix [%v] and etcd_bargains_rush_product_key [%v] cat't be empty", key_prefix, product_key)
	}

	if strings.HasSuffix(key_prefix, "/") == false {
		key_prefix = key_prefix + "/"
	}
	appConfig.ProductKey = fmt.Sprintf("%s%s", key_prefix, product_key)

	var etcdConf common.EtcdConf
	etcdConf.Endpoints = []string{beego.AppConfig.String("etcd_addr")}
	if len(etcdConf.Endpoints) == 0 {
		return fmt.Errorf("etcd addr can't be empty");
	}

	if etcd_timeout, err := beego.AppConfig.Int("etcd_timeout"); err != nil {
		etcdConf.Timeout = 0
	} else {
		etcdConf.Timeout = etcd_timeout
	}
	appConfig.EtcdConf = etcdConf

	appConfig.LogPath = beego.AppConfig.String("log_path")
	appConfig.LogLevel = beego.AppConfig.String("log_level")

	//logs.Debug("appConfig %v", appConfig)

	return nil
}

func initApp() (error) {

	err := common.InitBeegoLogger(appConfig.LogPath, appConfig.LogLevel)
	if err != nil {
		return err
	}

	pool, err := common.InitRedis(appConfig.RedisConf)
	if err != nil {
		return err
	}
	redisPool = pool
	logs.Debug("init redisPool succeed!")

	client, err := common.InitEtcd(appConfig.EtcdConf)
	if err != nil {
		return err
	}
	etcdClient = client
	logs.Debug("init etcdClient succeed!")

	err = appConfig.GetBargainsRushConfFromEtcd(etcdClient)
	if err != nil {
		return err
	}

	service.InitService(appConfig)

	go watchBargainsRushProductKey()

	logs.Info("initApp succeed")

	return nil
}

func watchBargainsRushProductKey() {

	client, err := common.InitEtcd(appConfig.EtcdConf)
	if err != nil {
		logs.Error("init etcd failed, err:%v", err)
		return
	}

	logs.Debug("begin watch key:%s", appConfig.ProductKey)
	for {
		rch := client.Watch(context.Background(), appConfig.ProductKey)
		var secProductInfo []config.BargainsRushConf
		var getConfSucc = true

		for wresp := range rch {
			for _, ev := range wresp.Events {
				if ev.Type == mvccpb.DELETE {
					logs.Warn("key[%s] 's config deleted", appConfig.ProductKey)
					continue
				}

				if ev.Type == mvccpb.PUT && string(ev.Kv.Key) == appConfig.ProductKey {
					err = json.Unmarshal(ev.Kv.Value, &secProductInfo)
					if err != nil {
						logs.Error("key [%s], Unmarshal[%s], err:%v ", err)
						getConfSucc = false
						continue
					}
				}
				logs.Debug("get config from etcd, %s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			}

			if getConfSucc {
				logs.Debug("get config from etcd succ, %v", secProductInfo)
				appConfig.UpdateBargainsRushConf(secProductInfo)
			}
		}

	}
}
