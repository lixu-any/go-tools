package l_mysql

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	lconv "github.com/lixu-any/go-tools/conv"
	lencry "github.com/lixu-any/go-tools/encry"
	lredis "github.com/lixu-any/go-tools/redis"
)

type MapMysqlConn struct {
	Name   string
	SqlCon string
}

//初始化MySQL
func InitMysqls(conns []MapMysqlConn) (err error) {

	err = orm.RegisterDriver("mysql", orm.DRMySQL)
	if err != nil {
		logs.Error("InitMysqls", "orm.RegisterDriver::", err)
		return
	}

	for _, v := range conns {
		err = orm.RegisterDataBase(v.Name, "mysql", v.SqlCon)
		if err != nil {
			return
		}
	}

	return
}

// ---------------------------
type ListConfig struct {
	TableName    string
	Where        string
	ColumnString []string
	ColumnInt    []string
	OrderBy      string
	DBIndex      string
	CacheDB      string
	Orm          orm.Ormer
	Exexpire     int
	MaxCount     int
}

type ExecConfig struct {
	TableName string
	Where     string
	Data      map[string]interface{}
	DBIndex   string
	Orm       orm.Ormer
}

const (
	DEFAULT_EXPIRE   = 300 //默认缓存5分钟
	DEFAULT_MAXCOUNT = 100 //默认最大条数
)

func GetList(req ListConfig) (nums int64, lst []map[string]interface{}, err error) {

	if req.TableName == "" || req.Where == "" || req.DBIndex == "" {
		err = fmt.Errorf("req.TableName empty %s", req.TableName)
		return
	}

	var qstr string

	var colums string

	if len(req.ColumnString) == 0 && len(req.ColumnInt) == 0 {
		colums = "*"
	} else {
		ss := append(req.ColumnString, req.ColumnInt...)
		colums = strings.Join(ss, ",")
	}

	qstr = req.Where + colums

	if req.Exexpire == 0 {
		req.Exexpire = DEFAULT_EXPIRE
	}

	var (
		rkey   string
		_redis = lredis.LXredis{
			Name: rkey,
			DB:   req.CacheDB,
		}
	)

	if req.Exexpire > 0 {
		rkey = "db::" + req.TableName + "::" + lencry.MD5(qstr)

		_redis.Name = rkey

		if b, _ := _redis.IsExist(); b {
			str, _ := _redis.GET()
			if str != "" {
				json.Unmarshal([]byte(str), &lst)
				nums = int64(len(lst))
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
		for _, i := range req.ColumnInt {
			if v[i] != nil {
				v[i] = lconv.StrToInt64(v[i].(string))
				continue
			}
		}
		lst = append(lst, v)
	}

	if req.Exexpire == -1 {
		return
	}

	_redis.Val = lconv.JsonEncode(lst)

	_redis.SetTime(req.Exexpire)

	return
}

//执行SQL语句
func Exec(req ExecConfig) (id int64, err error) {

	if req.TableName == "" || req.DBIndex == "" {
		err = fmt.Errorf("req.TableName empty %s", req.TableName)
		return
	}

	var (
		sql   string
		count int
		keys  string
		vals  string
	)

	if req.Where == "" {

		for k, v := range req.Data {

			if count > 0 {
				keys += ","
				vals += ","
			}

			switch v := v.(type) {
			case string:
				vals += fmt.Sprintf("'%s'", v)
			case float64:
				vals += fmt.Sprintf("'%f'", v)
			case int, int32, int64, int8:
				vals += fmt.Sprintf("'%d'", v)
			default:
				vals += fmt.Sprintf("'%s'", v)
			}

			keys += k

			count++

		}

		sql = fmt.Sprintf("insert into %s (%s) values (%s)", req.TableName, keys, vals)

	} else {

		for k, v := range req.Data {

			if count > 0 {
				keys += ","
			}

			switch v := v.(type) {
			case string:
				keys += fmt.Sprintf("%s=%s", k, v)
			case float64:
				keys += fmt.Sprintf("%s=%f", k, v)
			case int, int32, int64, int8:
				keys += fmt.Sprintf("%s=%d", k, v)
			default:
				keys += fmt.Sprintf("%s=%s", k, v)
			}

			count++

		}
		sql = fmt.Sprintf("update %s set %s where %s", req.TableName, keys, req.Where)
	}

	if req.Orm == nil {
		req.Orm = orm.NewOrm()
	}

	req.Orm.Using(req.DBIndex)

	res, err := req.Orm.Raw(sql).Exec()

	if err != nil {
		return
	}

	if req.Where == "" {
		id, err = res.LastInsertId()
		if err != nil {
			return
		}
	}

	return
}
