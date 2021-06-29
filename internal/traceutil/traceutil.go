package traceutil

import (
	"path/filepath"
	"runtime"
	"strings"
)

var wharfCoreDir string

func init() {
	// file: /some/path/to/repo/wharf-core/internal/traceutil/traceutil.go
	_, file, _, ok := runtime.Caller(0)
	if ok {
		// wharfCoreDir: /some/path/to/repo/wharf-core
		wharfCoreDir = filepath.Dir(filepath.Dir(filepath.Dir(file)))
	} else {
		// no file will start with null char
		wharfCoreDir = "\x00"
	}
}

// CallerFileWithLineNum returns the filename with its direct parent directory,
// as well as the line number.
//
// This is done by traversing the call stack up from the caller of the caller of
// this function and ignoring all paths from inside this repository
// (wharf-core). unless it's also a test file ("*_test.go")
func CallerFileWithLineNum() (string, int) {
	const (
		// start on 2 to disregard this func and caller of this func
		startDepth = 2
		// the max is mostly arbitrary, but we don't want an infinite loop
		maxDepth = 15
	)
	for i := startDepth; i <= maxDepth; i++ {
		_, path, line, ok := runtime.Caller(i)

		if ok && isValidCallerFile(path) {
			return fileAndLastDir(path), line
		}
	}
	return "", 0
}

func isValidCallerFile(path string) bool {
	return strings.HasSuffix(path, "_test.go") || !strings.HasPrefix(path, wharfCoreDir) || strings.HasSuffix(path, "/main.go")
}

func fileAndLastDir(path string) string {
	const unknownDir = "???" + string(filepath.Separator)
	dir, file := filepath.Split(path)
	if dir == "" {
		return unknownDir + file
	}
	return filepath.Join(filepath.Base(dir), file)
}
