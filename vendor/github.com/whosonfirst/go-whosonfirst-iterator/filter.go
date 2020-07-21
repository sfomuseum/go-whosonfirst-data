package iterator

import (
	"context"
	"github.com/aaronland/go-roster"
	"net/url"
)

type Filter interface {
	IncludeResult(context.Context, Result) (bool, error)
}

type FilterInitializeFunc func(context.Context, string) (Filter, error)

var filters roster.Roster

func ensureFilters() error {

	if filters == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		filters = r
	}

	return nil
}

func RegisterFilter(ctx context.Context, scheme string, f FilterInitializeFunc) error {

	err := ensureFilters()

	if err != nil {
		return err
	}

	return filters.Register(ctx, scheme, f)
}

func NewFilter(ctx context.Context, uri string) (Filter, error) {

	err := ensureFilters()

	if err != nil {
		return nil, err
	}

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := filters.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	f := i.(FilterInitializeFunc)
	return f(ctx, uri)
}
