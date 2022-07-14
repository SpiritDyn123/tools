package log

import (
	"runtime"
	"strings"
)

func getCaller(callDepth int, suffixesToIgnore ...string) (file string, line int) {
	// bump by 1 to ignore the getCaller (this) stackframe
	callDepth++
outer:
	for {
		var ok bool
		_, file, line, ok = runtime.Caller(callDepth)
		if !ok {
			file = "???"
			line = 0
			break
		}

		if strings.Index(file, "game_framework/commons-go") >= 0 && strings.Index(file, "wlog") >= 0 {
			callDepth++
			continue outer
		}

		if strings.Index(file, "logrus") >= 0 {
			for _, s := range suffixesToIgnore {
				if strings.HasSuffix(file, s) {
					callDepth++
					continue outer
				}
			}
		}

		break
	}
	return
}

func getCallerIgnoringLogMulti(callDepth int) (string, int) {
	// the +1 is to ignore this (getCallerIgnoringLogMulti) frame
	//return getCaller(callDepth+1, "hooks.go", "entry.go", "logger.go", "exported.go",
	//	"console.go",
	//	"file.go",
	//	"goline.go",
	//	"kafka.go",
	//	"log.go",
	//	"logrus_impl.go",
	//	"asm_amd64.s")

	return getCaller(callDepth+1, "hooks.go", "entry.go", "logger.go", "exported.go", "asm_amd64.s")
}

