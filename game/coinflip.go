package game

import (
	"math/rand"
	"time"
)

// CoinSide represents the side of a coin
type CoinSide string

const (
	Heads CoinSide = "HEADS"
	Tails CoinSide = "TAILS"
)

// CoinFlip represents the coin flip game logic
type CoinFlip struct {
	Result CoinSide
}

// NewCoinFlip creates a new coin flip game
func NewCoinFlip() *CoinFlip {
	return &CoinFlip{}
}

// Flip flips the coin and returns the result
func (c *CoinFlip) Flip() CoinSide {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())
	
	// Generate a random number between 0 and 1
	if rand.Intn(2) == 0 {
		c.Result = Heads
	} else {
		c.Result = Tails
	}
	
	return c.Result
}
