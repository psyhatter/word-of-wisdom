package puzzle

import (
	"crypto/sha256"
	"fmt"
	"hash"
	"math/bits"
	"sync"

	"github.com/psyhatter/word-of-wisdom/pkg/models"
)

// CheckSolution verifies the solution to a puzzle by comparing the number of
// zero bits in the solution hash with the required difficulty level. It returns
// an error if the solution does not meet the required difficulty level.
// If an error is returned, it is always implemented by SolutionErr.
func CheckSolution(p models.Puzzle, solution models.Solution) error {
	h, _ := hashPool.Get().(hash.Hash)
	if h == nil {
		h = sha256.New()
	}
	defer hashPool.Put(h)

	h.Reset()

	h.Write(p.Nonce)
	h.Write(solution)

	b, _ := bufPool.Get().(*[]byte)
	if b == nil {
		b = &[]byte{}
	}
	defer bufPool.Put(b)

	*b = h.Sum((*b)[:0])

	var zeros int
	for i := range *b {
		zeros += bits.LeadingZeros8((*b)[i])
		if zeros >= p.Difficulty {
			return nil
		}
		if (*b)[i] != 0 {
			break
		}
	}

	return SolutionErr{Required: p.Difficulty, Actual: zeros}
}

// SolutionErr represents an error indicating that the solution to a puzzle does
// not meet the required difficulty level.
type SolutionErr struct{ Required, Actual int }

func (s SolutionErr) Error() string {
	return fmt.Sprintf("the solution contains %d zero bits, but should %d", s.Actual, s.Required)
}

var bufPool sync.Pool
var hashPool sync.Pool
