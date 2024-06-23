package client

import (
	"testing"
	"time"

	"github.com/psyhatter/word-of-wisdom/pkg/models"
	"github.com/psyhatter/word-of-wisdom/pkg/puzzle"
)

func TestDDOSProtectedClient_solvePuzzle(t *testing.T) {
	p := models.Puzzle{
		Difficulty: 10,
		Nonce:      []byte("nonce"),
	}

	var c DDOSProtectedClient
	solution, err := c.solvePuzzle(p)
	if err != nil {
		t.Fatal(err)
	}

	err = puzzle.CheckSolution(p, solution)
	if err != nil {
		t.Fatal(err)
	}
}

func BenchmarkDDOSProtectedClient_solvePuzzle(b *testing.B) {
	p := models.Puzzle{
		Difficulty: 10,
		Nonce:      []byte("nonce"),
	}

	var c DDOSProtectedClient

	for i := 0; i < b.N; i++ {
		_, _ = c.solvePuzzle(p)
	}
}

func FuzzDDOSProtectedClient_solvePuzzle(f *testing.F) {
	f.Fuzz(func(t *testing.T, difficulty int, nonce []byte) {
		p := models.Puzzle{
			Difficulty: difficulty % 10,
			Nonce:      nonce,
			Deadline:   time.Now().Add(10 * time.Second),
		}

		var c DDOSProtectedClient

		_, _ = c.solvePuzzle(p)
	})
}
