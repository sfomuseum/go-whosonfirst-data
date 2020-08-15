package oembed

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aaronland/go-wunderkammer/oembed"
	"github.com/jtacoma/uritemplates"
	"github.com/tidwall/gjson"
	"github.com/whosonfirst/go-geojson-svg"
	_ "github.com/whosonfirst/go-whosonfirst-sources"
	"github.com/whosonfirst/go-whosonfirst-uri"
	_ "log"
)

type Media struct {
	Extension string `json:"extension"`
	Height    int    `json:"height"`
	Mimetype  string `json:"mimetype"`
	Secret    string `json:"secret"`
	Width     int    `json:"width"`
}

type OEmbedOptions struct {
	AuthorName          string
	AuthorURITemplate   string
	ProviderName        string
	ProviderURL         string
	MediaURITemplate    string
	MediaLabel          string
	ThumbnailMediaLabel string
}

const WHOSONFIRST_URI_TEMPLATE string = "wof://id/{wofid}"

var whosonfirst_uri_template *uritemplates.UriTemplate

func init() {

	t, err := uritemplates.Parse(WHOSONFIRST_URI_TEMPLATE)

	if err != nil {
		panic(err)
	}

	whosonfirst_uri_template = t
}

func OEmbedRecordFromFeature(ctx context.Context, body []byte, opts *OEmbedOptions) (*oembed.Photo, error) {

	id_rsp := gjson.GetBytes(body, "properties.wof:id")

	if !id_rsp.Exists() {
		return nil, errors.New("Missing wof:id")
	}

	wof_id := id_rsp.Int()

	wof_tree, err := uri.Id2Path(wof_id)

	if err != nil {
		return nil, err
	}

	name_rsp := gjson.GetBytes(body, "properties.wof:name")

	if !name_rsp.Exists() {
		return nil, errors.New("Missing wof:name")
	}

	wof_name := name_rsp.String()

	height := -1
	width := -1

	var media_template_uri string

	template_rsp := gjson.GetBytes(body, "properties.media:properties.media:uri_template")

	if !template_rsp.Exists() {
		media_template_uri = fmt.Sprintf(opts.MediaURITemplate, wof_tree, wof_id)

	} else {
		media_template_uri = template_rsp.String()
	}

	media_template, err := uritemplates.Parse(media_template_uri)

	if err != nil {
		return nil, err
	}

	var oembed_url string

	media_label := opts.MediaLabel
	media_path := fmt.Sprintf("properties.media:properties.sizes.%s", media_label)
	media_rsp := gjson.GetBytes(body, media_path)

	if media_rsp.Exists() {

		media_body := media_rsp.String()

		var m *Media

		err = json.Unmarshal([]byte(media_body), &m)

		if err != nil {
			return nil, err
		}

		media_values := make(map[string]interface{})
		media_values["secret"] = m.Secret
		media_values["extension"] = m.Extension
		media_values["label"] = media_label

		url, err := media_template.Expand(media_values)

		if err != nil {
			return nil, err
		}

		oembed_url = url

		height = m.Height
		width = m.Width

	} else {

		s := svg.New()
		// s.Mercator = *mercator

		err := s.AddFeature(string(body))

		if err != nil {
			return nil, err
		}

		width = 800
		height = 640

		svg_body := s.Draw(float64(width), float64(height),
			svg.WithAttribute("xmlns", "http://www.w3.org/2000/svg"),
			svg.WithAttribute("viewBox", fmt.Sprintf("0 0 %d %d", width, height)),
		)

		content_type := "image/svg+xml"

		b64_data := base64.StdEncoding.EncodeToString([]byte(svg_body))
		oembed_url = fmt.Sprintf("data:%s;base64,%s", content_type, b64_data)
	}

	author_name := opts.AuthorName

	author_template, err := uritemplates.Parse(opts.AuthorURITemplate)

	if err != nil {
		return nil, err
	}

	author_values := make(map[string]interface{})
	author_values["wof_tree"] = wof_tree
	author_values["wof_id"] = wof_id

	author_url, err := author_template.Expand(author_values)

	if err != nil {
		return nil, err
	}

	object_values := make(map[string]interface{})
	object_values["wofid"] = wof_id

	object_uri, err := whosonfirst_uri_template.Expand(object_values)

	provider_name := opts.ProviderName
	provider_url := opts.ProviderURL

	src_rsp := gjson.GetBytes(body, "properties.src:geom")

	if src_rsp.Exists() {

		/*
			src_name := src_rsp.String()
			src, err := sources.GetSourceByName(src_name)

			if err != nil {
				provider_name = src.Fullname
				provider_url = src.URL
			}
		*/
	}

	o := &oembed.Photo{
		Version:      "1.0",
		Type:         "photo",
		Height:       height,
		Width:        width,
		URL:          oembed_url,
		Title:        wof_name,
		AuthorName:   author_name,
		AuthorURL:    author_url,
		ProviderName: provider_name,
		ProviderURL:  provider_url,
		ObjectURI:    object_uri,
	}

	if opts.ThumbnailMediaLabel != "" {

		media_path := fmt.Sprintf("properties.media:properties.sizes.%s", opts.ThumbnailMediaLabel)
		media_rsp := gjson.GetBytes(body, media_path)

		if media_rsp.Exists() {

			media_body := media_rsp.String()

			var m *Media

			err = json.Unmarshal([]byte(media_body), &m)

			if err != nil {
				return nil, err
			}

			media_values := make(map[string]interface{})
			media_values["secret"] = m.Secret
			media_values["extension"] = m.Extension
			media_values["label"] = media_label

			url, err := media_template.Expand(media_values)

			if err != nil {
				return nil, err
			}

			o.ThumbnailURL = url

			o.ThumbnailHeight = m.Height
			o.ThumbnailWidth = m.Width
		} else {
			// do something here...
		}

	}

	return o, nil
}
