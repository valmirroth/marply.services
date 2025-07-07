package sqlRoth

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"vroth.rsdata/config"
)

const sqlServerSectionConfig = "SqlServerHost"
const sqlServerConfigUser = "SqlServerUser"
const sqlServerConfigPass = "SqlServerPassword"
const serviceName = "RSData Integration Service"

var mux sync.Mutex
var databases = make(map[string]*gorm.DB)

func GetDatabase(database string) *gorm.DB {
	DataDaseSelecionado := database

	database = DataDaseSelecionado
	var err error = nil

	token := strings.ToLower(database)
	db, ok := databases[token]
	if !ok {
		mux.Lock()
		defer mux.Unlock()
		db, ok = databases[token]
		if !ok {
			db, err = initializeDatabase(database)
			if err != nil {
				fmt.Println(err)
				db = nil
			}
		}
	}
	sqlDB, err := db.DB()
	sqlDB.SetMaxIdleConns(50)
	sqlDB.SetConnMaxLifetime(time.Duration(5) * time.Minute)
	sqlDB.SetConnMaxIdleTime(time.Duration(1) * time.Minute)
	return db
}

func initializeDatabase(database string) (*gorm.DB, error) {
	//	globalConfig := settings.GetInstange()

	connectionString, err := getSqlServerConnectionString("globalConfig", database)
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(sqlserver.Open(connectionString), &gorm.Config{SkipDefaultTransaction: true, Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		return nil, err
	}

	createDefaultIndex(db)
	return db, nil
}

func getSqlServerConnectionString(configs string, database string) (string, error) {
	config.LoadEnv()
	configSection := config.GetDBConfig()

	const connStringTemplate = "sqlserver://%s:%s@%s?database=%s&app name=%s"
	sqlServerIp := configSection.Host
	sqlServerUser := configSection.User

	sqlServerPass := configSection.Password

	return fmt.Sprintf(connStringTemplate, sqlServerUser, sqlServerPass, sqlServerIp, database, serviceName), nil
}

func invalidConfig(config string) error {
	var err = fmt.Sprintf("no configuration %s for sql found", config)
	return errors.New(err)
}

func createDefaultIndex(db *gorm.DB) {
	//	db.Exec(querys.DefaultIndex)
}
