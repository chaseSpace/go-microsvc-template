package orm

import (
	"fmt"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"microsvc/deploy"
	"microsvc/pkg/xlog"
)

var instMap = make(map[deploy.DBname]*gorm.DB)

func InitGorm(must bool) func(*deploy.XConfig, func(must bool, err error)) {
	return func(cc *deploy.XConfig, onEnd func(must bool, err error)) {

		gconf := &gorm.Config{
			Logger:          logger.Default.LogMode(logger.Info),
			CreateBatchSize: 100, // 批量插入时，分批进行
		}
		var db *gorm.DB
		var err error
		if len(cc.Mysql) == 0 {
			fmt.Println("### there is no mysql config found")
		} else {
			for _, v := range cc.Mysql {
				db, err = gorm.Open(mysql.Open(v.Dsn()), gconf)
				if err != nil {
					fmt.Printf("\n****** failed to connect to mysql: err:%v\n", err)
					fmt.Printf("****** mysql.dsn: %s\n\n", v.Dsn())
					break
				}
				instMap[v.DBname] = db
			}
		}

		if err == nil {
			err = setupSvcDB()
		}

		onEnd(must, err)
	}
}

type MysqlObj struct {
	name deploy.DBname
	*gorm.DB
	// 你可能希望在对象中包含一些其他自定义成员，在这里添加
}

func (m *MysqlObj) IsInvalid() bool {
	return m.DB == nil
}

func (m *MysqlObj) Stop() {
	db, _ := m.DB.DB()
	err := db.Close()
	if err != nil {
		xlog.Error("orm.Stop() failed", zap.Error(err))
	}
}

func (m *MysqlObj) String() string {
	return fmt.Sprintf("mysqlObj{name:%s, instExists:%v}", m.name, m.DB != nil)
}

var servicesDB []*MysqlObj

func setupSvcDB() error {
	for _, obj := range servicesDB {
		obj.DB = instMap[obj.name]
		if obj.IsInvalid() {
			return fmt.Errorf("orm.MysqlObj is invalid, %s", obj)
		}
	}
	return nil
}

func Stop() {
	for _, gdb := range instMap {
		db, _ := gdb.DB()
		_ = db.Close()
	}
	xlog.Debug("orm-mysql: resource released...")
}

func NewMysqlObj(dbname deploy.DBname) *MysqlObj {
	return &MysqlObj{name: dbname}
}

func RegSvcDB(obj ...*MysqlObj) {
	for _, o := range obj {
		if o.name == "" {
			panic(fmt.Sprintf("orm.RegSvcDB: need name"))
		}
	}
	servicesDB = obj
}

func IgnoreNil(err error) error {
	if err == nil || err == gorm.ErrRecordNotFound {
		return nil
	}
	return err
}
