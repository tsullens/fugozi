package server

import (
  "net/http"
  "io"
  "regexp"
  "fmt"
  "time"
  "gocached/database"
  "encoding/json"
)

// Accpeted Paths:
// /bucket/$bucket
// /bucket/$bucket/$doc

var validDocPath = regexp.MustCompile("^/bucket/([a-zA-Z0-9_-]+)/?([a-zA-Z0-9_-])*/?")

func docAPIHandler(w http.ResponseWriter, r *http.Request) {
  matches := validDocPath.FindStringSubmatch(r.URL.Path)

  if len(matches) < 2 {
    http.NotFound(w, r)
  }
  switch r.Method {
  case "GET":
    if len(matches) != 3 {
      http.Error(w, "", http.StatusBadRequest)
      return
    } else {
    docGetHandler(w, r, matches[1], matches[2])
    }
  case "POST":
    docPostHandler(w, r, matches[1])
  default:
    http.Error(w, "", http.StatusMethodNotAllowed)
  }
}

func docGetHandler(w http.ResponseWriter, r *http.Request, bucketId, docId string) {
  defer RequestLog(fmt.Sprintf("%s %s %s %s", r.Method, r.URL.Path, r.Proto, r.RemoteAddr), time.Now())

  doc, err := DocumentDatabase.Select(bucketId, docId)
  // if doc is nil, either the bucket doesn't exist,
  // or the docId doesn't exist... either way just return a 404
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  } else {
    // else let's write our doc content to the response
    w.Header().Set("Content-Type", "application")
    w.Write(doc)
  }
}

// if /$bucket exists, insert / update the doc
func docPostHandler(w http.ResponseWriter, r *http.Request, bucketId) {
  defer RequestLog(fmt.Sprintf("%s %s %s %s", r.Method, r.URL.Path, r.Proto, r.RemoteAddr), time.Now())

  var data database.DocumentData

  if (Config.Debug) {
    lgmsg := fmt.Sprintf("%v", r.ContentLength)
    Logger.Write(lgmsg)
  }
  if r.ContentLength < 1 {
    http.Error(w, "Request content length 0 or undeterminable", http.StatusBadRequest)
  }

  buf := make([]byte, r.ContentLength)
  _, err := io.ReadFull(r.Body, buf)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
  err = json.Unmarshal(buf, &data)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  var doc = Document{
    timestamp: time.Now(),
    data: data,
  }
  if err := DocumentDatabase.Update(bucketId, doc); err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
}

/*
func dbDeleteHandler(w http.ResponseWriter, r *http.Request) {
  defer RequestLog(fmt.Sprintf("%s %s %s %s %s", r.Method, r.URL.Path, w.Header().Get("status"), r.Proto, r.RemoteAddr), time.Now())

}
*/
