package database

import (
  "sync"
  "go-cached/util"
)

type Database struct {
  sync.RWMutex
  buckets map[string]*Bucket
  Logger *util.LumberJack
}

func NewDatabase() (*Database) {
  return &Database{
    buckets: make(map[string]*Bucket),
    Logger: util.NewLumberJack(util.Config.DbLog),
  }
}
