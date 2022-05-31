package config

import (
	"fmt"
	"log"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var once sync.Once
var configInfo = new(Specification)

type Specification struct {
	Port int
	DB   struct {
		Typ   string `yaml:"type"`
		DSN   string `yaml:"dsn"`
		Debug bool
	}
	Cache struct {
		Addr     string `yaml:"addr"`
		Password string `yaml:"password"`
	}
	JWT struct {
		Timeout int64  `yaml:"timeout"`
		SignKey string `yaml:"signKey"`
	}
	SMS struct {
		APIKey   string `yaml:"apiKey"`
		Template string `yaml:"template"`
	}
	Zap struct {
		Level        string `yaml:"level"`          // 级别
		Format       string `yaml:"format"`         // 输出
		Prefix       string `yaml:"prefix"`         // 日志前缀
		Director     string `yaml:"director"`       // 日志文件夹
		LogInConsole bool   `yaml:"log-in-console"` // 输出控制台
	}
	ExcelDir string `yaml:"excel-dir"` // excel的导出路径
	RSA      struct {
		Account      string `yaml:"account"`
		Password     string `yaml:"password"`
		RobotListUrl string `yaml:"robotListUrl"`
	} `yaml:"rsa"`
	WsPath string `yaml:"wsPath"`
}

func Get(fileName string) (*Specification, error) {
	once.Do(func() {
		if fileName == "" {
			fileName = "./conf/conf.yaml"
		}
		err := Refresh(fileName)
		if err != nil {
			log.Fatalln("read conf file: ", err)
		}
	})
	return configInfo, nil
}

func Refresh(fileName string) error {

	viper.SetConfigFile(fileName) // 指定配置文件路径
	err := viper.ReadInConfig()   // 读取配置信息
	if err != nil {               // 读取配置信息失败
		log.Fatalln("read conf file: ", err)
	}
	// 监控配置文件变化
	viper.WatchConfig()

	// 将读取的配置信息保存至全局变量Conf
	if err = viper.Unmarshal(configInfo); err != nil {
		log.Fatalln("read conf file: ", err)
	}

	// 注意！！！配置文件发生变化后要同步到全局变量Conf
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("夭寿啦~配置文件被人修改啦...")
		viper.Unmarshal(configInfo)
	})

	return nil
}
