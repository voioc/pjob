package init

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

var env string = "dev"

func init() {

	// 初始化数据模型
	// model.Init(StartTime)
	// model.Init(time.Now().Unix())
	// jobs.InitJobs()

	viper.SetDefault("envType", env)

	// 获取当前环境变量
	if realEnv := strings.ToLower(os.Getenv("envType")); realEnv != "" {
		env = realEnv
		viper.SetDefault("envType", env)
	}

	path, _ := filepath.Abs(filepath.Dir(""))        // 获取当前路径
	conf := path + "/config/config_" + env + ".toml" // 拼接配置文件
	fmt.Println(conf)

	configFile := flag.String("c", conf, "配置文件路径") // 手动置顶配置文件
	flag.Parse()

	viper.SetConfigFile(*configFile) // 读取配置文件
	fmt.Println("Loading config file " + *configFile)

	if err := viper.ReadInConfig(); err != nil { //是否读取成功
		log.Fatalln("打开配置文件失败", err)
	}
}
