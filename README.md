# go-double

Library that distincts stubs, spies and mocks that depends on [Testify](https://github.com/stretchr/testify)

## Purpose

POC to prepare a possible Pull Request on [Testify](https://github.com/stretchr/testify)

## Why

If we refer to a definition of [test doubles](https://martinfowler.com/bliki/TestDouble.html), there is an interest in distinguishing stubs, spies and mocks.

**Stubs** provide canned answers to calls made during the test. For example, a stub needs to return a value in response of a query. If the code under test changes and no longer needs to make that
query, there is no reason that the test should break.

**Spies** are stubs that also record some information based on how they were called. One form of this might be an email service that records how many messages it was sent.

**Mocks** are pre-programmed with expectations which form a specification of the calls they are expected to receive. They can throw an exception if they receive a call they don't expect and are
checked during verification to ensure they got all the calls they were expecting.

## Examples

### Stub

```go
package mypackage

/*
	In production code
*/

type InterfaceExample interface {
	DoSomething(number int) (int, error)
}

// SUT System Under Test
type SUT struct {
	dependency InterfaceExample
}

func (s SUT) MethodToTest(number int) error {
	n, err := s.dependency.DoSomething(number)
	if err != nil {
		return err
	}
	// do something with the result of dependency
	n = n % 100

	return nil
}

/*
	In test file
*/

type MyStubObject struct {
	double.Stub[MyStubObject]
}

func (m *MyStubObject) DoSomething(number int) (int, error) {
	args := m.Called(number)
	return args.Int(0), args.Error(1)
}

func TestStub(t *testing.T) {
	stub := double.New[MyStubObject](t)
	stub.On("DoSomething", 3).Return(4, nil)
	// the next line is not called but does not fail the test
	stub.On("DoSomething", 123).Return(-1, fmt.Errorf("stub error"))

	sut := SUT{stub}
	err := sut.MethodToTest(3)

	assert.Nil(t, err)
}
```