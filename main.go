package main

import (
  "fugozi/server"
)

func main() {
  srv := server.NewHttpServer()
  srv.SetHttpServerDebug(true)
  srv.RunServer()
}
