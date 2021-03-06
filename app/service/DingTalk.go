package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
	"github.com/voioc/cjob/utils"
	"github.com/voioc/coco/logzap"
)

type Dingtalk struct {
	Dingtalks map[string]string
	Content   map[string]interface{}
}

var DingtalkChan chan *Dingtalk
var DingtalkUrl string

func init() {
	DingtalkUrl = viper.GetString("ding.url")
	poolSize := 10 // viper.GetInt("dingtalk.pool")

	//创建通道
	DingtalkChan = make(chan *Dingtalk, poolSize)

	go func() {
		for {
			select {
			case m, ok := <-DingtalkChan:
				if !ok {
					return
				}
				if err := m.SendDingtalk(); err != nil {
					logzap.Ex(context.Background(), "SendDingtalk", err.Error())
					// beego.Error("SendDingtalk:", err.Error())

				}
			}
		}
	}()

}

func SendDingtalkToChan(dingtalks map[string]string, content map[string]interface{}) bool {
	dingTalk := &Dingtalk{
		Dingtalks: dingtalks,
		Content:   content,
	}

	select {
	case DingtalkChan <- dingTalk:
		return true
	case <-time.After(time.Second * 3):
		return false
	}
}

func (s *Dingtalk) SendDingtalk() error {

	for _, v := range s.Dingtalks {
		body, err := json.Marshal(s.Content)
		if err != nil {
			log.Println(err)
			return err
		}

		url := fmt.Sprintf(DingtalkUrl, v)
		_, resErr := utils.HttpPost(url, "application/json;charset=utf-8", bytes.NewBuffer(body))
		if resErr != nil {
			log.Println(resErr)
			return resErr
		}
		return nil
	}
	return nil
}
