package db

import (
	"fmt"

	"github.com/kdjuwidja/aishoppercommon/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MySQLConnectionPool struct {
	User         string
	Password     string
	Host         string
	Port         string
	DBName       string
	MaxOpenConns int
	MaxIdleConns int

	models []interface{}

	db *gorm.DB
}

func InitializeMySQLConnectionPool(user, password, host, port, dbName string, maxOpenConns, maxIdleConns int, models []interface{}) (*MySQLConnectionPool, error) {
	c := &MySQLConnectionPool{}
	c.User = user
	c.Password = password
	c.Host = host
	c.Port = port
	c.DBName = dbName
	c.MaxOpenConns = maxOpenConns
	c.MaxIdleConns = maxIdleConns
	c.models = models

	if c.User == "" || c.Password == "" || c.Host == "" || c.Port == "" || c.DBName == "" {
		return nil, fmt.Errorf("missing required configuration parameters")
	}

	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		c.User, c.Password, c.Host, c.Port, c.DBName)
	logger.Debugf("dsn: %s", dsn)

	c.db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := c.db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(c.MaxOpenConns)
	sqlDB.SetMaxIdleConns(c.MaxIdleConns)
	return c, nil
}

func (c *MySQLConnectionPool) AutoMigrate() error {
	if len(c.models) == 0 {
		return fmt.Errorf("no models to migrate")
	}

	err := c.db.AutoMigrate(c.models...)

	if err != nil {
		return err
	}

	return nil
}

func (c *MySQLConnectionPool) DropTables() error {
	if len(c.models) == 0 {
		return fmt.Errorf("no models to drop")
	}

	return c.db.Migrator().DropTable(c.models...)
}

func (c *MySQLConnectionPool) GetDB() *gorm.DB {
	return c.db
}

func (c *MySQLConnectionPool) Close() error {
	if c.db != nil {
		sqlDB, err := c.db.DB()
		if err != nil {
			return fmt.Errorf("error getting underlying *sql.DB: %v", err)
		}
		return sqlDB.Close()
	}
	return nil
}
