# go-double

[![Go Report Card](https://goreportcard.com/badge/github.com/laurentdutheil/go-double)](https://goreportcard.com/report/github.com/laurentdutheil/go-double) [![PkgGoDev](https://pkg.go.dev/badge/github.com/laurentdutheil/go-double)](https://pkg.go.dev/github.com/laurentdutheil/go-double)

Library that distincts stubs, spies and mocks that depends on [Testify](https://github.com/stretchr/testify)

## Purpose

POC to prepare a possible Pull Request on [Testify](https://github.com/stretchr/testify)

The purpose is to be as retro-compatible as possible

### TODO

- [x] TestData
- [ ] FindClosest
- [ ] FunctionalOptionsArguments (don't understand the behaviour and the implementation => need help)
- [x] NotBefore => NOT DO use InOrder instead
- [x] Maybe => NOT DO use Stub instead
- [x] Unset => NOT DO use Stub instead
- [ ] standardize error message
- [x] Mutex
- [x] Comments for documentation

## Why

If we refer to a definition of [test doubles](https://martinfowler.com/bliki/TestDouble.html), there is an interest in distinguishing stubs, spies and mocks.

**Stubs** provide prepared answers to calls made during the test. For example, a stub needs to return a value in response of a query. If the code under test changes and no longer needs to make that
query, there is no reason that the test should break.

**Spies** are stubs that also record some information based on how they were called. One form of this might be an email service that records how many messages it was sent.

**Mocks** are pre-programmed with expectations which form a specification of the calls they are expected to receive. They can throw an exception if they receive a call they don't expect and are
checked during verification to ensure they got all the calls they were expecting.

## Architecture Decision Record

As I code the POC, I have to choose some design.
I try to document them in ADR:

1. [Separate Stub, Spy and Mock](./adr/01-separate-stub-spy-mock.md)
2. [Propose a generic constructor](./adr/02-generic-constructor.md)
3. [Use reflection to know the number of return arguments](./adr/03-use-reflection-to-know-number-of-return-argument.md)
4. [Don't implement Maybe and Unset methods](./adr/04-dont-implement-Maybe-and-Unset.md)
5. [Encourage to write assertions at the end of the test](./adr/05-encourage-to-write-assertions-at-the-end-of-the-test.md)
6. [Replace NotBefore() by InOrder()](./adr/06-replace-NotBefore-by-InOrder.md)

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
	_, err := s.dependency.DoSomething(number)
	if err != nil {
		return err
	}
	// do something with the result of dependency

	return nil
}

/*
	In test file
*/

type MyStubObject struct {
	double.Stub
}

func (m *MyStubObject) DoSomething(number int) (int, error) {
	args := m.Called(number)
	return args.Int(0), args.Error(1)
}

func TestExample_Stub(t *testing.T) {
	t.Run("with On", func(t *testing.T) {
		stub := double.New[MyStubObject](t)
		stub.On("DoSomething", 3).Return(4, nil)
		// the next line is not called but does not fail the test
		stub.On("DoSomething", 123).Return(-1, fmt.Errorf("stub error"))

		sut := SUT{stub}
		err := sut.MethodToTest(3)

		assert.Nil(t, err)
	})

	t.Run("with When to have typing", func(t *testing.T) {
		stub := double.New[MyStubObject](t)
		stub.When(stub.DoSomething, 3).Return(4, nil)
		// the next line is not called but does not fail the test
		stub.When(stub.DoSomething, 123).Return(-1, fmt.Errorf("stub error"))

		sut := SUT{stub}
		err := sut.MethodToTest(3)

		assert.Nil(t, err)
	})
}
```

### Spy

```go
package mypackage

/*
	In production code
*/

type Dice interface {
	Roll() int
}

type SixDie struct{}

func (SixDie) Roll() int {
	return rand.Intn(6) + 1
}

type Game struct {
	position int
	dice     Dice
}

func (g *Game) Position() int {
	return g.position
}

func (g *Game) Play() {
	g.position += g.dice.Roll()
}

/*
	In test file
*/

func TestExample_Spy(t *testing.T) {
	t.Run("as a Stub", func(t *testing.T) {
		spy := double.New[SpyAsStub](t)
		game := Game{position: 12, dice: spy}
		spy.On("Roll").Return(4)

		game.Play()

		// check the state
		assert.Equal(t, 16, game.Position())
		// and/or check the call
		assert.Equal(t, 1, spy.NumberOfCalls("Roll"))
	})

	t.Run("as a spy of the real implementation", func(t *testing.T) {
		spy := double.New[SpyRealDice](t)
		spy.spied = SixDie{}
		game := Game{position: 12, dice: spy}

		game.Play()

		// check that it is a six die
		assert.GreaterOrEqual(t, game.Position(), 12+1)
		assert.LessOrEqual(t, game.Position(), 12+6)
		// and/or check the call
		assert.Equal(t, 1, spy.NumberOfCalls("Roll"))
	})
}

type SpyAsStub struct {
	double.Spy
}

func (s *SpyAsStub) Roll() int {
	arguments := s.Called()
	return arguments.Int(0)
}

type SpyRealDice struct {
	double.Spy
	spied SixDie
}

func (s *SpyRealDice) Roll() int {
	s.AddActualCall()
	return s.spied.Roll()
}

```

### Mock

```go
package mypackage

/*
In production code
*/
type Interface interface {
	GetSomething(number int) (int, error)
	DoSomething(number int)
}

type ObjectToTest struct {
	dependency Interface
}

func (o ObjectToTest) MethodToTest(number int) error {
	n, err := o.dependency.GetSomething(number)
	if err != nil {
		return err
	}
	n = (n + 42) % 10
	o.dependency.DoSomething(n)
	return nil
}

/*
In test file
*/
type MyMockObject struct {
	double.Mock
}

func (m *MyMockObject) GetSomething(number int) (int, error) {
	args := m.Called(number)
	return args.Int(0), args.Error(1)
}

func (m *MyMockObject) DoSomething(number int) {
	m.Called(number)
}

func TestExample_Mock(t *testing.T) {
	t.Run("Old fashion way with AssertExpectations", func(t *testing.T) {
		mock := double.New[MyMockObject](t)
		mock.On("GetSomething", 3).Return(4, nil)
		mock.On("DoSomething", 6)

		objectToTest := ObjectToTest{mock}
		_ = objectToTest.MethodToTest(3)

		mock.AssertExpectations(t)
	})

	t.Run("Stub the requests and mock the command", func(t *testing.T) {
		mock := double.New[MyMockObject](t)
		mock.On("GetSomething", 3).Return(4, nil)

		objectToTest := ObjectToTest{mock}
		_ = objectToTest.MethodToTest(3)

		mock.AssertCalled(t, "DoSomething", 6)
	})

	t.Run("with AssertNotCalled", func(t *testing.T) {
		mock := double.New[MyMockObject](t)
		expectedError := fmt.Errorf("mock error")
		mock.On("GetSomething", double.Anything).Return(0, expectedError)

		objectToTest := ObjectToTest{mock}
		err := objectToTest.MethodToTest(3)

		assert.Equal(t, expectedError, err)
		mock.AssertNotCalled(t, "DoSomething", double.Anything)
	})

	t.Run("Assert the order of the calls", func(t *testing.T) {
		mock1 := double.New[MyMockObject](t)
		mock1.On("GetSomething", double.Anything).Return(1, nil)
		mock2 := double.New[MyMockObject](t)
		mock2.On("GetSomething", double.Anything).Return(2, nil)
		inOrder := double.InOrder(mock1, mock2)
		objectToTest1 := ObjectToTest{mock1}
		objectToTest2 := ObjectToTest{mock2}

		_ = objectToTest1.MethodToTest(1)
		_ = objectToTest2.MethodToTest(2)

		inOrder.AssertCalled(t, mock1, "GetSomething", 1)
		inOrder.AssertCalled(t, mock1, "DoSomething", double.AnythingOfType("int"))
		inOrder.AssertCalled(t, mock2, "GetSomething", 2)
		inOrder.AssertCalled(t, mock2, "DoSomething", double.AnythingOfType("int"))
		inOrder.AssertNoMoreExpectations(t)
	})
}
```
