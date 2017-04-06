package server

import (
  "net/http"
  "fmt"
  "io"
  "time"
  "regexp"
  "encoding/json"
  "go-cached/database"
)

// Accpeted Paths:
// GET /db/bucket/$bucket_name -> bucket info?
// PUT /db/bucket/$bucket_name -> creates bucket with bucket name + req data
//      IF bucket does not already exist.
// ???
// GET /db/stats -> print each bucket from the database in json
// GET /db/stats/bucket -> bucket - specific stats??
// GET /db/info | GET /db/info/$bucket_name | GET /db/info(?bucket=$bucket_name)
// ???

/*
{
  "bucketId": "Name of Bucket"
  "primayKey": "Field for Primary Key / index"
  "engine": "Name of DB engine to use"
  "secondaryKeys": ["List of fields", "for secodary indexes"]
}
*/

var validDBPath = regexp.MustCompile("^/db/bucket/?([a-zA-Z0-9-_]*)/?$")

func dbBucketAPIHandler(w http.ResponseWriter, r *http.Request) {

  matches := validDBPath.FindStringSubmatch(r.URL.Path)

  if len(matches) < 1 {
    http.NotFound(w, r)
    return
  }
  // Edge cases here: GET or DELETE with no bucketId in uri (len(matches) == 2), ??
  switch r.Method {
  case "PUT":
    bucketPutHandler(w, r)
  case "DELETE":
    bucketDeleteHandler(w, r, matches[1])
  case "GET":
    bucketStatsHandler(w, r, matches[1])
  default:
    http.Error(w, "", http.StatusMethodNotAllowed)
  }
  return
}

func bucketStatsHandler(w http.ResponseWriter, r *http.Request, bucketId string) {
  defer RequestLog(fmt.Sprintf("%s %s %s %s", r.Method, r.URL.Path, r.Proto, r.RemoteAddr), time.Now())
  bmd, err := DocumentDatabase.GetBucketMetaData(bucketId)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError) // should we return a 404? idk
    return
  }
  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(bmd)
}

func bucketPutHandler(w http.ResponseWriter, r *http.Request) {
  defer RequestLog(fmt.Sprintf("%s %s %s %s", r.Method, r.URL.Path, r.Proto, r.RemoteAddr), time.Now())

  var bmd database.BucketMetaData

  if r.ContentLength < 1 {
    http.Error(w, "Request content length 0 or undeterminable", http.StatusBadRequest)
  }
  buf := make([]byte, r.ContentLength)
  _, err := io.ReadFull(r.Body, buf)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  err = json.Unmarshal(buf, &bmd)
  if err != nil {
    Logger.Write("Error decoding request")
    return
  }
  err = DocumentDatabase.AddBucket(bmd)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  return
}

func bucketDeleteHandler(w http.ResponseWriter, r *http.Request, bucketId string) {
  defer RequestLog(fmt.Sprintf("%s %s %s %s", r.Method, r.URL.Path, r.Proto, r.RemoteAddr), time.Now())

  err := DocumentDatabase.DeleteBucket(bucketId)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  return
}
