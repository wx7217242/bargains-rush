package config

import (
	"context"
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"github.com/wx7217242/go-common"
	etcd_client "go.etcd.io/etcd/clientv3"
	"sync"
)

type AppConfig struct {
	RedisConf  util.RedisConf
	EtcdConf   util.EtcdConf
	ProductKey string
	LogPath    string
	LogLevel   string

	BargainsRushConfMap map[int]*BargainsRushConf
	RWProductLock       sync.RWMutex
}

type BargainsRushConf struct {
	ProductId int
	StartTime int
	EndTime   int
	Status    int
	Total     int
	Left      int
}

func (config *AppConfig) UpdateBargainsRushConf(conf []BargainsRushConf) {
	var tmp map[int]*BargainsRushConf = make(map[int]*BargainsRushConf, 1024)
	for _, v := range conf {
		tmp[v.ProductId] = &v
	}

	config.RWProductLock.Lock()
	config.BargainsRushConfMap = tmp
	config.RWProductLock.Unlock()
}

func (config *AppConfig) GetBargainsRushConfFromEtcd(etcdClient *etcd_client.Client) error {
	response, err := etcdClient.Get(context.Background(), config.ProductKey)
	if err != nil {
		logs.Error("get [%s] from etcd failed, err:%v", config.ProductKey, err)
		return err
	}

	var conf []BargainsRushConf
	for k, v := range response.Kvs {

		logs.Debug("key [%s] value [%s]", k, v)
		err = json.Unmarshal(v.Value, &conf)
		if err != nil {
			logs.Error("Unmarshal bargains rush info failed, err:%v", err)
			continue
		}

		logs.Debug("bargains rush info is [%v]", conf)
	}

	config.UpdateBargainsRushConf(conf)
	return nil
}
