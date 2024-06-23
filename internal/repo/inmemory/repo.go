package inmemory

import (
	"bytes"
	"context"
	_ "embed"
	"math/rand"
	"sync/atomic"

	"github.com/psyhatter/word-of-wisdom/internal/repo"
)

// Quotes taken from the site https://blog.hubspot.com/sales/famous-quotes
//
//go:embed quotes.txt
var quotes []byte

func NewRepo() *Repo {
	q := bytes.FieldsFunc(quotes, func(r rune) bool { return r == '\n' })
	rand.Shuffle(len(q), func(i, j int) { q[i], q[j] = q[j], q[i] })
	return &Repo{quotes: q}
}

var _ repo.Repo = (*Repo)(nil)

// Repo implements repo.Repo.
type Repo struct {
	counter atomic.Int32
	quotes  [][]byte
}

func (r *Repo) GetRandomQuote(_ context.Context) ([]byte, error) {
	b := r.quotes[r.counter.Add(1)%int32(len(quotes))]
	return append([]byte{}, b...), nil
}
