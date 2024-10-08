package double_test

import (
	"github.com/stretchr/testify/assert"
	"testing"

	. "github.com/laurentdutheil/go-double/double"
)

func TestRuntimeUtil(t *testing.T) {
	t.Run("GetCallingFunctionName", func(t *testing.T) {
		t.Run("Skip specified number of stack frames", func(t *testing.T) {
			var spiedNumberOfFrames int

			beforeMonkeyPatch := RuntimeCallerFunc
			defer func() { RuntimeCallerFunc = beforeMonkeyPatch }()
			RuntimeCallerFunc = func(skip int) (pc uintptr, file string, line int, ok bool) {
				spiedNumberOfFrames = skip
				return 0, "", 0, true
			}

			specifiedNumberOfFrames := 3
			GetCallingFunctionName(specifiedNumberOfFrames)

			assert.Equal(t, specifiedNumberOfFrames, spiedNumberOfFrames)
		})

		t.Run("Panic on runtime.Caller error", func(t *testing.T) {
			beforeMonkeyPatch := RuntimeCallerFunc
			defer func() { RuntimeCallerFunc = beforeMonkeyPatch }()
			RuntimeCallerFunc = func(skip int) (pc uintptr, file string, line int, ok bool) {
				return 0, "", 0, false
			}

			assert.PanicsWithValue(t, "Couldn't get the caller information", func() { GetCallingFunctionName(2) })
		})

		t.Run("Extract function name", func(t *testing.T) {
			beforeMonkeyPatch := RuntimeFuncForPCNameFunc
			defer func() { RuntimeFuncForPCNameFunc = beforeMonkeyPatch }()
			RuntimeFuncForPCNameFunc = func(pc uintptr) string {
				return "github.com/laurentdutheil/go-double/double_test.(*StubExample).MethodWithReturnArguments"
			}

			assert.Equal(t, "MethodWithReturnArguments", GetCallingFunctionName(2))
		})

		t.Run("Extract function name with GCCGO", func(t *testing.T) {
			beforeMonkeyPatch := RuntimeFuncForPCNameFunc
			defer func() { RuntimeFuncForPCNameFunc = beforeMonkeyPatch }()
			RuntimeFuncForPCNameFunc = func(pc uintptr) string {
				return "github_com_laurentdutheil_go-double_store_double_test.MethodWithReturnArguments.pN39_github_com_laurentdutheil_go-double_store_double_test.StubExample"
			}

			assert.Equal(t, "MethodWithReturnArguments", GetCallingFunctionName(2))
		})
	})

	t.Run("GetFunctionName", func(t *testing.T) {
		t.Run("Get the function name", func(t *testing.T) {
			stub := &StubExample{}

			functionName, _ := GetFunctionName(stub.Method)

			assert.Equal(t, "Method", functionName)
		})

		t.Run("Error if pass anything other than a function ", func(t *testing.T) {

			_, err := GetFunctionName("not a function")

			assert.EqualError(t, err, "\"not a function\" of type 'string' is not a function")
		})
	})

	t.Run("GetCallingMethodInformation", func(t *testing.T) {
		t.Run("Find method name and number of return arguments", func(t *testing.T) {
			stubExample := &StubExample{}
			stubMethodCalled := "MethodWithReturnArguments"

			beforeMonkeyPatch := RuntimeFuncForPCNameFunc
			defer func() { RuntimeFuncForPCNameFunc = beforeMonkeyPatch }()
			RuntimeFuncForPCNameFunc = func(pc uintptr) string {
				return stubMethodCalled
			}

			methodInformation, _ := GetCallingMethodInformation(stubExample)

			assert.Equal(t, stubMethodCalled, methodInformation.Name)
			assert.Equal(t, 2, methodInformation.NumOut)
		})

		t.Run("Return an error if method is private or if method does not exist", func(t *testing.T) {
			stubExample := &StubExample{}
			stubMethodCalled := "privateMethod"

			beforeMonkeyPatch := RuntimeFuncForPCNameFunc
			defer func() { RuntimeFuncForPCNameFunc = beforeMonkeyPatch }()
			RuntimeFuncForPCNameFunc = func(pc uintptr) string {
				return stubMethodCalled
			}

			_, err := GetCallingMethodInformation(stubExample)

			assert.Error(t, err)
			assert.ErrorContains(t, err, "couldn't get the caller method information. 'privateMethod' is private or does not exist")
		})
	})
}
