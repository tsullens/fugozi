package server

import (
  "net/http"
  "io"
  "regexp"
  "fmt"
  "time"
)

// Accpeted Paths:
// /bucket/$bucket
// /bucket/$bucket/$doc

var validPath = regexp.MustCompile("^/(bucket)/([a-zA-Z0-9_]+)/([a-zA-Z0-9_]+)/{0,1}$")

func dbHandler(w http.ResponseWriter, r *http.Request) {
  switch r.Method {
  case "GET":
    dbGetHandler(w, r)
  case "PUT":
    dbPutHandler(w, r)
  }
}

func dbGetHandler(w http.ResponseWriter, r *http.Request) {
  defer RequestLog(fmt.Sprintf("%s %s %s %s %s", r.Method, r.URL.Path, w.Header().Get("status"), r.Proto, r.RemoteAddr), time.Now())

  // try to get a match on our endpoint
  // if no match is found, it's a bad request
  m := validPath.FindStringSubmatch(r.URL.Path)
  if m == nil {
    http.NotFound(w, r)
    return
  }

  doc := docDB.Select(m[2], m[3])
  // if doc is nil, either the bucket doesn't exist,
  // or the docId doesn't exist... either way just return a 404
  if doc == nil {
    http.NotFound(w, r)
    return
  } else {
    // else let's write our doc content to the response
    w.Header().Set("Content-Type", "application")
    w.Write(doc)
  }
}

// if /$bucket exists, insert / update the doc

func dbPutHandler(w http.ResponseWriter, r *http.Request) {
  defer RequestLog(fmt.Sprintf("%s %s %s %s %s", r.Method, r.URL.Path, w.Header().Get("status"), r.Proto, r.RemoteAddr), time.Now())

  m := validPath.FindStringSubmatch(r.URL.Path)
  if m == nil {
    http.NotFound(w, r)
    return
  }

  if (self.Debug) {
    lgmsg := fmt.Sprintf("%v", r.ContentLength)
    self.Write(lgmsg)
  }
  if r.ContentLength < 1 {
    http.Error(w, "Request content length 0 or undeterminable", http.StatusBadRequest)
    return
  }
  buf := make([]byte, r.ContentLength)
  _, err := io.ReadFull(r.Body, buf)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  docDB.Insert(m[2], m[3], buf)
  return
}

/*
func dbDeleteHandler(w http.ResponseWriter, r *http.Request) {
  defer RequestLog(fmt.Sprintf("%s %s %s %s %s", r.Method, r.URL.Path, w.Header().Get("status"), r.Proto, r.RemoteAddr), time.Now())

}
*/
