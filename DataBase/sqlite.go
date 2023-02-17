package DataBase

import (
	"Core/Log"
	"Core/Stl"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

var sqliteDB *sqlx.DB
var sqliteQueue = Stl.NewQueue()
var sqliteQueueIndex = 0

func CloseSqlLite() {
	sqliteDB.Close()
}

func InitSqlLite(fileName string) {
	db, err := sqlx.Open("sqlite3", fileName)
	if err != nil {
		Log.ErrorLog("Failed to InitSqlLite, err = %v, fileName = %v", err, fileName)
		return
	}
	sqliteDB = db
	go sqliteExec()
}

func SqliteInsert(tableName string, data map[string]interface{}) {
	fields, values := insertToString(data)
	sql := fmt.Sprintf("insert into %s (%s) values (%s)", tableName, fields, values)
	sqliteIntoQueue(sql)
}

func SqliteUpdate(tabName string, condition map[string]interface{}, data map[string]interface{}) {
	sql := fmt.Sprintf("update %s set %s where %s", tabName, dataToString(data), conditionToString(condition))
	sqliteIntoQueue(sql)
}

func SqliteDelete(tabName string, condition map[string]interface{}) {
	sql := fmt.Sprintf("delete from %s where %s", tabName, conditionToString(condition))
	sqliteIntoQueue(sql)
}

func SqliteGetData(tabName string, condition map[string]interface{}, data interface{}, fields ...string) {
	sql := selectToString(tabName, condition, fields...)
	err := sqliteDB.Select(data, sql)
	if err != nil {
		Log.ErrorLog("数据库查询失败: err = %v, sql = %v", err, sql)
	}
}

func SqliteTableExist(tabName string) bool {
	var count []int
	if err := sqliteDB.Select(&count, fmt.Sprintf("SELECT COUNT(*) FROM sqlite_master where type = 'table' and name = '%v'", tabName)); err != nil {
		return false
	}
	return len(count) == 1 && count[0] > 0
}

func SqliteExec(sql string) {
	sqliteIntoQueue(sql)
}

func SqliteExecSyn(sql string) {
	sqliteexec(sql)
}

func sqliteExec() {
	for {
		if sqliteQueue.Head() != nil {
			sqliteexec("BEGIN")
			for sqliteQueue.Head() != nil {
				sqliteexec(sqliteQueue.Dequeue().(string))
			}
			sqliteexec("COMMIT")
		}
		time.Sleep(time.Microsecond)
	}
}

func sqliteIntoQueue(sql string) {
	if len(sql) < 1 {
		Log.ErrorLog("Failed to sqliteIntoQueue, sql is empty")
		return
	}
	sqliteQueue.Enqueue(sql)
}

func sqliteexec(sql string) {
	_, err := sqliteDB.Exec(sql)
	if err != nil {
		Log.ErrorLog("Failed to sqliteExec, err = %v, sql = %v", err, sql)
	}
}
