package database

import (
  "sync"
  "errors"
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
  Name string
  database *lock_collection `json:"-"`
}

func NewBucket(name string) (*Bucket) {
  return &Bucket{
    Name: name,
    database: NewLockCollection(),
  }
}

func (bucket *Bucket) Get(key string) ([]byte, error) {
  bucket.database.RLock()
  defer bucket.database.RUnlock()
  if doc, ok := bucket.database.collection[key]; ok {
    return doc, nil
  }
  return nil, errors.New("Key not found")
}

func (bucket *Bucket) Update(key string, doc []byte) {
  bucket.database.Lock()
  defer bucket.database.Unlock()
  bucket.database.collection[key] = bytes.ToLower(doc)
  return
}

func (bucket *Bucket) Delete(key string) {
  bucket.database.Lock()
  defer bucket.database.Unlock()
  delete(bucket.database.collection, key)
}
