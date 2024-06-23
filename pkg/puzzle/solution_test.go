package puzzle_test

import (
	"testing"

	"github.com/psyhatter/word-of-wisdom/pkg/models"
	"github.com/psyhatter/word-of-wisdom/pkg/puzzle"
)

func TestCheckSolution(t *testing.T) {
	tests := []struct {
		name      string
		puzzle    models.Puzzle
		solution  models.Solution
		expectErr bool
	}{

		{
			name: "Invalid Solution",
			puzzle: models.Puzzle{
				Difficulty: 10,
				Nonce:      []byte("nonce"),
			},
			solution:  []byte("solution"),
			expectErr: true,
		},
		{
			name: "Difficulty failed",
			puzzle: models.Puzzle{
				Difficulty: 256,
				Nonce:      []byte("nonce"),
			},
			solution:  []byte{0, 0, 0, 0, 0, 208, 165, 206},
			expectErr: true,
		},
		{
			name: "Valid Solution",
			puzzle: models.Puzzle{
				Difficulty: 28,
				Nonce:      []byte("nonce"),
			},
			solution:  []byte{0, 0, 0, 0, 0, 208, 165, 206},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := puzzle.CheckSolution(tt.puzzle, tt.solution)
			if (err != nil) != tt.expectErr {
				t.Errorf("CheckSolution() error = %v, expectErr %v", err, tt.expectErr)
			}
		})
	}
}

func BenchmarkCheckSolution(b *testing.B) {
	p := models.Puzzle{
		Difficulty: 7,
		Nonce:      []byte("nonce"),
	}
	solution := []byte{0, 0, 0, 0, 0, 208, 165, 206}
	for i := 0; i < b.N; i++ {
		_ = puzzle.CheckSolution(p, solution)
	}
}

func FuzzCheckSolution(f *testing.F) {
	f.Fuzz(func(t *testing.T, difficulty int, nonce, solution []byte) {
		p := models.Puzzle{
			Difficulty: difficulty % 256,
			Nonce:      nonce,
		}

		_ = puzzle.CheckSolution(p, solution)
	})
}
