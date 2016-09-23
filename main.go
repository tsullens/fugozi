package main

import (
  "fugozi/httpserver"
  "log"
)

func main() {
  srv := httpserver.NewHttpServer()
  srv.SetHttpServerDebug(true)
  log.Printf("%v %v", srv.Status, srv.Debug)
  srv.RunServer()
}
