package lalilog

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/aliyun/aliyun-log-go-sdk/producer"
	"github.com/astaxie/beego/logs"
	lconv "github.com/lixu-any/go-tools/conv"
	ltime "github.com/lixu-any/go-tools/time"
)

var (
	AliProducer      *producer.Producer
	AliDefaultConfig MapAlilogConfig //默认配置
)

type MapAlilogConfig struct {
	Endpoint        string
	AccessKeyId     string
	AccessKeySecret string
	Project         string
	On              string
	Debug           string
}

//初始化阿里云日志
func InitAliLog(config MapAlilogConfig) {

	if config.On != "1" {
		return
	}

	producerConfig := producer.GetDefaultProducerConfig()

	producerConfig.Endpoint = config.Endpoint
	producerConfig.AccessKeyID = config.AccessKeyId
	producerConfig.AccessKeySecret = config.AccessKeySecret

	AliProducer = producer.InitProducer(producerConfig)

	ch := make(chan os.Signal)

	signal.Notify(ch)

	AliProducer.Start()

	//AliProducer.SafeClose() // 安全关闭
}

func formatLog(f interface{}, v ...interface{}) string {
	var msg string
	switch f := f.(type) {
	case string:
		msg = f
		if len(v) == 0 {
			return msg
		}
		if strings.Contains(msg, "%") && !strings.Contains(msg, "%%") {
			//format string
		} else {
			//do not contain format char
			msg += strings.Repeat(" %v", len(v))
		}
	default:
		msg = fmt.Sprint(f)
		if len(v) == 0 {
			return msg
		}
		msg += strings.Repeat(" %v", len(v))
	}
	return fmt.Sprintf(msg, v...)
}

//写日志到阿里云
func Writelog(s string, topic string, ip string, config map[string]string) (err error) {

	if AliDefaultConfig.On != "1" {
		return
	}

	log := producer.GenerateLog(uint32(time.Now().Unix()), config)

	err = AliProducer.SendLog(AliDefaultConfig.Project, s, topic, ip, log)

	if err != nil {
		return
	}

	return

}

// 错误日志
func Errors(fun string, f interface{}, v ...interface{}) {

	logd := make(map[string]string)

	logd["fun"] = fun

	logd["msg"] = formatLog(f, v)

	logd["time"] = lconv.Int64ToStr(ltime.UninxTime())

	logs.Error(fun, logd["msg"], logd["time"])

	go Writelog("error", fun, "127.0.0.1", logd)
}

// 调试日志
func Debug(fun string, f interface{}, v ...interface{}) {

	logd := make(map[string]string)

	logd["fun"] = fun

	logd["msg"] = formatLog(f, v)

	logd["time"] = lconv.Int64ToStr(ltime.UninxTime())

	logs.Error(fun, logd["msg"], logd["time"])

	if AliDefaultConfig.Debug == "1" {
		go Writelog("debug", fun, "127.0.0.1", logd)
	}
}
