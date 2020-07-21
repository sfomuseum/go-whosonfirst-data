package iterator

import (
	"context"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/feature"
)

func FeatureFromResult(ctx context.Context, result Result) (geojson.Feature, error) {

	b, err := result.Bytes()

	if err != nil {
		return nil, err
	}

	f, err := feature.LoadFeature(b)

	if err != nil {
		return nil, err
	}

	return f, nil
}
