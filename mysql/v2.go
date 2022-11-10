package lmysql

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego/orm"
	lconv "github.com/lixu-any/go-tools/conv"
	lencry "github.com/lixu-any/go-tools/encry"
	lredis "github.com/lixu-any/go-tools/redis"
	"github.com/mitchellh/mapstructure"
)

type ListConfigV2 struct {
	TableName string    //表名
	Where     string    // where 条件, isdel=0
	Columns   string    //字符串字段
	OrderBy   string    //排序
	DBIndex   string    //指定数据库
	CacheDB   string    //数据库缓存 Redis数据库索引
	CacheName string    //缓存的名字
	Orm       orm.Ormer //实例化 ORM
	Exexpire  int       //缓存过期时间 -1不使用缓存
	MaxCount  int       //最大条数 默认100
}

func GetListV2(req ListConfigV2, retdata interface{}) (nums int64, err error) {

	var datalist []map[string]interface{}

	if req.TableName == "" || req.Where == "" || req.DBIndex == "" {
		err = fmt.Errorf("GetListV2 req.TableName empty %s", req.TableName)
		return
	}

	var qstr string

	var colums string

	if len(req.Columns) == 0 {
		colums = "*"
	} else {
		colums = req.Columns
	}

	qstr = req.Where + colums

	if req.Exexpire == 0 {
		req.Exexpire = DEFAULT_EXPIRE
	}

	if req.CacheDB == "" {
		req.CacheDB = LMysqlCacheDB
	}

	var (
		rkey   string
		_redis = lredis.LXredis{
			Name: rkey,
			DB:   req.CacheDB,
		}
	)

	if req.Exexpire > 0 {

		if req.CacheName != "" {
			_redis.Name = req.CacheName
		} else {
			_redis.Name = "db::" + req.TableName + "::" + lencry.MD5(qstr)
		}

		if b, _ := _redis.IsExist(); b {
			str, _ := _redis.GET()
			if str != "" {
				json.Unmarshal([]byte(str), &datalist)
				nums = int64(len(datalist))
				mapstructure.WeakDecode(datalist, &retdata)
				return
			}
		}

	}

	if req.Exexpire == 0 {
		req.Exexpire = DEFAULT_EXPIRE
	}

	if req.MaxCount == 0 {
		req.MaxCount = DEFAULT_MAXCOUNT //默认最大取100条
	}

	var (
		where, sql, orderby, limit string
	)

	if req.Orm == nil {
		req.Orm = orm.NewOrm()
	}

	if req.Where != "" {
		where = fmt.Sprintf(" where %s", req.Where)
	}

	if req.OrderBy != "" {
		orderby = fmt.Sprintf(" order by %s", req.OrderBy)
	}

	limit = fmt.Sprintf(" limit %d", req.MaxCount)

	sql = fmt.Sprintf("select %s from %s %s %s %s", colums, req.TableName, where, orderby, limit)

	req.Orm.Using(req.DBIndex)

	var (
		dblst []orm.Params
	)

	nums, err = req.Orm.Raw(sql).Values(&dblst)

	if err != nil {
		return
	}

	if nums < 1 {
		return
	}

	for _, v := range dblst {
		datalist = append(datalist, v)
	}

	err = mapstructure.WeakDecode(datalist, &retdata)
	if err != nil {
		return
	}

	if req.Exexpire == -1 {
		return
	}

	_redis.Val = lconv.JsonEncode(datalist)

	_redis.SetTime(req.Exexpire)

	return
}
