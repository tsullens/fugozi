package server

import (
  "fugozi/database"
  "net/http"
  "io"
  "encoding/json"
  "regexp"
  "fmt"
)

// Accpeted Paths:
// /bucket/$bucket
// /bucket/$bucket/$doc

var validPath = regexp.MustCompile("^/(bucket)/([a-zA-Z0-9_]+)/{0,1}([a-zA-Z0-9_]+){0,1}/{0,1}$")

func dbHandler(w http.ResponseWriter, r *http.Request) {
  switch r.Method {
  case "GET":
    dbGetHandler(w, r)
  case "POST":
    dbPostHandler(w, r)
  }
}

func dbGetHandler(w http.ResponseWriter, r *http.Request) {

  lgmsg := fmt.Sprintf("%s %s %s %s", r.Method, r.URL.Path, r.Proto, r.RemoteAddr)
  self.Logger.WriteLog(lgmsg)

  m := validPath.FindStringSubmatch(r.URL.Path)
  if m == nil {
    http.NotFound(w, r)
    return
  }

  buckets.RLock()
  bucket, exists := buckets.m[m[2]]
  buckets.RUnlock()

  if (self.Debug) {
    lgmsg := fmt.Sprintf("%v %v %v", m, len(m), exists)
    self.Logger.WriteLog(lgmsg)
  }

  // if the bucket key exists, let's serve the content
  if exists {
    // Check if a document was requested (.e.g /buckets/$bucket/$doc)
    // if so, fetch the document from the bucket's DocDB and serve it
    // if the $doc key isn't found, return an error.
    if m[3] != "" {
      doc, err := bucket.Get(m[3])
      if err == nil {
        // $doc key exits, serve the content
        w.Header().Set("Content-Type", "application")
        w.Write(doc)
      } else {
        // $doc key was not found, serve back a 404
        http.NotFound(w, r)
        return
      }
    } else {
        // No $doc key was requested, just serve the $bucket's stats
        w.Header().Set("Content-Type", "application/json")
        js, err := json.MarshalIndent((*bucket), "", "  ")
        if err != nil {
          http.Error(w, err.Error(), http.StatusInternalServerError)
          return
        }
        w.Write(js)
      }
    } else {
      // $bucket wasn't found, return a 404
      http.NotFound(w, r)
      return
    }
}

// if /$bucket exists, insert / update the doc

func dbPostHandler(w http.ResponseWriter, r *http.Request) {

  var bucket *database.Bucket
  var exists bool

  lgmsg := fmt.Sprintf("%s %s %s %s", r.Method, r.URL.Path, r.Proto, r.RemoteAddr)
  self.Logger.WriteLog(lgmsg)

  m := validPath.FindStringSubmatch(r.URL.Path)
  if m == nil {
    http.NotFound(w, r)
    return
  }

  buckets.RLock()
  bucket, exists = buckets.m[m[2]]
  buckets.RUnlock()

  if (self.Debug) {
    lgmsg := fmt.Sprintf("%v %v %v", m, len(m), exists)
    self.Logger.WriteLog(lgmsg)
  }

  // bucket doesn't exist, let's create one and add the doc if /$doc is present
  if !exists {
    buckets.Lock()
    buckets.m[m[2]] = database.NewBucket(m[2])
    buckets.Unlock()
    buckets.RLock()
    bucket, exists = buckets.m[m[2]]
    buckets.RUnlock()
  }
  if !exists {
    http.Error(w, "Bucket could not be added", http.StatusInternalServerError)
    return
  }
  // if we have a doc to insert / update (e.g. r.Path = /$bucket/$doc)
  if m[3] != "" {
    if (self.Debug) {
      lgmsg := fmt.Sprintf("%v", r.ContentLength)
      self.Logger.WriteLog(lgmsg)
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
    if (self.Debug) {
      lgmsg := fmt.Sprintf("%s\n", buf)
      self.Logger.WriteLog(lgmsg)
    }
    bucket.Update(m[3], buf)
  }
}
