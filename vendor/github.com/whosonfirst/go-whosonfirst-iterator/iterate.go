package iterator

import (
	"context"
)

type IterateCallbackFunc func(context.Context, Result, error) error

func IterateWithCallback(ctx context.Context, it Iterator, uri string, cb IterateCallbackFunc) error {

	iter_result_ch := make(chan Result)
	iter_err_ch := make(chan error)
	iter_done_ch := make(chan bool)

	err_ch := make(chan error)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go it.IterateWithChannels(ctx, uri, iter_result_ch, iter_err_ch, iter_done_ch)

	working := true

	for working {
		select {

		case <-ctx.Done():
			working = false
		case <-iter_done_ch:
			working = false
		case iter_err := <-iter_err_ch:

			go func() {
				// see what's happening here? we're letting the callback func
				// decide whether the error should stop iteration

				err := cb(ctx, nil, iter_err)

				if err != nil {
					err_ch <- err
				}
			}()

		case result := <-iter_result_ch:

			go func() {

				err := cb(ctx, result, nil)

				if err != nil {
					err_ch <- err
				}
			}()

		default:
			//
		}

		if !working {
			break
		}
	}

	return nil
}

func TestResultWithFilters(ctx context.Context, result Result, filters ...Filter) (bool, error) {

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	iter_err_ch := make(chan error)
	ok_ch := make(chan bool)
	iter_done_ch := make(chan bool)

	for _, f := range filters {

		go func(result Result, f Filter) {

			defer func() {
				iter_done_ch <- true
			}()

			ok, err := f.IncludeResult(ctx, result)

			if err != nil {
				iter_err_ch <- err
				return
			}

			ok_ch <- ok

		}(result, f)
	}

	remaining := len(filters)

	for remaining > 0 {
		select {
		case <-ctx.Done():
			return false, nil
		case <-iter_done_ch:
			remaining -= 1
		case err := <-iter_err_ch:
			return false, err
		case ok := <-ok_ch:
			if !ok {
				return false, nil
			}
		default:
			// pass
		}
	}

	return true, nil
}
