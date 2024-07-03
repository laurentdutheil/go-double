package double_test

import (
	"github.com/stretchr/testify/assert"
	"testing"

	. "github.com/laurentdutheil/go-double/double"
)

func TestGetCallingFunctionName_skipStackFramesOnCallerCall(t *testing.T) {
	beforeMonkeyPatch := RuntimeCallerFunc
	defer func() { RuntimeCallerFunc = beforeMonkeyPatch }()
	skipSpy := 0
	RuntimeCallerFunc = func(skip int) (pc uintptr, file string, line int, ok bool) {
		skipSpy = skip
		return 0, "", 0, true
	}

	GetCallingFunctionName(3)

	assert.Equal(t, 3, skipSpy)
}

func TestGetCallingFunctionName_panicOnCallerError(t *testing.T) {
	beforeMonkeyPatch := RuntimeCallerFunc
	defer func() { RuntimeCallerFunc = beforeMonkeyPatch }()
	RuntimeCallerFunc = func(skip int) (pc uintptr, file string, line int, ok bool) {
		return 0, "", 0, false
	}

	assert.PanicsWithValue(t, "Couldn't get the caller information", func() { GetCallingFunctionName(2) })
}

func TestGetCallingFunctionName_extractFunctionName(t *testing.T) {
	beforeMonkeyPatch := RuntimeFuncForPCNameFunc
	defer func() { RuntimeFuncForPCNameFunc = beforeMonkeyPatch }()
	RuntimeFuncForPCNameFunc = func(pc uintptr) string {
		return "github.com/laurentdutheil/go-double/double_test.(*StubExample).MethodWithReturnArguments"
	}

	assert.Equal(t, "MethodWithReturnArguments", GetCallingFunctionName(2))
}

func TestGetCallingFunctionName_extractFunctionNameWithGCCGO(t *testing.T) {
	beforeMonkeyPatch := RuntimeFuncForPCNameFunc
	defer func() { RuntimeFuncForPCNameFunc = beforeMonkeyPatch }()
	RuntimeFuncForPCNameFunc = func(pc uintptr) string {
		return "github_com_laurentdutheil_go-double_store_double_test.MethodWithReturnArguments.pN39_github_com_laurentdutheil_go-double_store_double_test.StubExample"
	}

	assert.Equal(t, "MethodWithReturnArguments", GetCallingFunctionName(2))
}

func TestGetCallingMethod_(t *testing.T) {
	beforeMonkeyPatch := RuntimeFuncForPCNameFunc
	defer func() { RuntimeFuncForPCNameFunc = beforeMonkeyPatch }()
	RuntimeFuncForPCNameFunc = func(pc uintptr) string {
		return "github.com/laurentdutheil/go-double/double_test.(*StubExample).Method"
	}

	stubExample := &StubExample{}

	method := GetCallingMethod(stubExample)

	assert.Equal(t, "Method", method.Name)
	assert.Equal(t, 0, method.Type.NumOut())
}
