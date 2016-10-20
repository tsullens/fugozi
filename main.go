package main

import (
  "fugozi/server"
  "fugozi/util"
  "flag"
  "regexp"
  "fmt"
  "os"
)

var ipAddress string
var port string
var logFile string
var debug bool

func main() {

  flag.StringVar(&ipAddress, "ip", "0.0.0.0", "IP address to listen on")
  flag.StringVar(&port, "p", "3341", "Port to listen on")
  flag.StringVar(&logFile, "log", "srv.log", "Log file name to write to")
  flag.BoolVar(&debug, "debug", false, "Debug on / off (=true | =false)")
  flag.Parse()

  lggr := util.NewLumberJack(logFile)

  validIP := regexp.MustCompile("(^[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}$)")
  validPort := regexp.MustCompile("(^[0-9]{2,5}$)")

  if !validIP.MatchString(ipAddress) || !validPort.MatchString(port) {
    lggr.Write(fmt.Sprintf("IP and/or Port parameter is invalid: using ip %s port %s", ipAddress, port))
    os.Exit(1)
  }

  srv := server.NewHttpServer(ipAddress, port, lggr, debug)
  srv.RunServer()
}
