package models

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/project-miko/miko/conf"

	"github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

const (
	ModelYes = 1
	ModelNo  = 0
)

var (
	// global db instances
	dbs = make(map[string]*gorm.DB)
	// global redis instances
	redises = make(map[string]*RedisClient)
	pgs     = make(map[string]*pgxpool.Pool)
)

func InitModel() {

	// db_main
	err := registerDB("db_main")
	if err != nil {
		panic(err)
	}

	// rdb_main
	err = registerRedis("rdb_main")
	if err != nil {
		panic(err)
	}

	// pg_main
	err = registerPG("pg_main")
	if err != nil {
		panic(err)
	}
}

// get main db instance
func GetDbInst() *gorm.DB {
	return dbs["db_main"]
}

// get main redis instance
func GetRdbInst() *RedisClient {
	return redises["rdb_main"]
}

// get main pg instance
func GetPGInst(sectionName string) *pgxpool.Pool {
	return pgs[sectionName]
}

func registerPG(sectionName string) error {

	// "postgres://user:password@host:port/dbname"
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", conf.GetConfigString(sectionName, "user"),
		conf.GetConfigString(sectionName, "password"),
		conf.GetConfigString(sectionName, "host"),
		conf.GetConfigString(sectionName, "port"),
		conf.GetConfigString(sectionName, "name"))

	var err error
	pg, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return err
	}
	pgs[sectionName] = pg

	return nil
}

// register db instance
func registerDB(sectionName string) error {
	//"user:password@tcp(host:port)/dbname?charset=utf8&parseTime=True&loc=Local"
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4", conf.GetConfigString(sectionName, "user"),
		conf.GetConfigString(sectionName, "password"),
		conf.GetConfigString(sectionName, "host"),
		conf.GetConfigString(sectionName, "port"),
		conf.GetConfigString(sectionName, "name"))

	var err error
	dbs[sectionName], err = gorm.Open("mysql", connStr)
	if err != nil {
		return err
	}

	// todo set connection pool, actual setting reference production usage
	//db.DB().SetMaxIdleConns()

	return nil
}

// register redis instance
func registerRedis(sectionName string) error {
	// establish connection pool
	redisHost := conf.GetConfigString(sectionName, "host")
	redisPort := conf.GetConfigString(sectionName, "port")
	redisPwd := conf.GetConfigString(sectionName, "password")
	redisDb, err := conf.GetConfigInt(sectionName, "db")

	if err != nil {
		redisDb = 0
	}

	rc := &redis.Pool{
		MaxIdle:     10,               // max idle connections, these connections will be kept in the pool before they are closed
		MaxActive:   50,               // max active connections, the maximum number of connections in the pool
		IdleTimeout: 60 * time.Second, // idle connection timeout, should be shorter than redis server timeout
		Wait:        true,             // if max connections are reached, return an error or wait for a connection to be released
		Dial: func() (redis.Conn, error) {
			con, err := redis.Dial("tcp", redisHost+":"+redisPort,
				redis.DialPassword(redisPwd),
				redis.DialDatabase(int(redisDb)),
				redis.DialConnectTimeout(2*time.Second),
				redis.DialReadTimeout(2*time.Second),
				redis.DialWriteTimeout(2*time.Second))
			if err != nil {
				return nil, err
			}
			return con, nil
		},
	}

	redises[sectionName] = &RedisClient{client: rc}

	return nil
}
