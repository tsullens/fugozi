package database

import (
  "sync"
  "bytes"
  "time"
)
/*
type document struct {
  timestamp time.Time
  data []byte
}
*/
//type document []byte

type DocumentData map[string]{}interface

type Document struct {
  timestamp time.Time
  data DocumentData
}

type BucketMetaData struct {
  sync.RWMutex
  BucketId        string `json:"bucketid" binding:"required"`
  docCount        int
  PrimaryKey      string `json:"primarykey" binding:"required"`
  SecondaryKeys   []string
  Engine          string `json:"engine" binding:"required"`
}

func (bmd *BucketMetaData) updateDocCount(n int) {
  bmd.Lock()
  defer bmd.Unlock()
  bmd.docCount = bmd.docCount + n
}

type Bucket struct {
  *BucketMetaData
  collector
}

func (bmd *BucketMetaData) GetMetaData() (BucketMetaData) {
  bmd.RLock()
  defer bmd.RUnlock()
  return *bmd
}

func NewBucket(b BucketMetaData) (*Bucket, error) {
  switch b.Engine {
  case "syncmap":
    return &Bucket{
      BucketMetaData: &b,
      collector: NewSyncMapCollector(),
    }, nil
  default:
    return nil, &BucketEngineError{engineName: b.Engine}
  }
}

type collector interface {
  Get(string) document
  Update(string, []byte)
  Delete(string)
}

type syncmap_collector struct {
  sync.RWMutex
  collection map[string]document
}

func NewSyncMapCollector() (*syncmap_collector) {
  return &syncmap_collector{
    collection: make(map[string]document),
  }
}

func (c *syncmap_collector) Get(key string) (document) {
  c.RLock()
  defer c.RUnlock()
  if doc, exists := c.collection[key]; exists {
    return doc
  }
  return nil
}

func (c *syncmap_collector) Update(doc Document) {
  c.Lock()
  defer c.Unlock()
  c.collection[key] = bytes.ToLower(doc)
  return
}

func (c *syncmap_collector) Delete(key string) {
  c.Lock()
  defer c.Unlock()
  delete(c.collection, key)
}
