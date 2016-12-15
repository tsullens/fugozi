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
  self *httpServer
  docDB *database.Database
)
const (
  timeLayout = "2006-01-02 15:04:05.000 MST"
)

type httpServer struct {
  IpAddr string
  Port string
  *util.LumberJack `json:"-"`
  Status string
  StartTime string
  Debug bool
}

func NewHttpServer() (*httpServer) {
  return &httpServer{
    IpAddr: util.Config.IpAddress,
    Port: util.Config.Port,
    util.NewLumberJack(util.Config.HttpLog),
    Status: "Initialized",
    Debug: util.Config.Debug,
  }
}

/*
  Simple abstraction to handle writing request times to the logger (LumberJack)
*/
func RequestLog(msg string, start time.Time) {
  elapsed := time.Since(start)
  self.Write(fmt.Sprintf("%s %s", msg, elapsed))
}

/*
 * Deprecated - not used, but available *
*/
func (srv *httpServer) SetHttpServerDebug(val bool) {
  srv.Debug = val
}

func (srv *httpServer) RunServer() {

  srv.Status = "Running"
  srv.StartTime = time.Now().Format(timeLayout)
  self = srv
  docDB = database.NewDatabase()
  binding := []string{srv.IpAddr, srv.Port}

  initialize(docDB)

  // Route Handlers
  http.HandleFunc("/status/", statusHandler)

  http.HandleFunc("/bucket/", dbHandler)
  http.HandleFunc("/", rootHandler)

  lgmsg := fmt.Sprintf("Listening on %s", strings.Join(binding, ":"))
  self.Write(lgmsg)

  // Start the server

  err := http.ListenAndServe(strings.Join(binding, ":"), nil)
  self.Write(err.Error())
  os.Exit(1)
}
