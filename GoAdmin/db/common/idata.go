package common

import (
	"github.com/7058011439/haoqbb/Log"
	"github.com/jinzhu/gorm"
	"sync"
	"time"
)

type Manager struct {
	mysql    *gorm.DB
	memData  map[string]map[interface{}]IDataDB // map[表名]map[id]表记录
	memMutex sync.Mutex
}

func (m *Manager) MysqlDB() *gorm.DB {
	return m.mysql
}

func (m *Manager) ResetMemData() {
	m.memData = map[string]map[interface{}]IDataDB{}
}

func (m *Manager) loadMemData(db IDataDB) IDataDB {
	m.memMutex.Lock()
	defer m.memMutex.Unlock()
	if tabData, ok := m.memData[db.TableName()]; ok {
		return tabData[db.GetKey()]
	}
	return nil
}

func (m *Manager) loadMysqlData(db IDataDB) {
	if err := m.mysql.Where(db).Find(db).Error; err != nil && err.Error() != "record not found" {
		Log.Warn("查询[%v]失败, err = %v, data = %v", db.TableName(), err, db)
	}
}

func (m *Manager) LoadData(db IDataDB) IDataDB {
	if data := m.loadMemData(db); data == nil {
		if m.loadMysqlData(db); db.IsValid() {
			m.insertDataToMem(db)
		}
		return db
	} else {
		db = data
		return data
	}
}

func (m *Manager) insertDataToMysql(db IDataDB) error {
	return m.mysql.Create(db).Error
}

func (m *Manager) insertDataToMem(db IDataDB) error {
	m.memMutex.Lock()
	defer m.memMutex.Unlock()
	if tabData, ok := m.memData[db.TableName()]; ok {
		tabData[db.GetKey()] = db
	} else {
		tabData = map[interface{}]IDataDB{}
		tabData[db.GetKey()] = db
		m.memData[db.TableName()] = tabData
	}
	return nil
}

func (m *Manager) InsertData(db IDataDB) error {
	if err := m.insertDataToMysql(db); err != nil {
		Log.Warn("插入[%v]失败, err = %v, data = %v", db.TableName(), err, db)
		return err
	} else {
		m.insertDataToMem(db)
	}
	return nil
}

func (m *Manager) updateMysql(db IDataDB) error {
	return m.mysql.Model(db).Omit(db.Omit()...).Update(db).Error
}

func (m *Manager) UpdateData(db IDataDB) error {
	if err := m.updateMysql(db); err != nil {
		Log.Warn("更新[%v]失败, err = %v, data = %v", db.TableName(), err, db)
		return err
	} else {
		return m.deleteMemData(db)
	}
}

func (m *Manager) deleteMysqlData(db IDataDB, ids ...int64) error {
	if db.IsValid() {
		ids = append(ids, db.GetId())
	}
	for _, id := range ids {
		db.SetId(id)
		if err := m.mysql.Delete(db).Error; err != nil {
			return err
		}
	}
	return nil
}

func (m *Manager) deleteMemData(db IDataDB, ids ...int64) error {
	m.memMutex.Lock()
	defer m.memMutex.Unlock()

	if tabData, ok := m.memData[db.TableName()]; ok {
		if db.IsValid() {
			delete(tabData, db.GetKey())
		}
		for _, id := range ids {
			db.SetId(id)
			delete(tabData, db.GetKey())
		}
	}
	return nil
}

func (m *Manager) DeleteData(db IDataDB, ids ...int64) error {
	if err := m.deleteMysqlData(db, ids...); err != nil {
		Log.Warn("删除[%v]失败, err = %v, db = %v, ids = %v", db.TableName(), err, db, ids)
		// 这个地方的err有值，可能只是因为部分记录删除失败(还有部分成功了)，因为一些原因，不好筛选，所以干脆直接清理内存
		m.deleteMemData(db, ids...)
		return err
	} else {
		m.deleteMemData(db, ids...)
	}
	return nil
}

func (m *Manager) ClearMemData(db IDataDB, ids ...int64) error {
	return m.deleteMemData(db, ids...)
}

func NewManager(mysql *gorm.DB) *Manager {
	return &Manager{
		mysql:   mysql,
		memData: map[string]map[interface{}]IDataDB{},
	}
}

type IChild interface {
	GetParentId() int64
	AddChild(child IChild)
	GetId() int64
}

type IDataDB interface {
	TableName() string
	IsValid() bool
	GetKey() interface{}
	Omit() []string
	SetId(int64)
	GetId() int64
	Reset()
}

type Model struct {
	ID        int64     `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (m *Model) GetKey() interface{} {
	return m.ID
}

func (m *Model) SetId(id int64) {
	m.ID = id
}

func (m *Model) IsValid() bool {
	return m.ID > 0
}

func (m *Model) Omit() []string {
	return []string{"id", "created_at"}
}

func (m *Model) GetId() int64 {
	return m.ID
}

func (m *Model) Reset() {
	m.ID = 0
}

type IUpdateData interface {
	IDataDB
	SetUpdateBy(int64)
	SetCreateBy(int64)
}

type ControlBy struct {
	CreateBy int64 `json:"createBy" gorm:"index:idx_createBy;comment:'创建者'"`
	UpdateBy int64 `json:"updateBy" gorm:"index:idx_updateBy;comment:'更新者'"`
}

func (c *ControlBy) SetCreateBy(id int64) {
	c.CreateBy = id
}

func (c *ControlBy) SetUpdateBy(id int64) {
	c.UpdateBy = id
}
