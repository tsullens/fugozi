package database

import (
  "sync"
  "bytes"
)

type lock_collection struct {
  sync.RWMutex
  collection map[string][]byte
}

func NewLockCollection() (*lock_collection) {
  return &lock_collection{
    collection: make(map[string][]byte),
  }
}

type Bucket struct {
  sync.RWMutex
  BucketId string
  //database *lock_collection `json:"-"`
  collection map[string][]byte
}

func NewBucket(name string) (*Bucket) {
  return &Bucket{
    BucketId: name,
    //database: NewLockCollection(),
    collection: make(map[string][]byte),
  }
}

func (bucket *Bucket) Get(key string) ([]byte) {
  bucket.RLock()
  defer bucket.RUnlock()
  if doc, exists := bucket.collection[key]; exists {
    return doc
  }
  return nil
}

func (bucket *Bucket) Update(key string, doc []byte) {
  bucket.Lock()
  defer bucket.Unlock()
  bucket.collection[key] = bytes.ToLower(doc)
  return
}

func (bucket *Bucket) Delete(key string) {
  bucket.Lock()
  defer bucket.Unlock()
  delete(bucket.collection, key)
}
