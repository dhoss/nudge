package main
import (
  "nudge"
  "fmt"
  "os"
)

func usage() {
  fmt.Fprintf(os.Stderr, "usage: %s http://url.to/post-to\n", os.Args[0])
  os.Exit(2)
}

func main() {
  if len(os.Args) < 2 {
    usage()
  }
  n := nudge.NewRequest(os.Args[1])
  n.Enqueue()
}

