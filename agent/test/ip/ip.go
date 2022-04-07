package main

import (
	"github.com/astaxie/beego/logs"
	"github.com/voioc/cjob/utils"
)

func main() {
	logs.Info(utils.PublicIp())

}
