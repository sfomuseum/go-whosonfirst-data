package iterator

import (
	"context"
	"github.com/aaronland/go-roster"
	"net/url"
)

type Iterator interface {
	IterateWithChannels(context.Context, string, chan Result, chan error, chan bool)
}

type IteratorInitializeFunc func(context.Context, string) (Iterator, error)

var iterators roster.Roster

func ensureIterators() error {

	if iterators == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		iterators = r
	}

	return nil
}

func RegisterIterator(ctx context.Context, scheme string, f IteratorInitializeFunc) error {

	err := ensureIterators()

	if err != nil {
		return err
	}

	return iterators.Register(ctx, scheme, f)
}

func NewIterator(ctx context.Context, uri string) (Iterator, error) {

	err := ensureIterators()

	if err != nil {
		return nil, err
	}

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := iterators.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	f := i.(IteratorInitializeFunc)
	return f(ctx, uri)
}
