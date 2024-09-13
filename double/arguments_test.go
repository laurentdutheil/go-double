package double_test

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"

	. "github.com/laurentdutheil/go-double/double"
)

func TestArguments(t *testing.T) {
	t.Run("Get", func(t *testing.T) {
		t.Run("Return the value of the given argument", func(t *testing.T) {
			var args = Arguments([]interface{}{"hello", 123, true})

			assert.Equal(t, "hello", args.Get(0).(string))
			assert.Equal(t, 123, args.Get(1).(int))
			assert.Equal(t, true, args.Get(2).(bool))
		})

		t.Run("Panic if the argument does not exist", func(t *testing.T) {
			var args = Arguments([]interface{}{"hello"})

			assert.PanicsWithValue(t, "assert: arguments: Cannot call Get(1) because there are 1 argument(s).", func() {
				args.Get(1)
			})
		})
	})

	t.Run("String", func(t *testing.T) {
		t.Run("normal String() method - return a string representation of the args", func(t *testing.T) {
			var args = Arguments([]interface{}{"hello", 123, true})

			assert.Equal(t, `string,int,bool`, args.String())
		})

		t.Run("Return the value of the given argument", func(t *testing.T) {
			var args = Arguments([]interface{}{"hello", 123, true})

			assert.Equal(t, "hello", args.String(0))
		})

		t.Run("Panic if the argument does not exist", func(t *testing.T) {
			var args = Arguments([]interface{}{"hello"})

			assert.PanicsWithValue(t, "assert: arguments: Cannot call Get(1) because there are 1 argument(s).", func() {
				args.String(1)
			})
		})

		t.Run("Panic if the argument is not a string", func(t *testing.T) {
			var args = Arguments([]interface{}{true})

			assert.PanicsWithValue(t, "assert: arguments: String(0) failed because object wasn't correct type: %!s(bool=true)", func() {
				args.String(0)
			})
		})

		t.Run("Panic if we pass more than one index", func(t *testing.T) {
			var args = Arguments([]interface{}{"hello", "world"})

			assert.PanicsWithValue(t, "assert: arguments: Wrong number of arguments passed to String.  Must be 0 or 1, not 2", func() {
				args.String(0, 1)
			})
		})
	})

	t.Run("Int", func(t *testing.T) {
		t.Run("Return the value of the given argument", func(t *testing.T) {
			var args = Arguments([]interface{}{"hello", 123, true})

			assert.Equal(t, 123, args.Int(1))
		})

		t.Run("Panic if the argument does not exist", func(t *testing.T) {
			var args = Arguments([]interface{}{123})

			assert.PanicsWithValue(t, "assert: arguments: Cannot call Get(1) because there are 1 argument(s).", func() {
				args.Int(1)
			})
		})

		t.Run("Panic if the argument is not an int", func(t *testing.T) {
			var args = Arguments([]interface{}{true})

			assert.PanicsWithValue(t, "assert: arguments: Int(0) failed because object wasn't correct type: %!s(bool=true)", func() {
				args.Int(0)
			})
		})
	})

	t.Run("Error", func(t *testing.T) {
		t.Run("Return the value of the given argument", func(t *testing.T) {
			var err = errors.New("an Error")
			var args = Arguments([]interface{}{"string", 123, true, err})

			assert.Equal(t, err, args.Error(3))
		})

		t.Run("Return nil of the given argument is nil", func(t *testing.T) {
			var args = Arguments([]interface{}{"string", 123, true, nil})

			assert.Equal(t, nil, args.Error(3))
		})

		t.Run("Panic if the argument does not exist", func(t *testing.T) {
			var args = Arguments([]interface{}{123})

			assert.PanicsWithValue(t, "assert: arguments: Cannot call Get(1) because there are 1 argument(s).", func() {
				_ = args.Error(1)
			})
		})

		t.Run("Panic if the argument is not an error", func(t *testing.T) {
			var args = Arguments([]interface{}{true})

			assert.PanicsWithValue(t, "assert: arguments: Error(0) failed because object wasn't correct type: %!s(bool=true)", func() {
				_ = args.Error(0)
			})
		})
	})

	t.Run("Bool", func(t *testing.T) {
		t.Run("Return the value of the given argument", func(t *testing.T) {
			var args = Arguments([]interface{}{"hello", 123, true})

			assert.Equal(t, true, args.Bool(2))
		})

		t.Run("Panic if the argument does not exist", func(t *testing.T) {
			var args = Arguments([]interface{}{"hello"})

			assert.PanicsWithValue(t, "assert: arguments: Cannot call Get(1) because there are 1 argument(s).", func() {
				args.Bool(1)
			})
		})

		t.Run("Panic if the argument is not a boolean", func(t *testing.T) {
			var args = Arguments([]interface{}{123})

			assert.PanicsWithValue(t, "assert: arguments: Bool(0) failed because object wasn't correct type: %!s(int=123)", func() {
				args.Bool(0)
			})
		})
	})

	t.Run("Matches", func(t *testing.T) {
		t.Run("t.Helper is called", func(t *testing.T) {
			st := &SpiedTestingT{}

			var args = Arguments{123}
			args.Matches(st, 123)

			assert.True(t, st.helperCalled)
		})

		t.Run("false if length are different", func(t *testing.T) {
			st := &SpiedTestingT{}

			var args = Arguments{123}
			assert.False(t, args.Matches(st, 123, true))
			assert.Contains(t, st.logMessages, "Arguments have not the same size: len(expected) == 1 ; len(actual) == 2")

			args = Arguments{123, true, 1.0}
			assert.False(t, args.Matches(st, 123, true))
			assert.Contains(t, st.logMessages, "Arguments have not the same size: len(expected) == 3 ; len(actual) == 2")
		})

		t.Run("compare primitives", func(t *testing.T) {
			st := &SpiedTestingT{}

			var args = Arguments{123, true, 1.0}

			assert.True(t, args.Matches(st, 123, true, 1.0))
			assert.Contains(t, st.logMessages, "\t0: PASS: (int=123) == (int=123)")
			assert.Contains(t, st.logMessages, "\t1: PASS: (bool=true) == (bool=true)")
			assert.Contains(t, st.logMessages, "\t2: PASS: (float64=1) == (float64=1)")

			assert.False(t, args.Matches(st, 456, true, 1.0))
			assert.Contains(t, st.logMessages, "\t0: FAIL: (int=456) != (int=123)")

			assert.False(t, args.Matches(st, 123, false, 1.0))
			assert.Contains(t, st.logMessages, "\t1: FAIL: (bool=false) != (bool=true)")

			assert.False(t, args.Matches(st, 123, true, 1.2))
			assert.Contains(t, st.logMessages, "\t2: FAIL: (float64=1.2) != (float64=1)")

		})

		t.Run("compare string", func(t *testing.T) {
			st := &SpiedTestingT{}

			var args = Arguments{"String"}

			assert.True(t, args.Matches(st, "String"))
			assert.Contains(t, st.logMessages, "\t0: PASS: (string=String) == (string=String)")

			assert.False(t, args.Matches(st, ""))
			assert.Contains(t, st.logMessages, "\t0: FAIL: (string=) != (string=String)")
		})

		t.Run("compare list", func(t *testing.T) {
			st := &SpiedTestingT{}

			var args = Arguments{[]int{1, 2, 3}}

			assert.True(t, args.Matches(st, []int{1, 2, 3}))
			assert.Contains(t, st.logMessages, "\t0: PASS: ([]int=[1 2 3]) == ([]int=[1 2 3])")

			assert.False(t, args.Matches(st, []int{4, 2, 3}))
			assert.Contains(t, st.logMessages, "\t0: FAIL: ([]int=[4 2 3]) != ([]int=[1 2 3])")
		})

		t.Run("compare with Anything argument", func(t *testing.T) {
			st := &SpiedTestingT{}

			var args = Arguments{1, Anything, 3}

			assert.True(t, args.Matches(st, 1, 2, 3))
			assert.Contains(t, st.logMessages, "\t1: PASS: (int=2) == (string=double.Anything)")

			assert.True(t, args.Matches(st, 1, "String", 3))
			assert.Contains(t, st.logMessages, "\t1: PASS: (string=String) == (string=double.Anything)")

			assert.True(t, args.Matches(st, 1, 4.5, 3))
			assert.Contains(t, st.logMessages, "\t1: PASS: (float64=4.5) == (string=double.Anything)")

			assert.True(t, args.Matches(st, 1, 2, Anything))
			assert.Contains(t, st.logMessages, "\t2: PASS: (string=double.Anything) == (int=3)")
		})

		t.Run("compare with AnythingOfType argument", func(t *testing.T) {
			st := &SpiedTestingT{}

			var args = Arguments{AnythingOfType("int"), AnythingOfType("string"), AnythingOfType("*double_test.ExampleType")}

			assert.True(t, args.Matches(st, 1, "String", &ExampleType{true}))
			assert.Contains(t, st.logMessages, "\t0: PASS: (int=1) matches (double.anythingOfTypeArgument=int)")
			assert.Contains(t, st.logMessages, "\t1: PASS: (string=String) matches (double.anythingOfTypeArgument=string)")
			assert.Contains(t, st.logMessages, "\t2: PASS: (*double_test.ExampleType=&{true}) matches (double.anythingOfTypeArgument=*double_test.ExampleType)")

			assert.True(t, args.Matches(st, 2, "any string", &ExampleType{false}))
			assert.Contains(t, st.logMessages, "\t0: PASS: (int=2) matches (double.anythingOfTypeArgument=int)")
			assert.Contains(t, st.logMessages, "\t1: PASS: (string=any string) matches (double.anythingOfTypeArgument=string)")
			assert.Contains(t, st.logMessages, "\t2: PASS: (*double_test.ExampleType=&{false}) matches (double.anythingOfTypeArgument=*double_test.ExampleType)")

			assert.False(t, args.Matches(st, "any string", "any string", &ExampleType{false}))
			assert.Contains(t, st.logMessages, "\t0: FAIL: (string=any string) doesn't match (double.anythingOfTypeArgument=int)")

			assert.False(t, args.Matches(st, 2, 2, &ExampleType{false}))
			assert.Contains(t, st.logMessages, "\t1: FAIL: (int=2) doesn't match (double.anythingOfTypeArgument=string)")

			assert.False(t, args.Matches(st, 2, "any string", ExampleType{false}))
			assert.Contains(t, st.logMessages, "\t2: FAIL: (double_test.ExampleType={false}) doesn't match (double.anythingOfTypeArgument=*double_test.ExampleType)")
		})

		t.Run("compare IsType argument", func(t *testing.T) {
			st := &SpiedTestingT{}

			var args = Arguments([]interface{}{"string", IsType(0), true})

			assert.True(t, args.Matches(st, "string", 123, true))
			assert.Contains(t, st.logMessages, "\t1: PASS: (int=123) matches (*double.isTypeArgument=int)")

			assert.False(t, args.Matches(st, "string", "string", true))
			assert.Contains(t, st.logMessages, "\t1: FAIL: (string=string) doesn't match (*double.isTypeArgument=int)")
		})

		t.Run("compare Matcher argument", func(t *testing.T) {
			st := &SpiedTestingT{}

			matchFn := func(a int) bool {
				return a == 123
			}
			var args = Arguments{"string", MatchedBy(matchFn), true}

			assert.True(t, args.Matches(st, "string", 123, true))
			assert.Contains(t, st.logMessages, "\t1: PASS: (int=123) matches (double.functionMatcherArgument=func(int) bool)")

			assert.False(t, args.Matches(st, "string", 124, true))
			assert.Contains(t, st.logMessages, "\t1: FAIL: (int=124) doesn't match (double.functionMatcherArgument=func(int) bool)")

			assert.False(t, args.Matches(st, "string", false, true))
			assert.Contains(t, st.logMessages, "\t1: FAIL: (bool=false) doesn't match (double.functionMatcherArgument=func(int) bool)")

			assert.False(t, args.Matches(st, "string", nil, true))
			assert.Contains(t, st.logMessages, "\t1: FAIL: (<nil>=<nil>) doesn't match (double.functionMatcherArgument=func(int) bool)")
		})

		t.Run("compare Matcher argument with type nillable", func(t *testing.T) {
			st := &SpiedTestingT{}

			matchFn := func(a []int) bool {
				return a == nil
			}
			var args = Arguments{"string", MatchedBy(matchFn), true}

			assert.True(t, args.Matches(st, "string", nil, true))
			assert.Contains(t, st.logMessages, "\t1: PASS: (<nil>=<nil>) matches (double.functionMatcherArgument=func([]int) bool)")

			assert.False(t, args.Matches(st, "string", 123, true))
			assert.Contains(t, st.logMessages, "\t1: FAIL: (int=123) doesn't match (double.functionMatcherArgument=func([]int) bool)")

			assert.False(t, args.Matches(st, "string", false, true))
			assert.Contains(t, st.logMessages, "\t1: FAIL: (bool=false) doesn't match (double.functionMatcherArgument=func([]int) bool)")
		})
	})

}
