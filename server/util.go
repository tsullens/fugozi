package server

import (
  "fugozi/database"
  "net/http"
  "log"
  "time"
  "fmt"
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

/*
  Simple abstraction to handle writing request times to the Logger (LumberJack)
*/
func RequestLog(msg string, start time.Time) {
  elapsed := time.Since(start)
  self.Logger.Write(fmt.Sprintf("%s %s", msg, elapsed))
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
