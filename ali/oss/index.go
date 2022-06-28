package loss

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

//声明全局 oss对象
var OssClient *oss.Client
var OssBucketDefault *oss.Bucket
var OssConfig MapOssConfig

type MapOssConfig struct {
	Endpoint        string
	AccessKeyId     string
	AccessKeySecret string
	Bucket          string
	Domain          string
}

//初始化阿里云oss
func InitOss(config MapOssConfig) (err error) {
	// 创建OSSClient实例。
	// yourEndpoint填写Bucket对应的Endpoint，以华东1（杭州）为例，填写为https://oss-cn-hangzhou.aliyuncs.com。其它Region请按实际情况填写。
	// 阿里云账号AccessKey拥有所有API的访问权限，风险很高。强烈建议您创建并使用RAM用户进行API访问或日常运维，请登录RAM控制台创建RAM用户。

	OssConfig = config

	OssClient, err = oss.New(OssConfig.Endpoint, OssConfig.AccessKeyId, OssConfig.AccessKeySecret)
	if err != nil {
		return
	}

	OssBucketDefault, err = OssClient.Bucket(OssConfig.Bucket)
	if err != nil {
		return
	}

	return

}

//上传文件
//localpath 本地文件路径
//objectpath oss路径
func UploadFile(localpath, objectpath string) (fullpath string, err error) {

	err = OssBucketDefault.PutObjectFromFile(objectpath, localpath)
	if err != nil {
		return
	}

	fullpath = OssConfig.Domain + objectpath

	return
}
