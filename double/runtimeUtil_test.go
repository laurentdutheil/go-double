package double

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetCallingFunctionName_skip2StackFramesOnCallerCall(t *testing.T) {
	beforeMonkeyPatch := runtimeCallerFunc
	defer func() { runtimeCallerFunc = beforeMonkeyPatch }()
	skipSpy := 0
	runtimeCallerFunc = func(skip int) (pc uintptr, file string, line int, ok bool) {
		skipSpy = skip
		return 0, "", 0, true
	}

	getCallingFunctionName()

	assert.Equal(t, 2, skipSpy)
}

func TestGetCallingFunctionName_panicOnCallerError(t *testing.T) {
	beforeMonkeyPatch := runtimeCallerFunc
	defer func() { runtimeCallerFunc = beforeMonkeyPatch }()
	runtimeCallerFunc = func(skip int) (pc uintptr, file string, line int, ok bool) {
		assert.Equal(t, 2, skip)
		return 0, "", 0, false
	}

	assert.PanicsWithValue(t, "Couldn't get the caller information", func() { getCallingFunctionName() })
}

func TestGetCallingFunctionName_extractFunctionName(t *testing.T) {
	beforeMonkeyPatch := runtimeFuncForPCNameFunc
	defer func() { runtimeFuncForPCNameFunc = beforeMonkeyPatch }()
	runtimeFuncForPCNameFunc = func(pc uintptr) string {
		return "github.com/docker/libkv/store/mock.(*Mock).WatchTree"
	}

	assert.Equal(t, "WatchTree", getCallingFunctionName())
}

func TestGetCallingFunctionName_extractFunctionNameWithGCCGO(t *testing.T) {
	beforeMonkeyPatch := runtimeFuncForPCNameFunc
	defer func() { runtimeFuncForPCNameFunc = beforeMonkeyPatch }()
	runtimeFuncForPCNameFunc = func(pc uintptr) string {
		return "github_com_docker_libkv_store_mock.WatchTree.pN39_github_com_docker_libkv_store_mock.Mock"
	}

	assert.Equal(t, "WatchTree", getCallingFunctionName())
}
