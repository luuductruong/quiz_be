package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

type Config struct {
	Host     string `config:"host"`
	Port     string `config:"port"`
	User     string `config:"user"`
	Password string `config:"password"`
	DbName   string `config:"dbname"`
	SSL      bool   `config:"ssl"`
}

type SQL interface {
	GetDB() *gorm.DB
	Close() error
}

type sqlStruct struct {
	DB *gorm.DB
}

func (s *sqlStruct) GetDB() *gorm.DB {
	return s.DB
}

func (s *sqlStruct) Close() error {
	s.DB = nil
	return nil
}

func NewSQL(c *Config) (SQL, error) {
	host := "host=" + c.Host
	user := "user=" + c.User
	password := "password=" + c.Password
	dbname := "dbname=" + c.DbName
	port := "port=" + c.Port
	connectionString := host + " " + user + " " + password + " " + port + " " + dbname
	if c.SSL {
		connectionString += " sslmode=require"
	} else {
		connectionString += " sslmode=disable"
	}
	db, err := gorm.Open(postgres.Open(connectionString))
	if err != nil {
		return nil, err
	}
	return &sqlStruct{DB: db}, nil
}

func InitDB(dsn string) (*gorm.DB, error) {
	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	return DB, nil
}
