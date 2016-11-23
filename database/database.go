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

func (*Database db) Select(bucketId string, docId string) ([]byte, error){

}

func (*Database db) Insert(bucketId, docId string, doc []byte) (error) {

}

func (*Database db) Delete(bucketId, docId string) (error) {

}
