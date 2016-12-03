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

func initialize(docDB *database.Database) {
  docDB.Insert("_default", "default", []byte(`
    {
    "Name": "default doc"
    }
    `))
  docDB.Insert("_default", "test1", []byte(`
    {
    "Name": "test1"
    }
    `))
}
