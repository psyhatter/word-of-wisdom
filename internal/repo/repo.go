package repo

import (
	"context"
)

type Repo interface {
	// GetRandomQuote is a method that retrieves a random quote from the repository.
	GetRandomQuote(_ context.Context) ([]byte, error)
}
