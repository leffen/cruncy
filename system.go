package cruncy

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
)

// WriteHeapFile dumps memory allocation to disc for analytics
func WriteHeapFile(fileName string) error {
	runtime.GC()

	EnsureFileSave(fileName)

	tfile, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("Unable to make perf file with error %s", err)

	}
	pprof.WriteHeapProfile(tfile)
	tfile.Close()
	return nil
}
