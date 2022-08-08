package lmysql

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	lconv "github.com/lixu-any/go-tools/conv"
	lencry "github.com/lixu-any/go-tools/encry"
	lredis "github.com/lixu-any/go-tools/redis"
)

var LMysqlCacheDB string //默认缓存数据库

type MapMysqlConn struct {
	Name   string
	SqlCon string
}

//初始化MySQL
func InitMysqls(conns []MapMysqlConn, cachedb string) (err error) {

	LMysqlCacheDB = cachedb

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
	TableName    string    //表名
	Where        string    // where 条件, isdel=0
	ColumnString []string  //字符串字段
	ColumnInt    []string  //数组字段
	OrderBy      string    //排序
	DBIndex      string    //指定数据库
	CacheDB      string    //数据库缓存 Redis数据库索引
	Orm          orm.Ormer //实例化 ORM
	Exexpire     int       //缓存过期时间 -1不使用缓存
	MaxCount     int       //最大条数 默认100
}

type ExecConfig struct {
	TableName  string                 //表名
	Where      string                 // where 条件, isdel=0
	Data       map[string]interface{} //要添加 修改的数据，会自动判断变量类型
	DBIndex    string                 //指定数据库
	Orm        orm.Ormer              //实例化 ORM
	CreaeSql   string                 //创建表的SQL语句
	CreatePath string                 //创建表的SQL文件路径
}

const (
	DEFAULT_EXPIRE   = 300 //默认缓存5分钟
	DEFAULT_MAXCOUNT = 100 //默认最大条数
)

func _getList(req ListConfig) (nums int64, lst []map[string]interface{}, err error) {

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

func GetList(req ListConfig) (nums int64, lst []map[string]interface{}, err error) {
	return _getList(req)
}

//获取缓存的key名称
func GetCacheName(tbname, where string, colstr, colint []string, cachedb ...string) string {

	var qstr string

	var colums string

	if len(colstr) == 0 && len(colint) == 0 {
		colums = "*"
	} else {
		ss := append(colstr, colint...)
		colums = strings.Join(ss, ",")
	}

	qstr = where + colums

	return "db::" + tbname + "::" + lencry.MD5(qstr)

}

//删除cache
func DeleteCache(tbname, where string, colstr, colint []string, cachedb ...string) (err error) {

	rkey := GetCacheName(tbname, where, colstr, colint)

	var cdb = LMysqlCacheDB
	if len(cachedb) > 0 {
		cdb = cachedb[0]
	}

	_redis := lredis.LXredis{
		Name: rkey,
		DB:   cdb,
	}

	err = _redis.DELETE()

	return
}

//执行SQL语句
func Exec(req ExecConfig) (id int64, err error) {

	if req.TableName == "" || req.DBIndex == "" {
		err = fmt.Errorf("req.TableName empty %s", req.TableName)
		return
	}

	var (
		sql    string
		count  int
		keys   string
		vals   string
		tcount int
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
				keys += fmt.Sprintf("%s='%s'", k, v)
			case float64:
				keys += fmt.Sprintf("%s='%f'", k, v)
			case int, int32, int64, int8:
				keys += fmt.Sprintf("%s='%d'", k, v)
			default:
				keys += fmt.Sprintf("%s='%s'", k, v)
			}

			count++

		}
		sql = fmt.Sprintf("update %s set %s where %s", req.TableName, keys, req.Where)
	}

	if req.Orm == nil {
		req.Orm = orm.NewOrm()
	}

RSTART:

	req.Orm.Using(req.DBIndex)

	res, err := req.Orm.Raw(sql).Exec()

	if err != nil {
		if b := strings.Contains(err.Error(), "doesn't exist"); b && (req.CreaeSql != "" || req.CreatePath != "") && tcount == 0 {

			if req.CreaeSql == "" {
				req.CreaeSql, _ = readFilesql(req.CreatePath, req.TableName)
			}

			//表不存在 去创建
			err = CreateTable(req.CreaeSql, req.Orm, req.DBIndex)
			if err != nil {
				return
			}

			tcount++
			//重新执行一次
			goto RSTART

		}
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

type ExecSqlConfig struct {
	Sql        string    //表名
	TableName  string    //表名
	DBIndex    string    //指定数据库
	Orm        orm.Ormer //实例化 ORM
	CreaeSql   string    //创建表的SQL语句
	CreatePath string    //创建表的SQL文件路径
}

func ExecSql(config ExecSqlConfig) (err error) {

	var (
		tcount int
	)

	if config.Orm == nil {
		config.Orm = orm.NewOrm()
	}

	if config.DBIndex == "" {
		config.DBIndex = "default"
	}

RSTART:

	config.Orm.Using(config.DBIndex)

	_, err = config.Orm.Raw(config.Sql).Exec()
	if err != nil {
		if b := strings.Contains(err.Error(), "doesn't exist"); b && (config.CreaeSql != "" || config.CreatePath != "") && tcount == 0 {
			if config.CreaeSql == "" {
				config.CreaeSql, _ = readFilesql(config.CreatePath, config.TableName)
			}

			//表不存在 去创建
			err = CreateTable(config.CreaeSql, config.Orm, config.DBIndex)
			if err != nil {
				return
			}

			tcount++
			//重新执行一次
			goto RSTART
		}

		return

	}

	return
}

func CreateTable(csql string, o orm.Ormer, dbindex string) (err error) {

	o.Using(dbindex)

	_, err = o.Raw(csql).Exec()
	if err != nil {
		return
	}

	return
}

func readFilesql(pth, tbname string) (sql string, err error) {
	filePtr, err := os.Open(pth)
	if err != nil {
		return
	}
	defer filePtr.Close()

	by, err := ioutil.ReadAll(filePtr)
	if err != nil {
		return
	}

	sql = strings.Replace(string(by), "{tablename}", tbname, -1)

	return
}
