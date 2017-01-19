package server

import (
  "go-cached/database"
  "go-cached/util"
  "net/http"
  "time"
  "fmt"
  "strings"
  "os"
)

var (
  DocumentDatabase *database.Database
  Logger *util.LumberJack
  Config *Configuration
)
const (
  timeLayout = "2006-01-02 15:04:05.000 MST"
)

type Configuration struct {
  IpAddress string
  Port string
  HttpLog string
  DbLog string
  Debug bool
}

/*
  Simple abstraction to handle writing request times to the logger (LumberJack)
*/
func RequestLog(msg string, start time.Time) {
  elapsed := time.Since(start)
  Logger.Write(fmt.Sprintf("%s %s", msg, elapsed))
}

func RunServer(config *Configuration) {
  Config = config
  DocumentDatabase = database.NewDatabase(Config.DbLog)
  Logger = util.NewLumberJack(Config.HttpLog)
  binding := []string{Config.IpAddress, Config.Port}
  //srv.Status = "Running"
  //srv.StartTime = time.Now().Format(timeLayout)
  //self = srv
  //initialize(DocumentDatabase)

  // Route Handlers
  http.HandleFunc("/status/", statusHandler)
  http.HandleFunc("/bucket/", docAPIHandler)
  http.HandleFunc("/db/bucket/", dbAPIHandler)
//  http.HandleFunc("/", rootHandler)

  lgmsg := fmt.Sprintf("Listening on %s", strings.Join(binding, ":"))
  Logger.Write(lgmsg)

  // Start the server

  err := http.ListenAndServe(strings.Join(binding, ":"), nil)
  Logger.Write(err.Error())
  os.Exit(1)
}
