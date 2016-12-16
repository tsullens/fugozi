package util

import (
  "os"
  "fmt"
  "time"
)

type LumberJack struct {
  Channel chan string
}

const (
  timeLayout = "2006-01-02 15:04:05.000 MST"
)

func NewLumberJack(fpath string) (*LumberJack) {
  fi, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
  if err != nil {
    panic(err)
  }
  ch := make(chan string, 100)
  go run(ch, fi)
  return &LumberJack {
    Channel: ch,
  }
}

func (lggr *LumberJack) Write(msg string) {
  lggr.Channel <- msg
}

func run(ch <-chan string, fi *os.File) {
  defer fi.Close()

  for {
    select {
    case msg := <-ch:
      fmsg := fmt.Sprintf("%s %s\n", time.Now().Format(timeLayout), msg)
      if _, err :=  fi.WriteString(fmsg); err != nil {
        panic(err)
      }
      fi.Sync()
    }
  }
}
