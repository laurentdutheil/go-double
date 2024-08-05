package double

type MocksInOrder struct {
	mocks             []IMock
	actualCalls       []ActualCall
	assertCursor      int
	expectationsCount int
}

func InOrder(mocks ...IMock) *MocksInOrder {
	result := &MocksInOrder{mocks: mocks}
	for _, mock := range mocks {
		mock.InOrder(result)
	}
	return result
}

func (i *MocksInOrder) AssertNumberOfCallsWithArguments(t TestingT, mock IMock, expectedCalls int, methodName string, arguments ...interface{}) bool {
	expected := true
	for c := 0; c < expectedCalls; c++ {
		called := i.AssertCalled(t, mock, methodName, arguments...)
		expected = expected && called
	}
	return expected
}

func (i *MocksInOrder) AssertCalled(t TestingT, mock IMock, methodName string, arguments ...interface{}) bool {
	call, callExists := i.popCurrentCall()

	if callExists && mock.AssertCalled(t, methodName, arguments...) &&
		call.isEqual(methodName, arguments) {

		i.expectationsCount++
		return true
	}
	t.Errorf("InOrder: %s with arguments %v is not called in right order (expected %d)", methodName, arguments, i.assertCursor)
	return false
}

func (i *MocksInOrder) AssertNoMoreExpectations(t TestingT) bool {
	if i.assertCursor != i.expectationsCount {
		t.Errorf("InOrder: there are still expectations to call")
	}
	return i.assertCursor == i.expectationsCount
}

func (i *MocksInOrder) addCall(call ActualCall) {
	i.actualCalls = append(i.actualCalls, call)
}

func (i *MocksInOrder) popCurrentCall() (*ActualCall, bool) {
	var currentCall *ActualCall
	if len(i.actualCalls) > i.assertCursor {
		currentCall = &(i.actualCalls[i.assertCursor])
	}
	i.assertCursor++
	return currentCall, currentCall != nil
}
