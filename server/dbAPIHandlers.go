package server

import (
  "net/http"
  "io"
  "fmt"
  "time"
  "encoding/json"
  "go-cached/database"
)

// Accpeted Paths:
// /db/bucket

/*
{
  "bucketId": "Name of Bucket"
  "primayKey": "Field for Primary Key / index"
  "engine": "Name of DB engine to use"
  "secondaryKeys": ["List of fields", "for secodary indexes"]
}
*/

func dbAPIHandler(w http.ResponseWriter, r *http.Request) {
  switch r.Method {
  case "GET":
    bucketGetHandler(w, r)
  case "PUT":
    bucketPutHandler(w, r)
//  case "DELETE":
//    bucketDeleteHandler(w, r)
  default:
    http.Error(w, "", http.StatusMethodNotAllowed)
  }
}

func bucketGetHandler(w http.ResponseWriter, r *http.Request) {
  defer RequestLog(fmt.Sprintf("%s %s %s %s", r.Method, r.URL.Path, r.Proto, r.RemoteAddr), time.Now())

}

func bucketPutHandler(w http.ResponseWriter, r *http.Request) {
  defer RequestLog(fmt.Sprintf("%s %s %s %s", r.Method, r.URL.Path, r.Proto, r.RemoteAddr), time.Now())
  Logger.Write()

  buf := make([]byte, r.ContentLength)
  _, err := io.ReadFull(r.Body, buf)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
  bmd := database.BucketMetaData{}
  if err := json.Unmarshal(buf, &bmd); err != nil {
    if Config.Debug {
      Logger.Write(fmt.Sprintf("Bad data on bucket creation: size: %d, data: [ %s ]", r.ContentLength, string(buf)))
    }
    http.Error(w, err.Error(), http.StatusBadRequest)
  }
  DocumentDatabase.AddBucket(bmd)
}
