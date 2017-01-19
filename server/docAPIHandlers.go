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

func docAPIHandler(w http.ResponseWriter, r *http.Request) {
  switch r.Method {
  case "GET":
    docGetHandler(w, r)
  case "PUT":
    docPutHandler(w, r)
//  case "DELETE":
//    docDeleteHandler(w, r)
  default:
    http.Error(w, "", http.StatusMethodNotAllowed)
  }
}

func docGetHandler(w http.ResponseWriter, r *http.Request) {
  defer RequestLog(fmt.Sprintf("%s %s %s %s", r.Method, r.URL.Path, r.Proto, r.RemoteAddr), time.Now())
  // try to get a match on our endpoint
  // if no match is found, it's a bad request
  m := validPath.FindStringSubmatch(r.URL.Path)
  if m == nil {
    http.NotFound(w, r)
  }
  doc, err := DocumentDatabase.Select(m[2], m[3])
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
func docPutHandler(w http.ResponseWriter, r *http.Request) {
  defer RequestLog(fmt.Sprintf("%s %s %s %s", r.Method, r.URL.Path, r.Proto, r.RemoteAddr), time.Now())

  m := validPath.FindStringSubmatch(r.URL.Path)
  if m == nil {
    http.NotFound(w, r)
  }
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
  if err := DocumentDatabase.Update(m[2], m[3], buf); err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
}

/*
func dbDeleteHandler(w http.ResponseWriter, r *http.Request) {
  defer RequestLog(fmt.Sprintf("%s %s %s %s %s", r.Method, r.URL.Path, w.Header().Get("status"), r.Proto, r.RemoteAddr), time.Now())

}
*/
