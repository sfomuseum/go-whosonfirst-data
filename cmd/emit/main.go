package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/aaronland/go-json-query"
	"github.com/tidwall/pretty"
	"github.com/whosonfirst/go-whosonfirst-iterator"
	_ "github.com/whosonfirst/go-whosonfirst-iterator-fs"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
	"sync/atomic"
)

func main() {

	iter_uri := flag.String("uri", "directory:///", "A valid whosonfirst/go-whosonfirst-iterator URI.")

	to_stdout := flag.Bool("stdout", true, "Emit to STDOUT")
	to_devnull := flag.Bool("null", false, "Emit to /dev/null")

	as_json := flag.Bool("json", false, "Emit a JSON list.")
	format_json := flag.Bool("format-json", false, "Format JSON output for each record.")

	// as_oembed := flag.Bool("oembed", false, "Emit results as OEmbed records")

	var queries query.QueryFlags
	flag.Var(&queries, "query", "One or more {PATH}={REGEXP} parameters for filtering records.")

	valid_modes := strings.Join([]string{query.QUERYSET_MODE_ALL, query.QUERYSET_MODE_ANY}, ", ")
	desc_modes := fmt.Sprintf("Specify how query filtering should be evaluated. Valid modes are: %s", valid_modes)

	query_mode := flag.String("query-mode", query.QUERYSET_MODE_ALL, desc_modes)

	flag.Parse()

	ctx := context.Background()

	iter, err := iterator.NewIterator(ctx, *iter_uri)

	if err != nil {
		log.Fatalf("Failed to create new iterator, %v", err)
	}

	writers := make([]io.Writer, 0)

	if *to_stdout {
		writers = append(writers, os.Stdout)
	}

	if *to_devnull {
		writers = append(writers, ioutil.Discard)
	}

	if len(writers) == 0 {
		log.Fatal("Nothing to write to.")
	}

	wr := io.MultiWriter(writers...)

	var qs *query.QuerySet

	if len(queries) > 0 {

		qs = &query.QuerySet{
			Queries: queries,
			Mode:    *query_mode,
		}
	}

	mu := new(sync.RWMutex)

	counter := int32(0)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	cb := func(ctx context.Context, result iterator.Result, iter_err error) error {

		if iter_err != nil {
			return iter_err
		}

		body, err := result.Bytes()

		if err != nil {
			return err
		}

		// maybe just pass bytes.NewReader(body) to https://github.com/aaronland/go-jsonl/blob/master/walk/reader.go ?
		// (20200721/thisisaaronland)

		if qs != nil {

			matches, err := query.Matches(ctx, qs, body)

			if err != nil {
				return err
			}

			if !matches {
				return nil
			}
		}

		var stub interface{}

		err = json.Unmarshal(body, &stub)

		if err != nil {
			return err
		}

		body, err = json.Marshal(stub)

		if err != nil {
			return err
		}

		if *format_json {
			body = pretty.Pretty(body)
		}

		body = bytes.TrimSpace(body)

		mu.Lock()
		defer mu.Unlock()

		new_count := atomic.AddInt32(&counter, 1)

		if *as_json && new_count > 1 {
			wr.Write([]byte(","))
		}

		wr.Write(body)
		wr.Write([]byte("\n"))
		return nil
	}

	uris := flag.Args()

	if *as_json {
		wr.Write([]byte("["))
	}

	for _, uri := range uris {

		err := iterator.IterateWithCallback(ctx, iter, uri, cb)

		if err != nil {
			log.Fatalf("Failed to iterate URI '%s', %v", uri, err)
		}
	}

	if *as_json {
		wr.Write([]byte("]"))
	}

}
