package fs

import (
	"context"
	"io"
	"os"
)

const STDIN string = "STDIN"

func readerFromPath(ctx context.Context, abs_path string) (io.ReadCloser, error) {

	if abs_path == STDIN {
		return os.Stdin, nil
	}

	fh, err := os.Open(abs_path)

	if err != nil {
		return nil, err
	}

	return fh, nil
}
