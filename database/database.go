package database

import (
  "go-cached/util"
  "sync"
  "fmt"
  "time"
)

type Database struct {
  sync.RWMutex
  buckets map[string]*Bucket
  log *util.LumberJack
}

func (db *Database) transactionLog(msg string, dbTransactionStart time.Time) {
  // Format here is generally for a message to be thus:
  // $bucketId <METHOD> $docId [bucketTransactionTime] | [totalTransactionTime]
  db.log.Write(fmt.Sprintf("%s | Total Transaction %s", msg, time.Since(dbTransactionStart)))
}

func NewDatabase(dblog string) (*Database) {
    return &Database{
      buckets: make(map[string]*Bucket),
      log: util.NewLumberJack(dblog),
    }
}

func (db *Database) AddBucket(b BucketMetaData) (error) {
  // Create bucket
  db.Lock()
  defer db.Unlock()
  _, exists := db.buckets[b.bucketId]
  if exists {
    return &BucketExistsError{bucketId: b.bucketId}
  }
  bucket, err := NewBucket(b)
  if err != nil {
    return err
  }
  db.buckets[b.bucketId] = bucket
  return nil
}

func (db *Database) Select(bucketId string, docId string) ([]byte, error) {
  var msg string
  defer db.transactionLog(msg, time.Now())

  db.RLock()
  bucket, exists := db.buckets[bucketId]
  db.RUnlock()
  if !exists {
    return nil, &BucketNotFoundError{bucketId: bucketId}
  } else {
    bucketTransactionStart := time.Now()
    doc := bucket.Get(docId)
    bucketTransactionTime := time.Since(bucketTransactionStart)
    if doc != nil {
      msg = fmt.Sprintf("%s SELECT %s %s", bucketId, docId, bucketTransactionTime)
    } else {
      msg = fmt.Sprintf("%s SELECT MISS %s %s", bucketId, docId, bucketTransactionTime)
    }
    return doc, nil
  }
}

func (db *Database) Update(bucketId, docId string, doc []byte) (error) {
  var (
    bucket *Bucket
    exists bool
    msg string
  )
  defer db.transactionLog(msg, time.Now())
  db.RLock()
  bucket, exists = db.buckets[bucketId]
  db.RUnlock()
  if !exists {
    return &BucketNotFoundError{bucketId: bucketId}
  } else {
    bucketTransactionStart := time.Now()
    bucket.Update(docId, doc)
    bucketTransactionTime := time.Since(bucketTransactionStart)
    msg = fmt.Sprintf("%s INSERT %s %s", bucketId, docId, bucketTransactionTime)
    return nil
  }
}


func (db *Database) Delete(bucketId, docId string) (error) {
  var msg string
  defer db.transactionLog(msg, time.Now())

  db.RLock()
  bucket, exists := db.buckets[bucketId]
  db.RUnlock()
  if !exists {
    return &BucketNotFoundError{bucketId: bucketId}
  } else {
    bucketTransactionStart := time.Now()
    bucket.Delete(docId)
    bucketTransactionTime := time.Since(bucketTransactionStart)
    msg = fmt.Sprintf("%s DELETE %s %s", bucketId, docId, bucketTransactionTime)
    return nil
  }

}

/*
  Collect statistics about the bucket
  return as json data in []byte

func BucketStats(bucketId) ([]byte, error) {

}
*/
