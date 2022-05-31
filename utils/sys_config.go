package utils

import (
	"errors"
	"pluto/global"
	"pluto/model/table"
	"strconv"
	"time"
)

var configMap map[string]table.Config

func init() {
	configMap = make(map[string]table.Config)
	go refreshConfigMap()
}

func GetConfigValueInt64(key string, defaultValue int64) (res int64) {

	if item, ok := configMap[key]; ok {
		res, _ = strconv.ParseInt(item.Value, 10, 64)
		return
	}

	var config table.Config
	if global.GVA_DB.Where("key = ?", key).First(&config).RecordNotFound() {
		return defaultValue
	}
	res, _ = strconv.ParseInt(config.Value, 10, 64)
	return
}

func GetConfigValueInt(key string, defaultValue int) (res int) {

	if item, ok := configMap[key]; ok {
		res, _ = strconv.Atoi(item.Value)
		return
	}

	var config table.Config
	if global.GVA_DB.Where("key = ?", key).First(&config).RecordNotFound() {
		return defaultValue
	}
	res, _ = strconv.Atoi(config.Value)
	return
}

func GetConfigValueFloat64(key string, defaultValue float64) (res float64) {

	if item, ok := configMap[key]; ok {
		res, _ = strconv.ParseFloat(item.Value, 64)
		return
	}

	var config table.Config
	if global.GVA_DB.Where("key = ?", key).First(&config).RecordNotFound() {
		return defaultValue
	}
	res, _ = strconv.ParseFloat(config.Value, 64)
	return
}

func GetConfigValueString(key string, defaultValue string) (res string) {
	if item, ok := configMap[key]; ok {
		return item.Value
	}
	var config table.Config
	if global.GVA_DB.Where("key = ?", key).First(&config).RecordNotFound() {
		return defaultValue
	}
	return config.Value
}

func GetConfigInfo(key string) (table.Config, error) {
	if item, ok := configMap[key]; ok {
		return item, nil
	}

	var config table.Config
	if global.GVA_DB.Where("key = ?", key).First(&config).RecordNotFound() {
		return config, errors.New("获取配置信息失败")
	}
	return config, nil
}

func refreshConfigMap() {

	// 用time.NewTicker 实现定时触发 2m执行一次
	t := time.NewTicker(time.Minute * 2)

	for {
		select {
		case <-t.C:
			var configs []table.Config
			global.GVA_DB.Find(&configs)

			for _, item := range configs {
				configMap[item.Key] = item
			}
		}
	}
}
