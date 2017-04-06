package database

import (
  "sync"
)

type syncmap_collector struct {
  sync.RWMutex
  collection map[string]*Document
}

func NewSyncMapCollector() (*syncmap_collector) {
  return &syncmap_collector{
    collection: make(map[string]*Document),
  }
}

func (c *syncmap_collector) Get(docId string) (*Document) {
  c.RLock()
  defer c.RUnlock()
  if doc, exists := c.collection[docId]; exists {
    return doc
  }
  return nil
}

func (c *syncmap_collector) Update(key string, doc *Document) {
  c.Lock()
  defer c.Unlock()
  c.collection[key] = doc
  return
}

func (c *syncmap_collector) Delete(key string) {
  c.Lock()
  defer c.Unlock()
  delete(c.collection, key)
}
