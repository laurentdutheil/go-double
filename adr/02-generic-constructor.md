- Date: 2024-07-12
- Decision taken by: Laurent DUTHEIL
- Status: accepted

# Context

As the `Called()` and `MethodCalled()` methods are used in the implementation of the Mock, we cannot pass the testing.T to these method.

Then, if the user don't set the testing.T with the `Test()` method, the `Called()` and `MethodCalled()` have to panic instead of failing the test.

The consequence is that all the remain tests are not executed.

Questions:

- Do we have to force developers to set the testing.T attribute?
- If yes, how do we do it?

# Considered options

1. Option 1:
    - Propose a generic constructor to set the testing.T attribute
    - Panic in `Called()` and `MethodCalled()` if the attribute is not set with an error message that explain how to use the constructor
    - ✅ **Advantage:**
        - The testing.T is always set and the test don't panic anymore
        - We can use the `t.Errorf()`, `t.FailNow()` or `t.Logf()` without failover in all the other methods of the mock
    - 🚫 **Disadvantage:**
        - The developer use to create Mocks with `Mock{}`
        - Test panics if the developer don't use constructor
2. Option 2:
    - Propose a generic constructor to set the testing.T attribute.
    - Log in `Called()` and `MethodCalled()` if the attribute is not set with an error message that explain how to use the constructor
    - Propose a failover if the attribute is not set
    - ✅ **Advantage:**
        - The developer can create Mocks with `Mock{}`
        - The failover can be used in all the other methods of the mock
    - 🚫 **Disadvantage:**
        - The failover have to panic in some case if the testing.T attribute is not set

# Decision

Option 1:

- Propose a generic constructor to set the testing.T attribute
- Panic in `Called()` and `MethodCalled()` if the attribute is not set with an error message that explain how to use the constructor

```go
// New is a constructor for Stub, Spy and Mock
//
//	type MockExample struct {
//		Mock
//	}
//	...
//	func TestExample(t *testing.T) {
//		myMock := New[MockExample](t)
//		...
//	}
func New[T any](t TestingT) *T {
var result interface{} = new(T)
tester := result.(Tester)
tester.Test(t)
return result.(*T)
}

type Tester interface {
Test(t TestingT)
}
```

# Consequences

No failover. We can use testing.T everywhere 