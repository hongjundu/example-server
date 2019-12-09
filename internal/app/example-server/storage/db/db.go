package db

import (
	_ "example-server/internal/app/example-server/storage/model"
	"example-server/internal/pkg/env"
	"fmt"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/hongjundu/go-level-logger"
)

func Init() error {
	logger.Debugf("[db] Init")

	connString := GetConnString()
	logger.Debugf("[storage] db connect string: %s", connString)

	err := orm.RegisterDataBase("default", "mysql", connString, 30)

	if err != nil {
		logger.Errorf("[db] ORM RegisterDataBase failed: %+v", err)
		return err
	}

	// Create tables if not exist
	if err = createTables(); err != nil {
		logger.Errorf("[db] create tables failed: %+v", err)
		return err
	}

	return nil

}

func GetConnString() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&loc=Local", env.Env.MySqlUser, env.Env.MySqlPassword, env.Env.MySqlHost, env.Env.MySqlDb)
}

func GetCasbinConnString() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/", env.Env.MySqlUser, env.Env.MySqlPassword, env.Env.MySqlHost)
}

func createTables() error {
	logger.Debugf("[db] createTables")

	name := "default" // database alias
	force := false    // do not drop table before creating
	verbose := true   // print execute process

	var err error
	if err = orm.RunSyncdb(name, force, verbose); err == nil {
		logger.Infoln("[storage] create tables successfully")
	}

	return err
}
