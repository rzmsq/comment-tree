package db

import (
	"CommentTree/commentTree/config"
	"CommentTree/commentTree/pkg/logger"

	"github.com/wb-go/wbf/dbpg"
)

type DB struct {
	log logger.Interface
	db  *dbpg.DB
}

func NewDB(log logger.Interface, cfg config.DBConfig) (*DB, error) {
	opts := &dbpg.Options{MaxOpenConns: cfg.MaxOpenConns, MaxIdleConns: cfg.MaxIdleConns}
	db, err := dbpg.New(cfg.MasterDSN, nil, opts)
	if err != nil {
		return nil, err
	}

	return &DB{
		log: log,
		db:  db,
	}, nil
}
