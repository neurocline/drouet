package z

import (
	"bufio"
	"os"
)

var Log *bufio.Writer

func init() {
	f, err := os.Create("zlog_drouet.txt")
	if err != nil {
		panic(err)
	}
	Log = bufio.NewWriter(f)
}

// Intended to be the target of a defer in program main
func Shutdown() {
	Log.Flush()
	Log = nil
}
