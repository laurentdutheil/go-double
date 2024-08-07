- Date: 2024-08-05
- Decision taken by: Laurent DUTHEIL
- Status: proposed

# Context

`NotBefore()` method is based on prepared method of the Stub.
But the assertion responsibility is on the Mock.

Here an proposition of another design to assert the order of mocks calls.

```go
   mock1 := double.New[MyMockObject](t)
mock2 := double.New[MyMockObject](t)
inOrder := double.InOrder(mock1, mock2)

mock1.DoSomething(1)
mock2.DoSomething(2)

inOrder.AssertCalled(t, mock1, "DoSomething", 1)
inOrder.AssertCalled(t, mock2, "DoSomething", 2)
inOrder.AssertNoMoreExpectations(t)
```

# Considered options

1. Option 1:
    - Implement the InOrder design
    - ✅ **Advantage:**
        - all code is in the Mock structure
        - InOrder is an assertion on the actual calls
        - easier to read
    - 🚫 **Disadvantage:**
        - break the retro-compatibility
        - explain how to migration from NotBefore to InOrder

# Decision

Option 1: Implement the InOrder design
