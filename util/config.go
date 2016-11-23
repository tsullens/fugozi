package util

type Configuration struct {
  IpAddress string
  Port string
  HttpLog string
  DbLog string
  Debug bool
}

var Config *Configuration
