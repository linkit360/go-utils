package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/jinzhu/configor"
)

type DataBaseConfig struct {
	ConnMaxLifetime  int    `default:"-1" yaml:"conn_ttl"`
	MaxOpenConns     int    `default:"15" yaml:"max_open_conns"`
	MaxIdleConns     int    `default:"5" yaml:"max_idle_conns"`
	ReconnectTimeout int    `default:"10" yaml:"timeout"`
	User             string `default:""`
	Pass             string `default:""`
	Port             string `default:""`
	Name             string `default:""`
	Host             string `default:""`
	SSLMode          string `default:"disable" yaml:"ssl_mode"`
	TablePrefix      string `default:"xmp_" yaml:"table_prefix"`
}

func (dbConfig DataBaseConfig) GetConnStr() string {
	dsn := "postgres://" +
		dbConfig.User + ":" +
		dbConfig.Pass + "@" +
		dbConfig.Host + ":" +
		dbConfig.Port + "/" +
		dbConfig.Name + "?sslmode=" +
		dbConfig.SSLMode
	return dsn
}

func Init(conf DataBaseConfig) *sql.DB {
	var dbConn *sql.DB

	var err error
	dbConn, err = sql.Open("postgres", conf.GetConnStr())
	if err != nil {
		fmt.Printf("open error %s, dsn: %s", err.Error(), conf.GetConnStr())
		log.WithField("error", err.Error()).Fatal("db connect")
	}
	if err = dbConn.Ping(); err != nil {
		fmt.Printf("ping error %s, dsn: %s", err.Error(), conf.GetConnStr())
		log.WithField("error", err.Error()).Fatal("db ping")
	}

	dbConn.SetMaxOpenConns(conf.MaxOpenConns)
	dbConn.SetMaxIdleConns(conf.MaxIdleConns)
	dbConn.SetConnMaxLifetime(time.Second * time.Duration(conf.ConnMaxLifetime))

	log.WithFields(log.Fields{
		"host": conf.Host, "dbname": conf.Name, "user": conf.User}).Info("database connected")
	return dbConn
}

func InitName(e, path string) *sql.DB {
	var cfg map[string]DataBaseConfig
	if err := configor.Load(&cfg, path); err != nil {
		log.WithField("config", err.Error()).Fatal("config load error")
	}
	if dbConf, ok := cfg[e]; ok {
		return Init(dbConf)
	}
	log.WithField("error", "no such dbmap").Fatal("init database")
	return nil
}
