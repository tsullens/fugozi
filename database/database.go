package database

import (
  "go-cached/util"
  "sync"
  "fmt"
  "time"
)

type Database struct {
  sync.RWMutex
  buckets map[string]Bucket
  *util.LumberJack
}

func transactionLog(msg string, dbTransactionStart time.Time) {
  // Format here is generally for a message to be thus:
  // $bucketId <METHOD> $docId [bucketTransactionTime] | [totalTransactionTime]
  db.Write(fmt.Sprintf("%s | Total Transaction %s", msg, time.Since(dbTransactionStart)))
}

func NewDatabase() (*Database) {
    return &Database{
    buckets: make(map[string]Bucket),
    util.NewLumberJack(util.Config.DbLog),
  }
}

func (db *Database) Select(bucketId string, docId string) ([]byte) {
  var msg string
  defer transactionLog(msg, time.Now())

  db.RLock()
  bucket, exists := db.buckets[bucketId]
  db.RUnlock()
  if exists {
    bucketTransactionStart := time.Now()
    doc := bucket.Get(docId)
    bucketTransactionTime := time.Since(bucketTransactionStart)
    if doc != nil {
      msg = fmt.Sprintf("%s SELECT %s %s", bucketId, docId, bucketTransactionTime)
    } else {
      msg = fmt.Sprintf("%s SELECT MISS %s %s", bucketId, docId, bucketTransactionTime)
    }
    return doc
  } else {
    return nil
  }
}

func (db *Database) Insert(bucketId, docId string, doc []byte) {
  var (
    bucket Bucket
    exists bool
    msg string
  )
  defer transactionLog(msg, time.Now())
  // Try to get the bucket from our current buckets map, if it doesn't
  // exist, let's create it. At the end of this, bucket should always
  // reference a *Bucket
  db.RLock()
  bucket, exists = db.buckets[bucketId]
  db.RUnlock()
  if !exists {
    db.Lock()
    db.buckets[bucketId] = NewLockBucket(bucketId)
    bucket = db.buckets[bucketId]
    db.Unlock()
  }

  // Error Handling here?
  // What happens if a write to the bucket's collection fails? Is it possible?
  bucketTransactionStart := time.Now()
  bucket.Update(docId, doc)
  bucketTransactionTime := time.Since(bucketTransactionStart)
  msg = fmt.Sprintf("%s INSERT %s %s", bucketId, docId, bucketTransactionTime)
  return
}


func (db *Database) Delete(bucketId, docId string) (error) {
  var msg string
  defer transactionLog(msg, time.Now())

  db.RLock()
  bucket, exists := db.buckets[bucketId]
  db.RUnlock()
  if exists {
    bucketTransactionStart := time.Now()
    doc := bucket.Delete(docId)
    bucketTransactionTime := time.Since(bucketTransactionStart)
    msg = fmt.Sprintf("%s DELETE %s %s", bucketId, docId, bucketTransactionTime)
    return doc
  } else {
    return nil
  }

}

/*
  Collect statistics about the bucket
  return as json data in []byte

func BucketStats(bucketId) ([]byte, error) {

}
*/
