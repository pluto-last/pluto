package db

import (
	"log"
	"pluto/config"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var once = new(sync.Once)

func Instance(cfg *config.Specification) (*gorm.DB, error) {
	db, err := gorm.Open(cfg.DB.Typ, cfg.DB.DSN)
	if err != nil {
		return nil, err
	}
	if cfg.DB.Debug {
		db.LogMode(true)
	}
	if err := db.DB().Ping(); err != nil {
		return nil, err
	}
	once.Do(func() {
		db.DB().SetMaxIdleConns(50)
		db.DB().SetMaxOpenConns(200)
		//db.DB().SetConnMaxIdleTime(30 * time.Second)
		db.DB().SetConnMaxLifetime(time.Hour)
	})
	return db, nil
}

type Tx struct {
	*gorm.DB
	commit bool
}

func Begin(db *gorm.DB) *Tx {
	return &Tx{
		DB: db.Begin(),
	}
}

func (tx *Tx) RollbackIfFailed() {
	if tx.commit {
		return
	}
	if err := tx.Rollback().Error; err != nil {
		log.Println("rollback failed", err)
	}
}

func (tx *Tx) Commit() {
	if err := tx.DB.Commit().Error; err != nil {
		log.Println("commit failed", err)
		return
	}
	tx.commit = true
}

func WithUpdatedAt(kv ...interface{}) map[string]interface{} {
	var k string
	m := make(map[string]interface{})
	for i, v := range kv {
		if i%2 != 0 {
			m[k] = v
			continue
		}
		k = v.(string)
	}
	m["updated_at"] = time.Now()
	return m
}
