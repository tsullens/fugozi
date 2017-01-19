package server

import (
  "net/http"
  "fmt"
  "encoding/json"
  "time"
)

// Route declarations
func rootHandler(w http.ResponseWriter, r *http.Request) {
//  rlog("rootHandler", r)
  defer RequestLog(fmt.Sprintf("%s %s %s %s", r.Method, r.URL.Path, r.Proto, r.RemoteAddr), time.Now())

  http.Redirect(w, r, "/status", http.StatusMovedPermanently)
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
//  rlog("statusHandler", r)
  defer RequestLog(fmt.Sprintf("%s %s %s %s", r.Method, r.URL.Path, r.Proto, r.RemoteAddr), time.Now())

  w.Header().Set("Content-Type", "application/json")
  js, err := json.MarshalIndent(Config, "", "  ")
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
  w.Write(js)
}

/*
  Not used currently

func bucketsHandler(w http.ResponseWriter, r *http.Request) {
  defer RequestLog(fmt.Sprintf("%s %s %s %s", r.Method, r.URL.Path, r.Proto, r.RemoteAddr), time.Now())

  w.Header().Set("Content-Type", "application/json")
  js, err := json.MarshalIndent(buckets, "", "  ")
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
  w.Write(js)
}
*/
