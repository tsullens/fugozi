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

type PrimaryKeyNotFoundError struct {
  pKey string
  doc *Document
}
func (e *PrimaryKeyNotFoundError) Error() string {
  return fmt.Sprintf("Primary Key field %s not found in document %v", e.pKey, e.doc)
}

type PrimaryKeyTypeValueMismatchError struct {
  doc *Document
}
func (e *PrimaryKeyTypeValueMismatchError) Error() string {
  return fmt.Sprintf("Primary Keys currently MUST be a string type. Document: %v", e.doc)
}

type DocumentNotFoundError struct {
  docId string
}
func (e *DocumentNotFoundError) Error() string {
  return fmt.Sprintf("Document with Primary Key Value %s not found", e.docId)
}
