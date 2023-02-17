package DataBase

import (
	"Core/Log"
	"Core/Stl"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"strings"
	"time"
)

const (
	maxConn   = 120
	maxPacket = 3 * 1024 * 1024
)

type callBack func(sql string, err error, args ...interface{})

type mysqlCallBack struct {
	fun  callBack
	sql  string
	args []interface{}
}

var mysqlDb *sqlx.DB
var mysqlDbErr error
var mysqlQueue []*Stl.Queue
var connIndex int64

// 初始化链接
func MysqlInit(userName string, passWord string, ip string, port int, dbName string) {
	dbDSN := fmt.Sprintf("%s:%s@tcp(%s:%v)/%s?charset=%s", userName, passWord, ip, port, dbName, "utf8")

	// 打开连接失败
	mysqlDb, mysqlDbErr = sqlx.Open("mysql", dbDSN)
	if mysqlDbErr != nil {
		Log.ErrorLog("数据源配置不正确: " + mysqlDbErr.Error())
		return
	}

	// 最大连接数
	mysqlDb.SetMaxOpenConns(maxConn)
	// 闲置连接数
	mysqlDb.SetMaxIdleConns(20)
	// 最大连接周期
	mysqlDb.SetConnMaxLifetime(200 * time.Second)

	if mysqlDbErr = mysqlDb.Ping(); nil != mysqlDbErr {
		Log.ErrorLog("数据库链接失败: " + mysqlDbErr.Error())
		return
	}
	for i := 0; i < maxConn; i++ {
		mysqlQueue = append(mysqlQueue, Stl.NewQueue())
		go mysqlExec(mysqlQueue[i])
	}
}

func mysqlExec(queue *Stl.Queue) {
	for {
		if queue.Head() != nil {
			for queue.Head() != nil {
				data := queue.Dequeue().(*mysqlCallBack)
				err := mysqlexec(data.sql)
				if data.fun != nil {
					data.fun(data.sql, err, data.args...)
				}
			}
		}
		time.Sleep(time.Microsecond)
	}
}

func mysqlexec(sql string) error {
	_, err := mysqlDb.Exec(sql)
	if err != nil {
		Log.ErrorLog("Failed to mysqlExec, err = %v, sql = %v", err, sql)

	}
	return err
}

func mysqlIntoQueue(sql string, back callBack, param ...interface{}) {
	if len(sql) < 1 {
		Log.ErrorLog("Failed to mysqlIntoQueue, sql is empty")
		return
	}
	mysqlQueue[connIndex%maxConn].Enqueue(&mysqlCallBack{
		fun:  back,
		sql:  sql,
		args: param,
	})
	connIndex++
}

func MysqlExec(sql string, back callBack, param ...interface{}) {
	mysqlIntoQueue(sql, back, param...)
}

// 插入数据
func MysqlInsert(tabName string, data map[string]interface{}, back callBack, param ...interface{}) {
	fields, values := insertToString(data)
	sql := fmt.Sprintf("insert into %s (%s) values (%s)", tabName, fields, values)
	mysqlIntoQueue(sql, back, param...)
}

func MysqlInsertSyn(tabName string, data map[string]interface{}) (int64, error) {
	fields, values := insertToString(data)

	sql := fmt.Sprintf("insert into %s (%s) values (%s)", tabName, fields, values)

	r, err := mysqlDb.Exec(sql)
	if err != nil {
		Log.ErrorLog("Failed to MysqlInsertSyn, err = %v, sql = %v", err, sql)
		return 0, err
	}
	id, err := r.LastInsertId()
	if err != nil {
		Log.ErrorLog("Failed to MysqlInsertSyn LastInsertId, err = %v", err)
		return 0, err
	}
	return id, nil
}

func MysqlInsertMany(tabName string, manyData []map[string]interface{}, back callBack, args ...interface{}) {
	fields := ""
	var fieldList []string
	bFirst := true
	var sql strings.Builder
	sql.Grow(maxPacket * 2)
	for _, data := range manyData {
		values := ""
		if bFirst {
			sql.Reset()
			fields, values = insertToString(data)
			sql.WriteString(fmt.Sprintf("insert into %s (%s) values (%s)", tabName, fields, values))
			fieldList = strings.Split(fields, ",")
			bFirst = false
		} else {
			for _, field := range fieldList {
				value := data[field]
				if value == nil {
					values = fmt.Sprintf("%sNull,", values)
				} else {
					values = fmt.Sprintf("%s\"%v\",", values, value)
				}
			}
			values = values[:len(values)-1]
			sql.WriteString(fmt.Sprintf(",(%s)", values))
			if sql.Len() >= maxPacket {
				mysqlIntoQueue(sql.String(), back, args...)
				bFirst = true
			}
		}
	}

	mysqlIntoQueue(sql.String(), back, args...)
}

// 插入数据
func MysqlUpdate(tabName string, condition map[string]interface{}, data map[string]interface{}, back callBack, args ...interface{}) {
	sql := fmt.Sprintf("update %s set %s where %s", tabName, dataToString(data), conditionToString(condition))
	mysqlIntoQueue(sql, back, args...)
}

// 删除数据
func MysqlDelete(tabName string, condition map[string]interface{}, data map[string]interface{}, back callBack, args ...interface{}) {
	sql := fmt.Sprintf("delete from %s where %s", tabName, conditionToString(condition))
	mysqlIntoQueue(sql, back, args...)
}

func MySqlGetData(tableName string, condition map[string]interface{}, data interface{}, fields ...string) {
	sql := selectToString(tableName, condition, fields...)
	err := mysqlDb.Select(data, sql)
	if err != nil {
		Log.ErrorLog("数据库查询失败: err = %v, sql = %v", err, sql)
	}
}
