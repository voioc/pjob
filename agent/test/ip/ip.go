/************************************************************
** @Description: ip
** @Author: george hao
** @Date:   2019-06-27 09:22
** @Last Modified by:  george hao
** @Last Modified time: 2019-06-27 09:22
*************************************************************/
package main

import (
	"github.com/astaxie/beego/logs"
	"github.com/voioc/pjob/libs"
)

func main() {
	logs.Info(libs.PublicIp())

}
