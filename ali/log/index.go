package lalilog

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
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

	TableError string
	TableDebug string
	Env        string
}

//初始化阿里云日志
func InitAliLog(config MapAlilogConfig) {

	AliDefaultConfig = config

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

	if AliDefaultConfig.Env != "" {
		config["env"] = AliDefaultConfig.Env
	}

	log := producer.GenerateLog(uint32(time.Now().Unix()), config)

	err = AliProducer.SendLog(AliDefaultConfig.Project, s, topic, ip, log)

	if err != nil {
		return
	}

	return

}

// 获取正在运行的函数名
func runFuncName(l int) (string, string, int) {

	pc, file, line, _ := runtime.Caller(l)

	name := runtime.FuncForPC(pc).Name()

	split := strings.Split(name, ".")

	funname := split[len(split)-1]

	return funname, file, line
}

// 错误日志
func Errors(f interface{}, v ...interface{}) {

	funname, file, line := runFuncName(2)

	logd := make(map[string]string)

	logd["fun"] = funname

	logd["file"] = file

	logd["line"] = fmt.Sprintf("%d", line)

	logd["msg"] = formatLog(f, v)

	logd["time"] = lconv.Int64ToStr(ltime.UninxTime())

	logs.Error("file::", file, ",line::", line, ",fun::", funname, "[", logd["msg"], "],time::", logd["time"])

	if AliDefaultConfig.TableError != "" {
		go Writelog(AliDefaultConfig.TableError, funname, "127.0.0.1", logd)
	}

}

// 调试日志
func Debug(f interface{}, v ...interface{}) {

	funname, file, line := runFuncName(2)

	logd := make(map[string]string)

	logd["fun"] = funname

	logd["file"] = file

	logd["line"] = fmt.Sprintf("%d", line)

	logd["msg"] = formatLog(f, v)

	logd["time"] = lconv.Int64ToStr(ltime.UninxTime())

	logs.Debug("file::", file, ",line::", line, ",fun::", funname, "[", logd["msg"], "],time::", logd["time"])

	if AliDefaultConfig.Debug == "1" && AliDefaultConfig.TableDebug != "" {
		go Writelog(AliDefaultConfig.TableDebug, funname, "127.0.0.1", logd)
	}
}
