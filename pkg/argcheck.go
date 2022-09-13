package pkg

import (
	"fmt"
	"os"

	lua "github.com/yuin/gopher-lua"
)

func Argcheck(fnName string, nExpected int, l *lua.LState) {
	nObserved := l.GetTop()
	if nExpected != nObserved {
		fmt.Printf(
			"\"%s\" expected %d arguments, found %d\n", fnName, nExpected, nObserved,
		)
		os.Exit(1)
	}
}
