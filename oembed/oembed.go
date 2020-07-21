package oembed

import (
	"context"
	"errors"
	"fmt"
	"github.com/aaronland/go-wunderkammer/oembed"
	"github.com/jtacoma/uritemplates"
	"github.com/tidwall/gjson"
	"github.com/whosonfirst/go-whosonfirst-sources"
	"github.com/whosonfirst/go-whosonfirst-uri"
)

const WHOSONFIRST_URI_TEMPLATE string = "wof://id/{wofid}"

var whosonfirst_uri_template *uritemplates.UriTemplate

func init() {

	t, err := uritemplates.Parse(WHOSONFIRST_URI_TEMPLATE)

	if err != nil {
		panic(err)
	}

	whosonfirst_uri_template = t
}

func OEmbedRecordFromFeature(ctx context.Context, body []byte) (*oembed.Photo, error) {

	id_rsp := gjson.GetBytes(body, "properties.wof:id")

	if !id_rsp.Exists() {
		return nil, errors.New("Missing wof:id")
	}

	wof_id := id_rsp.Int()

	name_rsp := gjson.GetBytes(body, "properties.wof:name")

	if !name_rsp.Exists() {
		return nil, errors.New("Missing wof:name")
	}

	wof_name := name_rsp.String()

	author_name := "author..."

	author_url, err := uri.Id2RelPath(wof_id)

	if err != nil {
		return nil, err
	}

	values := make(map[string]interface{})
	values["wofid"] = wof_id

	object_uri, err := whosonfirst_uri_template.Expand(values)

	provider_name := "..."
	provider_url := "..."

	src_rsp := gjson.GetBytes(body, "properties.src:geom")

	if src_rsp.Exists() {

		src_name := src_rsp.String()

		provider_name = src_name
		provider_url = fmt.Sprintf("x-urn:src:geom#%s", src_name)

		src, err := sources.GetSourceByName(src_name)

		if err != nil {
			provider_name = src.Fullname
			provider_url = src.URL
		}
	}

	url := "fix me"

	o := &oembed.Photo{
		Version:      "1.0",
		Type:         "photo",
		Height:       -1, // https://github.com/Smithsonian/OpenAccess/issues/2
		Width:        -1, // see above
		URL:          url,
		Title:        wof_name,
		AuthorName:   author_name,
		AuthorURL:    author_url,
		ProviderName: provider_name,
		ProviderURL:  provider_url,
		ObjectURI:    object_uri,
	}

	return o, nil
}
