- Date: 2024-07-15
- Decision taken by: Laurent DUTHEIL
- Status: proposed

# Context

The [Arrange, Act, Assert](https://xp123.com/3a-arrange-act-assert/) pattern (or Given, When, Then) helps developers to write readable and maintainable tests.

The `AssertExpectations()` method don't respect this pattern as the expectations are in the beginning of the test (with the `On().Return()` lines).

Do we encourage developers to respect this pattern?
If yes, how?

# Considered options

1. Option 1:
    - Don't implement the `AssertExpectations()` method
    - ✅ **Advantage:**
        - developers have no choice
        - the other `Assert*` methods exists
        - developer tend to write less mock assertions making the test more readable
    - 🚫 **Disadvantage:**
        - break the retro-compatibility
        - `AssertExpectations()` method is used a lot
        - developers have to migrate a lot of tests to respect the Arrange, Act, Assert pattern
2. Option 2:
    - implement the `AssertExpectations()` method
    - ✅ **Advantage:**
        - retro-compatibility
        - facilitate the migration and the adoption
    - 🚫 **Disadvantage:**
        - developers continue to use `AssertExpectations()` method and don't respect the Arrange, Act, Assert pattern
        - have to explain/document the interest of the pattern
3. Option 3:
    - implement a deprecated `AssertExpectations()` method
    - ✅ **Advantage:**
        - retro-compatibility
        - facilitate the migration and the adoption
        - developers have time to migrate their tests to respect the Arrange, Act, Assert pattern
    - 🚫 **Disadvantage:**
        - `AssertExpectations()` method is used a lot

# Decision

Option 3: implement the `AssertExpectations()` method.

All the prepared called are considered as expectations.