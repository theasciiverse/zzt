package main

// Session logging: ZZT_LOG=1 zzt 2>zzt.log writes keys, boards, sound and
// a periodic snapshot to stderr. Off by default, since stderr is the screen.

import (
	"fmt"
	"os"
	"time"
)

var logOn = os.Getenv("ZZT_LOG") != ""
var logStart = time.Now()

func logf(format string, args ...interface{}) {
	if logOn {
		fmt.Fprintf(os.Stderr, "%8.3f %s\n", time.Since(logStart).Seconds(), fmt.Sprintf(format, args...))
	}
}
