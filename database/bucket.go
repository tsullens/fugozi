package database

import (
  "sync"
  "bytes"
)

type Bucket interface {
  Get() []byte
  Update()
  Delete()
}

type LockBucket struct {
  BucketId string
  *lock_collection `json:"-"`
}

func NewLockBucket(name string) (*LockBucket) {
  return &LockBucket{
    BucketId: name,
    //database: NewLockCollection(),
    collection: make(map[string][]byte),
  }
}

func (bucket LockBucket) Get(key string) ([]byte) {
  bucket.collection.RLock()
  defer bucket.RUnlock()
  if doc, exists := bucket.collection[key]; exists {
    return doc
  }
  return nil
}

func (bucket LockBucket) Insert(key string, doc []byte) (error){
  bucket.collection.Lock()
  defer bucket.collection.Unlock()
  if doc, exists := bucket.collection[key]; exists {
    return error
  } else {
    bucket.collection[key] = bytes.ToLower(doc)
    return nil
  }
}

func (bucket LockBucket) Update(key string, doc []byte) {
  bucket.collection.Lock()
  defer bucket..collection.Unlock()
  bucket.collection[key] = bytes.ToLower(doc)
  return
}

func (bucket *LockBucket) Delete(key string) {
  bucket.Lock()
  defer bucket.Unlock()
  delete(bucket.collection, key)
}


type lock_collection struct {
  sync.RWMutex
  collection map[string][]byte
}

func NewLockCollection() (*lock_collection) {
  return &lock_collection{
    collection: make(map[string][]byte),
  }
}
