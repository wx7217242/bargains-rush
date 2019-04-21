package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/wx7217242/bargains-rush/proxy/config"
	"github.com/wx7217242/go-common"
	"log"
	"time"
)

var (
	EtcdKey = "/bargains-rush/product"
)

func main() {

	etcdConf := util.EtcdConf{
		Endpoints: []string{"127.0.0.1:2379"},
		Timeout:   5,
	}

	etcdClient, err := util.InitEtcd(etcdConf)
	if err != nil {
		panic(err)
	}
	defer etcdClient.Close()
	log.Println("connect etcd succeed")

	var SecInfoConfArr []config.BargainsRushConf
	SecInfoConfArr = append(
		SecInfoConfArr,
		config.BargainsRushConf{
			ProductId: 1029,
			StartTime: 1505008800,
			EndTime:   1505012400,
			Status:    0,
			Total:     1000,
			Left:      1000,
		},
	)
	SecInfoConfArr = append(
		SecInfoConfArr,
		config.BargainsRushConf{
			ProductId: 1027,
			StartTime: 1505008800,
			EndTime:   1505012400,
			Status:    0,
			Total:     2000,
			Left:      1000,
		},
	)

	data, err := json.Marshal(SecInfoConfArr)
	if err != nil {
		fmt.Println("json failed, ", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	//cli.Delete(ctx, EtcdKey)
	//return
	_, err = etcdClient.Put(ctx, EtcdKey, string(data))
	cancel()
	if err != nil {
		fmt.Println("put failed, err:", err)
		return
	}

	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	resp, err := etcdClient.Get(ctx, EtcdKey)
	cancel()
	if err != nil {
		fmt.Println("get failed, err:", err)
		return
	}
	for _, ev := range resp.Kvs {
		fmt.Printf("%s : %s\n", ev.Key, ev.Value)
	}
}
