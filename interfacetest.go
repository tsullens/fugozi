package main

import (
	"errors"
	"fmt"
	"sync"
)

type Bucket interface {
	Get(string) (string, error)
	Put(string, string) (error)
}

type LockBucket struct {
	sync.RWMutex
	collection map[string]string
}

func (b *LockBucket) Get(key string) (string, error) {
	b.RLock()
	defer b.RUnlock()
	if item, exists := b.collection[key]; exists {
		return item, nil
	} else {
		return "", errors.New(fmt.Sprintf("%s not found", key))
	}
}

func (b *LockBucket) Put(key, item string) (error) {
	b.Lock()
	defer b.Unlock()
	if _, exists := b.collection[key]; exists {
		return errors.New("item exists")
	} else {
		b.collection[key] = item
		return nil
	}
}

type GenericBucket struct {
	collection map[string]string
}

func (b *GenericBucket) Get(key string) (string, error) {
	if item, exists := b.collection[key]; exists {
		return item, nil
	} else {
		return "", errors.New(fmt.Sprintf("%s not found", key))
	}
}

func (b *GenericBucket) Put(key, item string) (error) {
	if _, exists := b.collection[key]; exists {
		return errors.New("item exists")
	} else {
		b.collection[key] = item
		return nil
	}
}

func main() {
	fmt.Println("Hello, playground")
	buckets := make([]Bucket, 0)
	buckets = append(buckets, &LockBucket{collection: make(map[string]string)})
	buckets = append(buckets, &GenericBucket{collection: make(map[string]string)})

	items := map[string]string{
		"test1": "one",
		"test2": "two",
		"test3": "three",
	}
	for key, item := range items {
		for _, bucket := range buckets {
			err := bucket.Put(key, item)
      if err != nil {
        fmt.Println(err)
      }
		}
	}
	for index, bucket := range buckets {
		for key, _ := range items {
      item, err := bucket.Get(key)
      if err != nil {
        fmt.Println(err)
      } else {
        fmt.Printf("%v: %v -> %v\n", index, key, item)
      }
		}
	}
}
