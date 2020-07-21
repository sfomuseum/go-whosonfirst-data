package iterator

import (
	"context"
	"github.com/rs/xid"
	"io"
	"io/ioutil"
)

type Result interface {
	URI() string
	Bytes() ([]byte, error)
	ID() string
}

type IteratorResult struct {
	Result
	uri    string
	reader io.Reader
	bytes  []byte
	id     *xid.ID
}

func NewIteratorResultWithReader(ctx context.Context, uri string, reader io.Reader) (Result, error) {

	r := &IteratorResult{
		uri:    uri,
		reader: reader,
	}

	return r, nil
}

func (r *IteratorResult) ID() string {

	if r.id == nil {
		id := xid.New()
		r.id = &id
	}

	return r.id.String()
}

func (r *IteratorResult) URI() string {
	return r.uri
}

func (r *IteratorResult) Bytes() ([]byte, error) {

	if r.bytes == nil {
		b, err := ioutil.ReadAll(r.reader)

		if err != nil {
			return nil, err
		}

		r.bytes = b
	}

	return r.bytes, nil
}
