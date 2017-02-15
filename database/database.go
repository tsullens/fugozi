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
  db.log.Write(fmt.Sprintf("%s | Transaction Time %s", msg, time.Since(dbTransactionStart)))
}

func NewDatabase(dblog string) (*Database) {
    return &Database{
      buckets: make(map[string]*Bucket),
      log: util.NewLumberJack(dblog),
    }
}


// Bucket-level functions

func (db *Database) GetBucketMetaData(bucketId string) (BucketMetaData, error) {
  var logMsg string
  defer db.transactionLog(logMsg, time.Now())
  db.RLock()
  bucket, exists := db.buckets[bucketId]
  db.RUnlock()
  if !exists {
    err := &BucketNotFoundError{bucketId: bucketId}
    logMsg = fmt.Sprintf("Get Bucket MetaData %s failed: %s", bucketId, err.Error())
    return BucketMetaData{}, err
  }
  bmd := bucket.GetMetaData()
  logMsg = fmt.Sprintf("Get Bucket MetaData %s completed")
  return bmd, nil
}

func (db *Database) AddBucket(b BucketMetaData) (error) {
  db.log.Write(fmt.Sprintf("Adding bucket: %s", b.BucketId))
  // Create bucket
  var logMsg string
  defer db.transactionLog(logMsg, time.Now())
  db.Lock()
  defer db.Unlock()
  _, exists := db.buckets[b.BucketId]
  if exists {
    err := &BucketExistsError{bucketId: b.BucketId}
    logMsg = fmt.Sprintf("Add Bucket %s failed: %s", b.BucketId, err.Error())
    return err
  }
  bucket, err := NewBucket(b)
  if err != nil {
    logMsg = fmt.Sprintf("Add Bucket %s failed: %s", b.BucketId, err.Error())
    return err
  }
  db.buckets[b.BucketId] = bucket
  logMsg = fmt.Sprintf("Add Bucket %s completed", b.BucketId)
  return nil
}

func (db *Database) DeleteBucket(bucketId string) (error) {
  db.Lock()
  defer db.Unlock()
  _, exists := db.buckets[bucketId]
  if !exists {
    err := &BucketNotFoundError{bucketId: bucketId}
    db.log.Write(fmt.Sprintf("Delete Bucket %s failed: %s", bucketId, err.Error()))
    return err
  }
  delete(db.buckets, bucketId)
  db.log.Write(fmt.Sprintf("Deleted bucket: %s", bucketId))
  return nil
}

// Document-level Functions
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

func (db *Database) Update(bucketId string, doc Document) (error) {
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
    bucket.Update(doc)
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
