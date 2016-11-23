package server

import (
  "go-cached/database"
  "go-cached/util"
  "net/http"
  "sync"
  "time"
  "fmt"
  "strings"
  "os"
)

var (
  self *httpServer
  buckets = struct {
    sync.RWMutex
    m map[string]*database.Bucket
  }{m: make(map[string]*database.Bucket)}
)
const (
  timeLayout = "2006-01-02 15:04:05.000 MST"
)

type httpServer struct {
  IpAddr string
  Port string
  Logger *util.LumberJack `json:"-"`
  Status string
  StartTime string
  Debug bool
}

func NewHttpServer() (*httpServer) {
  return &httpServer{
    IpAddr: util.Config.IpAddress,
    Port: util.Config.Port,
    Logger: util.NewLumberJack(util.Config.HttpLog),
    Status: "Initialized",
    Debug: util.Config.Debug,
  }
}

/*
  Simple abstraction to handle writing request times to the Logger (LumberJack)
*/
func RequestLog(msg string, start time.Time) {
  elapsed := time.Since(start)
  self.Logger.Write(fmt.Sprintf("%s %s", msg, elapsed))
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
  binding := []string{srv.IpAddr, srv.Port}

  initialize()

  // Route Handlers
  http.HandleFunc("/status/", statusHandler)

  http.HandleFunc("/bucket/", dbHandler)
  http.HandleFunc("/", rootHandler)

  lgmsg := fmt.Sprintf("Listening on %s", strings.Join(binding, ":"))
  self.Logger.Write(lgmsg)

  // Start the server

  err := http.ListenAndServe(strings.Join(binding, ":"), nil)
  self.Logger.Write(err.Error())
  os.Exit(1)
}
