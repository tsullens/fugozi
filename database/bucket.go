package database

import (
  "sync"
  "time"
  "encoding/json"
  "bytes"
)
/*
type Document struct {
  timestamp time.Time
  data []byte
}
*/
//type Document []byte

type documentData map[string]interface{}

type Document struct {
  Timestamp time.Time
  Data documentData
}

func NewDocument(docData []byte) (*Document, error) {
  var data documentData
  err := json.Unmarshal(bytes.ToLower(docData), &data)
  if err != nil {
    return nil, err
  }
  return &Document{
    Timestamp: time.Now(),
    Data: data,
  }, nil
}

type collector interface {
  Get(string) *Document
  Update(string, *Document)
  Delete(string)
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
  store collector
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
      store: NewSyncMapCollector(),
    }, nil
  default:
    return nil, &BucketEngineError{engineName: b.Engine}
  }
}

func (b *Bucket) Get(docId string) (*Document) {
  return b.store.Get(docId)
}

func (b *Bucket) Update(doc *Document) (error) {
  bmd := b.GetMetaData()
  if v, exists := doc.Data[bmd.PrimaryKey]; exists {
    docId, ok := v.(string)
    if !ok {
      return &PrimaryKeyTypeValueMismatchError{doc: doc}
    }
    b.store.Update(docId, doc)
  } else {
    return &PrimaryKeyNotFoundError{pKey: bmd.PrimaryKey, doc: doc}
  }
  return nil
}

func (b *Bucket) Delete(key string) {
  b.store.Delete(key)
}
