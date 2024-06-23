package models

import (
	"time"
)

// Solution represents a solution to a puzzle.
type Solution []byte

// Puzzle represents a puzzle with its difficulty, nonce, and deadline.
type Puzzle struct {
	// Difficulty represents the difficulty level of a puzzle. It is an integer
	// value between 0 and 255.
	Difficulty int

	// Nonce is used to store a randomly generated value.
	Nonce []byte

	// Deadline represents the deadline time for a puzzle.
	Deadline time.Time
}
