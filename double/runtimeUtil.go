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

// GetCallingFunctionName get the calling function's name using runtime.Caller function.
// skipFrames is the number of stack frames to ascend, with 0 identifying the caller of Caller.
func GetCallingFunctionName(skipFrames int) string {
	pc, _, _, ok := RuntimeCallerFunc(skipFrames)
	if !ok {
		panic("Couldn't get the caller information")
	}
	functionPath := RuntimeFuncForPCNameFunc(pc)

	return extractFunctionName(functionPath)
}

// GetFunctionName get the name of the method passed in parameter.
// Return an error if the method parameter is not reflect.Func type.
//
//	GetFunctionName(stub.Method)
func GetFunctionName(method interface{}) (string, error) {
	valueOfMethod := reflect.ValueOf(method)
	if valueOfMethod.Type().Kind() != reflect.Func {
		return "", fmt.Errorf("%q of type '%T' is not a function", method, method)
	}

	return extractFunctionName(runtime.FuncForPC(valueOfMethod.Pointer()).Name()), nil
}

// GetCallingMethodInformation get the calling function's information such as name and number of return arguments.
// Use reflection to check if the calling function is part of the caller parameter.
// Return an error if the method is not part of the caller parameter
// or if the method is private.
func GetCallingMethodInformation(caller interface{}) (*MethodInformation, error) {
	functionName := GetCallingFunctionName(4)
	typeOfCaller := reflect.TypeOf(caller)
	method, ok := typeOfCaller.MethodByName(functionName)
	if !ok {
		return nil, fmt.Errorf("couldn't get the caller method information. '%s' is private or does not exist", functionName)
	}
	return &MethodInformation{functionName, method.Type.NumOut()}, nil
}

type MethodInformation struct {
	Name   string
	NumOut int
}

func extractFunctionName(functionPath string) string {
	// Next four lines are required to use GCCGO function naming conventions.
	// For Ex:  github_com_docker_libkv_store_mock.WatchTree.pN39_github_com_docker_libkv_store_mock.Mock
	// uses interface information unlike golang github.com/docker/libkv/store/mock.(*Mock).WatchTree
	// With GCCGO we need to remove interface information starting from pN<dd>.
	if gccgoRE.MatchString(functionPath) {
		functionPath = gccgoRE.Split(functionPath, -1)[0]
	}

	parts := strings.Split(functionPath, ".")
	functionName := parts[len(parts)-1]

	parts = strings.Split(functionName, "-")

	return parts[0]
}
