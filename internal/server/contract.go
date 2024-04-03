package server

import "context"

type quoteProvider interface {
	GetQuote(ctx context.Context) (string, error)
}
