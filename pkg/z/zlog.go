package z

import (
    "bufio"
    "bytes"
    "fmt"
    "os"
    "runtime"
    "strings"
)

var Log *bufio.Writer

func init() {
	f, err := os.Create("zlog_drouet.txt")
	if err != nil {
		panic(err)
	}
	Log = bufio.NewWriter(f)
}

func RawStack() string {
    buf := make([]byte, 10000)
    stackSize := runtime.Stack(buf, false)
    return string(buf[0:stackSize])
}

// This returns a pretty stack, with each line indented by two spaces
func Stack() string {
    // Get a stack crawl as a slice of lines
    buf := make([]byte, 10000)
    stackSize := runtime.Stack(buf, false)
    lines := bytes.Split(buf[0:stackSize], []byte("\n"))

    // Run through the lines in pairs, picking out funcname and sourcepath
    // We skip the first line because it's the name of the goroutine, and
    // we skip the next two lines because that's us.
    lines = lines[3:]

    var stack []string
    var L1, L2 []byte
    var prefix int = -1
    for i, _ := range lines {

        // Get two lines - trim leading whitespace from the second
        if L1 == nil {
            L1 = lines[i]
            continue
        }
        L2 = bytes.TrimLeft(lines[i], " \t")

        // Find the longest prefix of L1 in L2 - this will let us
        // skip the local hard disk path (this presumes GOPATH structure)
        n := 1
        idx := 0
        for n < len(L1) {
            q := bytes.Index(L2, L1[0:n])
            if q == -1 {
                break
            }
            idx = q
            n += 1
        }
        if prefix < 0 {
        	prefix = idx
        }

        // Pick out the funcname
        b1 := bytes.LastIndex(L1[0:n], []byte("/"))
        e1 := bytes.LastIndex(L1, []byte("("))
        funcname := L1[b1+1:e1]

        // Pick out the sourcepath
        e2 := bytes.LastIndex(L2, []byte(" "))
        sourcepath := L2[prefix:e2]

        // Pick out the address offset
       	offset := L2[e2+1:]

        // Add a line with a little indent
        stack = append(stack, fmt.Sprintf("  %s (%s %s)\n", sourcepath, funcname, offset))

        L1 = nil
    }

    // Return as a single string
    return strings.Join(stack, "")
}


// Intended to be the target of a defer in program main
func Shutdown() {
	Log.Flush()
	Log = nil
}
