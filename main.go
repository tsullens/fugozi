package main

import (
  "go-cached/server"
  "go-cached/util"
  "flag"
  "regexp"
  "fmt"
  "os"
)

var ipAddress string
var port string
var httpLog string
var dbLog string
var debug bool

func main() {

  flag.StringVar(&ipAddress, "ip", "0.0.0.0", "IP address to listen on")
  flag.StringVar(&port, "p", "3341", "Port to listen on")
  flag.StringVar(&httpLog, "http-log", "httpsrv.log", "Log file name to write Http logs")
  flag.StringVar(&dbLog, "db-log", "db.log", "Log file name to write database logs")
  flag.BoolVar(&debug, "debug", false, "Debug on / off (=true | =false)")
  flag.Parse()

  validIP := regexp.MustCompile("(^[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}$)")
  validPort := regexp.MustCompile("(^[0-9]{2,5}$)")

  if !validIP.MatchString(ipAddress) || !validPort.MatchString(port) {
    fmt.Sprintf("IP and/or Port parameter is invalid: using ip %s port %s", ipAddress, port)
    os.Exit(1)
  }

  util.Config = &util.Configuration{
    IpAddress: ipAddress,
    Port: port,
    HttpLog: httpLog,
    DbLog: dbLog,
    Debug: debug,
  }

  srv := server.NewHttpServer()
  srv.RunServer()
}
