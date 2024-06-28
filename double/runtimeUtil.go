package double

import (
	"regexp"
	"runtime"
	"strings"
)

// function variables to allow Monkey Patching in tests
var runtimeCallerFunc = runtime.Caller
var runtimeFuncForPCNameFunc = func(pc uintptr) string {
	return runtime.FuncForPC(pc).Name()
}

// regex for GCCGO functions
var gccgoRE = regexp.MustCompile(`\.pN\d+_`)

func getCallingFunctionName() string {
	pc, _, _, ok := runtimeCallerFunc(2)
	if !ok {
		panic("Couldn't get the caller information")
	}
	functionPath := runtimeFuncForPCNameFunc(pc)

	// Next four lines are required to use GCCGO function naming conventions.
	// For Ex:  github_com_docker_libkv_store_mock.WatchTree.pN39_github_com_docker_libkv_store_mock.Mock
	// uses interface information unlike golang github.com/docker/libkv/store/mock.(*Mock).WatchTree
	// With GCCGO we need to remove interface information starting from pN<dd>.
	if gccgoRE.MatchString(functionPath) {
		functionPath = gccgoRE.Split(functionPath, -1)[0]
	}

	parts := strings.Split(functionPath, ".")
	return parts[len(parts)-1]
}