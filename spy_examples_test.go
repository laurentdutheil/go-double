package main

import (
	"github.com/laurentdutheil/go-double/double"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

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
