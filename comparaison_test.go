package main

import (
	"github.com/laurentdutheil/go-double/double"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

/*
	In production code
*/

type Order struct {
	articleName           string
	requiredNumberOfItems uint
	isFilled              bool
}

type Warehouse interface {
	HasInventory(articleName string, requiredNumberOfItems uint) bool
	Remove(articleName string, requiredNumberOfItems uint)
}

func (o *Order) Fill(w Warehouse) {
	if w.HasInventory(o.articleName, o.requiredNumberOfItems) {
		w.Remove(o.articleName, o.requiredNumberOfItems)
		o.isFilled = true
	}
}

/*
	In test file
*/

func TestExample_Comparaison(t *testing.T) {
	t.Run("with mock.Mock", func(t *testing.T) {
		t.Run("remove article from inventory if enough stock", func(t *testing.T) {
			// Arrange
			warehouseMock := WarehouseMock{}
			warehouseMock.On("HasInventory", "fakeArticle", mock.AnythingOfType("uint")).Return(true)
			// if we remove the next line, the test fail
			warehouseMock.On("Remove", "fakeArticle", uint(50))

			// Act
			order := Order{articleName: "fakeArticle", requiredNumberOfItems: 50}
			order.Fill(&warehouseMock)

			// Assert
			assert.True(t, order.isFilled)
			warehouseMock.AssertNumberOfCalls(t, "Remove", 1)
			warehouseMock.AssertCalled(t, "Remove", "fakeArticle", uint(50))
		})

		t.Run("remove article from inventory if enough stock (with AssertExpectations)", func(t *testing.T) {
			// Arrange
			warehouseMock := WarehouseMock{}
			warehouseMock.On("HasInventory", "fakeArticle", mock.AnythingOfType("uint")).Return(true)
			// The assertion is in the 'arrange' part.
			warehouseMock.On("Remove", "fakeArticle", uint(50)).Once()

			// Act
			order := Order{articleName: "fakeArticle", requiredNumberOfItems: 50}
			order.Fill(&warehouseMock)

			// Assert
			assert.True(t, order.isFilled)
			warehouseMock.AssertExpectations(t)
		})

		t.Run("doesn't remove article from inventory if not enough stock", func(t *testing.T) {
			// Arrange
			warehouseMock := WarehouseMock{}
			warehouseMock.On("HasInventory", "fakeArticle", mock.AnythingOfType("uint")).Return(false)

			// Act
			order := Order{articleName: "fakeArticle", requiredNumberOfItems: 50}
			order.Fill(&warehouseMock)

			// Assert
			assert.False(t, order.isFilled)
			warehouseMock.AssertNotCalled(t, "Remove", mock.Anything, mock.Anything)
			// the AssertExpectations is useless as we want to verify that there is no call of 'remove'
		})
	})

	t.Run("with double.Mock", func(t *testing.T) {
		t.Run("remove article from inventory if enough stock", func(t *testing.T) {
			// Arrange
			warehouseMock := double.New[WarehouseDoubleMock](t)
			warehouseMock.On("HasInventory", "fakeArticle", double.AnythingOfType("uint")).Return(true)

			// Act
			order := Order{articleName: "fakeArticle", requiredNumberOfItems: 50}
			order.Fill(warehouseMock)

			// Assert
			assert.True(t, order.isFilled)
			warehouseMock.AssertNumberOfCalls(t, "Remove", 1)
			warehouseMock.AssertCalled(t, "Remove", "fakeArticle", uint(50))
		})

		t.Run("doesn't remove article from inventory if not enough stock", func(t *testing.T) {
			// Arrange
			warehouseMock := double.New[WarehouseDoubleMock](t)
			warehouseMock.On("HasInventory", "fakeArticle", double.AnythingOfType("uint")).Return(false)

			// Act
			order := Order{articleName: "fakeArticle", requiredNumberOfItems: 50}
			order.Fill(warehouseMock)

			// Assert
			assert.False(t, order.isFilled)
			warehouseMock.AssertNotCalled(t, "Remove", mock.Anything, mock.Anything)
		})
	})
}

type WarehouseMock struct {
	mock.Mock
}

func (w *WarehouseMock) HasInventory(articleName string, requiredNumberOfItems uint) bool {
	arguments := w.Called(articleName, requiredNumberOfItems)
	return arguments.Bool(0)
}

func (w *WarehouseMock) Remove(articleName string, requiredNumberOfItems uint) {
	w.Called(articleName, requiredNumberOfItems)
}

type WarehouseDoubleMock struct {
	double.Mock
}

func (w *WarehouseDoubleMock) HasInventory(articleName string, requiredNumberOfItems uint) bool {
	arguments := w.Called(articleName, requiredNumberOfItems)
	return arguments.Bool(0)
}

func (w *WarehouseDoubleMock) Remove(articleName string, requiredNumberOfItems uint) {
	w.Called(articleName, requiredNumberOfItems)
}
