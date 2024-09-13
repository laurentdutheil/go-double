package double

// InOrder allows verification in order.
//
//	firstMock := double.New[First](t)
//	secondMock := double.Now[Second](t)
//	inOrder := InOrder(firstMock, secondMock)
//
//	firstMock.add("was called first")
//	secondMock.add("was called second")
//
//	inOrder.AssertCalled(t, firstMock, "add", "was called first")
//	inOrder.AssertCalled(t, secondMock, "add", "was called second")
//	inOrder.AssertNoMoreExpectations(t)
func InOrder(mocks ...IMock) *InOrderValidator {
	result := &InOrderValidator{mocks: mocks}
	for _, mock := range mocks {
		mock.inOrder(result)
	}
	return result
}

type InOrderValidator struct {
	mocks             []IMock
	actualCalls       []ActualCall
	assertCursor      int
	expectationsCount int
}

// AssertCalled assert that the call of the mock was done in right order
func (i *InOrderValidator) AssertCalled(t TestingT, mock IMock, methodName string, arguments ...interface{}) bool {
	call, callExists := i.popCurrentCall()

	if callExists && mock.AssertCalled(t, methodName, arguments...) &&
		call.matches(t, methodName, arguments) {

		i.expectationsCount++
		return true
	}
	t.Errorf("InOrder: %s with arguments %v is not called in right order (expected %d)", methodName, arguments, i.assertCursor)
	return false
}

// AssertNumberOfCallsWithArguments assert that the number of calls of the mock was done in right order
func (i *InOrderValidator) AssertNumberOfCallsWithArguments(t TestingT, mock IMock, expectedCalls int, methodName string, arguments ...interface{}) bool {
	expected := true
	for c := 0; c < expectedCalls; c++ {
		called := i.AssertCalled(t, mock, methodName, arguments...)
		expected = expected && called
	}
	return expected
}

// AssertNoMoreExpectations assert that all the expectations defined in InOrder are verified
func (i *InOrderValidator) AssertNoMoreExpectations(t TestingT) bool {
	if i.assertCursor != i.expectationsCount {
		t.Errorf("InOrder: there are still expectations to call")
	}
	return i.assertCursor == i.expectationsCount
}

func (i *InOrderValidator) addCall(call ActualCall) {
	i.actualCalls = append(i.actualCalls, call)
}

func (i *InOrderValidator) popCurrentCall() (*ActualCall, bool) {
	var currentCall *ActualCall
	if len(i.actualCalls) > i.assertCursor {
		currentCall = &(i.actualCalls[i.assertCursor])
	}
	i.assertCursor++
	return currentCall, currentCall != nil
}
