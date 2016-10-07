package main

import (
  "fugozi/server"
  "log"
)

func main() {
  srv := server.NewHttpServer()
  srv.SetHttpServerDebug(true)
  log.Printf("%v %v", srv.Status, srv.Debug)
  srv.RunServer()
}
