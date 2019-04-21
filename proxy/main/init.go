package main

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/garyburd/redigo/redis"
	"github.com/wx7217242/go-common"
	etcd_client "go.etcd.io/etcd/clientv3"
	"strings"
)

type AppConfig struct {
	redisConf  util.RedisConf
	etcdConf   util.EtcdConf
	productKey string
	logPath    string
	logLevel   string
}

var (
	appConfig  AppConfig
	redisPool  *redis.Pool
	etcdClient *etcd_client.Client
)

func initConf() error {

	// 初始化redis配置
	var redisConf util.RedisConf
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
	appConfig.redisConf = redisConf

	// 初始化etcd
	key_prefix := beego.AppConfig.String("etcd_bargains_rush_key_prefix")
	product_key := beego.AppConfig.String("etcd_bargains_rush_product_key")

	if len(key_prefix) == 0 || len(product_key) == 0 {
		return fmt.Errorf("etcd_bargains_rush_key_prefix [%v] and etcd_bargains_rush_product_key [%v] cat't be empty", key_prefix, product_key)
	}

	if strings.HasSuffix(key_prefix, "/") == false {
		key_prefix = key_prefix + "/"
	}
	appConfig.productKey = fmt.Sprintf("%s%s", key_prefix, product_key)

	var etcdConf util.EtcdConf
	etcdConf.Endpoints = []string{beego.AppConfig.String("etcd_addr")}
	if len(etcdConf.Endpoints) == 0 {
		return fmt.Errorf("etcd addr can't be empty");
	}

	if etcd_timeout, err := beego.AppConfig.Int("etcd_timeout"); err != nil {
		etcdConf.Timeout = 0
	} else {
		etcdConf.Timeout = etcd_timeout
	}
	appConfig.etcdConf = etcdConf

	appConfig.logPath = beego.AppConfig.String("log_path")
	appConfig.logLevel = beego.AppConfig.String("log_level")

	//logs.Debug("appConfig %v", appConfig)

	return nil
}

func initApp() (error) {

	err := util.InitLogger(appConfig.logPath, appConfig.logLevel)
	if err != nil {
		return err
	}

	pool, err := util.InitRedis(appConfig.redisConf)
	if err != nil {
		return err
	}
	redisPool = pool
	logs.Debug("init redisPool succeed!")

	client, err := util.InitEtcd(appConfig.etcdConf)
	if err != nil {
		return err
	}
	etcdClient = client
	logs.Debug("init etcdClient succeed!")

	return nil
}
