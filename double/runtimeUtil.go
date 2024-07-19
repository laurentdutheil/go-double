package double

import (
	"fmt"
	"reflect"
	"regexp"
	"runtime"
	"strings"
)

// RuntimeCallerFunc function variable to allow Monkey Patching in tests
var RuntimeCallerFunc = runtime.Caller

// RuntimeFuncForPCNameFunc function variable to allow Monkey Patching in tests
var RuntimeFuncForPCNameFunc = func(pc uintptr) string {
	return runtime.FuncForPC(pc).Name()
}

// regex for GCCGO functions
var gccgoRE = regexp.MustCompile(`\.pN\d+_`)

func GetCallingFunctionName(skipFrames int) string {
	pc, _, _, ok := RuntimeCallerFunc(skipFrames)
	if !ok {
		panic("Couldn't get the caller information")
	}
	functionPath := RuntimeFuncForPCNameFunc(pc)

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

func GetCallingMethod(caller interface{}) Method {
	functionName := GetCallingFunctionName(3)
	typeOfCaller := reflect.TypeOf(caller)
	method, ok := typeOfCaller.MethodByName(functionName)
	if !ok {
		panic(fmt.Sprintf("Couldn't get the caller method information. '%s' is private or does not exist.", functionName))
	}
	return Method{functionName, method.Type.NumOut()}
}

type Method struct {
	Name   string
	NumOut int
}
