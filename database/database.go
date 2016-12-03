package database

import (
  "go-cached/util"
  "sync"
  "fmt"
)

type Database struct {
  sync.RWMutex
  buckets map[string]*Bucket
  logger *util.LumberJack
}

func NewDatabase() (*Database) {
    return &Database{
    buckets: make(map[string]*Bucket),
    logger: util.NewLumberJack(util.Config.DbLog),
  }
}

func (db *Database) Select(bucketId string, docId string) ([]byte) {
  db.RLock()
  bucket, exists := db.buckets[bucketId]
  db.RUnlock()
  if exists {
    doc := bucket.Get(docId)
    if doc != nil {
      db.logger.Write(fmt.Sprintf("Retrieved doc %s", docId))
    } else {
      db.logger.Write(fmt.Sprintf("Doc %s not found", docId))
    }
    return doc
  } else {
    return nil
  }
}

func (db *Database) Insert(bucketId, docId string, doc []byte) {
  var (
    bucket *Bucket
    exists bool
  )
  // Try to get the bucket from our current buckets map, if it doesn't
  // exist, let's create it. At the end of this, bucket should always
  // reference a *Bucket
  db.RLock()
  bucket, exists = db.buckets[bucketId]
  db.RUnlock()
  if !exists {
    db.Lock()
    db.buckets[bucketId] = NewBucket(bucketId)
    bucket = db.buckets[bucketId]
    db.Unlock()
    db.logger.Write(fmt.Sprintf("Created bucket %s", bucketId))
  }

  // Error Handling here?
  // What happens if a write to the bucket's collection fails? Is it possible?
  bucket.Update(docId, doc)
  db.logger.Write(fmt.Sprintf("Added doc %s", docId))
  return
}

/*
func (db *Database) Delete(bucketId, docId string) (error) {

}

/*
  Collect statistics about the bucket
  return as json data in []byte

func BucketStats(bucketId) ([]byte, error) {

}
*/
