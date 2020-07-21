package fs

import (
	"context"
	"errors"
	"github.com/whosonfirst/go-whosonfirst-crawl"
	"github.com/whosonfirst/go-whosonfirst-iterator"
	"net/url"
	"os"
	"path/filepath"
)

func init() {
	ctx := context.Background()
	err := iterator.RegisterIterator(ctx, "directory", NewDirectoryIterator)

	if err != nil {
		panic(err)
	}
}

type DirectoryIterator struct {
	iterator.Iterator
	root string
}

func NewDirectoryIterator(ctx context.Context, uri string) (iterator.Iterator, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	root := u.Path

	abs_root, err := filepath.Abs(root)

	if err != nil {
		return nil, err
	}

	info, err := os.Stat(abs_root)

	if err != nil {
		return nil, err
	}

	if !info.IsDir() {
		return nil, errors.New("Path is not a directory")
	}

	it := &DirectoryIterator{
		root: root,
	}

	return it, nil
}

func (it *DirectoryIterator) IterateWithChannels(ctx context.Context, uri string, result_ch chan iterator.Result, err_ch chan error, done_ch chan bool) {

	defer func() {
		done_ch <- true
	}()

	select {
	case <-ctx.Done():
		return
	default:
		// pass
	}

	// check for '*' and/or others?

	match, err := filepath.Match("..", uri)

	if err != nil {
		err_ch <- err
		return
	}

	if match {
		err_ch <- errors.New("Invalid URI")
		return
	}

	path := filepath.Join(it.root, uri)

	info, err := os.Stat(path)

	if err != nil {
		err_ch <- err
		return
	}

	if !info.IsDir() {
		err_ch <- errors.New("URI is not a directory")
		return
	}

	crawl_cb := func(path string, info os.FileInfo) error {

		select {
		case <-ctx.Done():
			return nil
		default:
			// pass
		}

		if info.IsDir() {
			return nil
		}

		fh, err := readerFromPath(ctx, path)

		if err != nil {
			return err
		}

		defer fh.Close()

		result, err := iterator.NewIteratorResultWithReader(ctx, path, fh)

		if err != nil {
			return err
		}

		result_ch <- result
		return nil
	}

	c := crawl.NewCrawler(path)
	err = c.Crawl(crawl_cb)

	if err != nil {
		err_ch <- err
		return
	}

	return
}
