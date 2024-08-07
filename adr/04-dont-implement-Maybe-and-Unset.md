- Date: 2024-07-15
- Decision taken by: Laurent DUTHEIL
- Status: accepted

# Context

As:

- the Stub structure handles the prepared calls
- the Mock structure handles the assertions
- a call of a stubbed method don't fail the test (except if we have to return something that was not predefined)
- a stubbed method that is not called don't fail the test

I think we don't need the `Maybe()` and `Unset()` methods

# Considered options

1. Option 1:
    - Don't implement `Maybe()` and `Unset()` methods
    - ✅ **Advantage:**
        - simpler API
        - simple implementation
        - these methods are not used a lot
    - 🚫 **Disadvantage:**
        - break the retro-compatibility
        - have to explain the migration
2. Option 2:
    - Implement `Maybe()` and `Unset()` methods
    - ✅ **Advantage:**
        - don't break the retro-compatibility
    - 🚫 **Disadvantage:**
        - the methods do nothing

# Decision

Option 1: Don't implement `Maybe()` and `Unset()` methods

# Consequences

TODO: write a migration documentation
