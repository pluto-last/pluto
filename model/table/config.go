package table

import (
	"pluto/global"
)

// Config 系统内部需要用到的一些配置
type Config struct {
	global.UUID
	Class string `json:"class" form:"class"` // 配置项分类
	Name  string `json:"name"  form:"name"`  // 配置项名称
	Key   string `json:"key"  form:"key"`
	Value string `json:"value" form:"value"`
	Unit  string `json:"unit" form:"unit"` // 配置项单位  元/个 ...
	Type  string `json:"type" form:"type"` // 配置项的数据类型  string/int/float ...
}

func (Config) TableName() string {
	return "sys_config"
}
