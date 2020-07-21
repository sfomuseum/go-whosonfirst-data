package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/aaronland/go-json-query"
	"github.com/sfomuseum/go-whosonfirst-data/oembed"
	"github.com/tidwall/pretty"
	"github.com/whosonfirst/go-whosonfirst-index"
	_ "github.com/whosonfirst/go-whosonfirst-index/fs"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
	"sync/atomic"
)

func main() {

	iter_uri := flag.String("uri", "directory:///", "A valid whosonfirst/go-whosonfirst-index URI.")

	to_stdout := flag.Bool("stdout", true, "Emit to STDOUT")
	to_devnull := flag.Bool("null", false, "Emit to /dev/null")

	as_json := flag.Bool("json", false, "Emit a JSON list.")
	format_json := flag.Bool("format-json", false, "Format JSON output for each record.")

	as_oembed := flag.Bool("oembed", false, "Emit results as OEmbed records")

	author_name := flag.String("oembed-author-name", "SFO Museum", "...")
	author_uri_template := flag.String("oembed-author-uri-template", "https://millsfield.sfomuseum.org/id/{wof_id}", "...")

	provider_name := flag.String("oembed-provider-name", "SFO Museum", "...")
	provider_url := flag.String("oembed-provider-url", "https://millsfield.sfomuseum.org/", "...")

	media_uri_template := flag.String("oembed-media-uri-template", "https://millsfield.sfomuseum.org/media/%s/%d_{secret}_{label}.{extension}", "...")
	media_label := flag.String("oembed-media-label", "z", "...")

	var queries query.QueryFlags
	flag.Var(&queries, "query", "One or more {PATH}={REGEXP} parameters for filtering records.")

	valid_modes := strings.Join([]string{query.QUERYSET_MODE_ALL, query.QUERYSET_MODE_ANY}, ", ")
	desc_modes := fmt.Sprintf("Specify how query filtering should be evaluated. Valid modes are: %s", valid_modes)

	query_mode := flag.String("query-mode", query.QUERYSET_MODE_ALL, desc_modes)

	flag.Parse()

	ctx := context.Background()

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

	cb := func(ctx context.Context, fh io.Reader, args ...interface{}) error {

		body, err := ioutil.ReadAll(fh)

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

		if *as_oembed {

			opts := &oembed.OEmbedOptions{
				AuthorName:        *author_name,
				AuthorURITemplate: *author_uri_template,
				ProviderName:      *provider_name,
				ProviderURL:       *provider_url,
				MediaURITemplate:  *media_uri_template,
				MediaLabel:        *media_label,
			}

			oembed_record, err := oembed.OEmbedRecordFromFeature(ctx, body, opts)

			if err != nil {
				return err
			}

			body, err = json.Marshal(oembed_record)

			if err != nil {
				return err
			}

		} else {

			var stub interface{}

			err = json.Unmarshal(body, &stub)

			if err != nil {
				return err
			}

			body, err = json.Marshal(stub)

			if err != nil {
				return err
			}

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

	idx, err := index.NewIndexer(*iter_uri, cb)

	if err != nil {
		log.Fatal(err)
	}

	uris := flag.Args()

	if *as_json {
		wr.Write([]byte("["))
	}

	err = idx.Index(ctx, uris...)

	if err != nil {
		log.Fatal(err)
	}

	if *as_json {
		wr.Write([]byte("]"))
	}

}
