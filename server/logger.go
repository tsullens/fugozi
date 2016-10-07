package server

import (
  "os"
  "fmt"
  "time"
)

type Logger struct {
  Channel chan string
}

func NewLogger(fpath string) (*Logger) {
  fi, err := os.Create(fpath)
  if err != nil {
    panic(err)
  }
  ch := make(chan string)
  go runLogger(ch, fi)
  return &Logger {
    Channel: ch,
  }
}

func (lggr *Logger) WriteLog(msg string) {
  lggr.Channel <- msg
}

func runLogger(ch <-chan string, fi *os.File) {
  defer fi.Close()
  fmt.Printf("starting logger\n")
  for {
    //fmt.Printf("waiting for messages...\n") // Debugging
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
