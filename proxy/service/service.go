package service

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/wx7217242/bargains-rush/proxy/config"
)

var (
	appConfig config.AppConfig
)

func InitService(config config.AppConfig) {
	appConfig = config

	logs.Debug("init service succeed")
}

func Info(productId int) (data map[string]interface{}, code int, err error) {

	appConfig.RWProductLock.RLock()
	defer appConfig.RWProductLock.RUnlock()

	v, ok := appConfig.BargainsRushConfMap[productId]
	if !ok {
		code = ErrInvalidProductId
		err = fmt.Errorf("not found product_id:%d", productId)
		return
	}

	data = make(map[string]interface{})
	data["product_id"] = productId
	data["start_time"] = v.StartTime
	data["end_time"] = v.EndTime
	data["status"] = v.Status

	return
}
