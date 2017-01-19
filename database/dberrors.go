package database

import(
  "fmt"
)

type BucketNotFoundError struct {
  bucketId string
}

func (e *BucketNotFoundError) Error() string {
  return fmt.Sprintf("bucket %s not found", e.bucketId)
}

type BucketExistsError struct {
  bucketId string
}

func (e *BucketExistsError) Error() string {
  return fmt.Sprintf("bucket %s already exists", e.bucketId)
}

type BucketEngineError struct {
  engineName string
}

func (e *BucketEngineError) Error() string {
  return fmt.Sprintf("database engine %s does not exist", e.engineName)
}
