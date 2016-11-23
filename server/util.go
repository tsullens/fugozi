package server

import (
  "go-cached/database"
  "net/http"
  "log"
)

/*
  Deprecated logging utility
*/
func rlog(message string, r *http.Request) {
  if r != nil {
    log.Printf("%s %s %s %s %s", r.Method, r.URL.Path, r.Proto, r.RemoteAddr, message)
    return
  }
  log.Printf("%s", message)
}

func initialize() {
  test := database.NewBucket("_default")
  test.Update("default", []byte(`
    {
    "Name": "default doc"
    }
    `))
  test.Update("test1", []byte(`
    {
    "Name": "test1"
    }
    `))
  buckets.Lock()
  buckets.m[test.Name] = test
  buckets.Unlock()
}
