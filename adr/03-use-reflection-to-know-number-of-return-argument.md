- Date: 2024-07-15
- Decision taken by: Laurent DUTHEIL
- Status: accepted

# Context

On `Called()` or `MethodCalled()` method call, this not necessary to fail if a stubbed method have no return parameters. Even if we don't predefine a call without return parameters, the test should
not fail.

The only way to know the number of return parameters is to use reflection. But we have to have a reference of the caller

# Considered options

1. Option 1:
    - Use the generic constructor to set the caller
    - ✅ **Advantage:**
        - with the caller we can easily do reflection and know how many return parameter on each method
        - the stubbed method don't have to fail the test
    - 🚫 **Disadvantage:**
        - The reflection don't work with private methods
2. Option 2:
    - No reflection
    - ✅ **Advantage:**
        - no impact on constructor or on `Called()` or `MethodCalled()` method call
    - 🚫 **Disadvantage:**
        - have to fail even if method have no return arguments

# Advices

<--Any advices worth mentioning-->

# Decision

Option 1: Use the generic constructor to set the caller

