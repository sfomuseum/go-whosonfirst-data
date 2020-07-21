# go-whosonfirst-iterator-fs

## Important

This is work in progress. It is meant to be a replacement for the `go-whosonfirst-index` packages. It will probably change.

## Background

* https://github.com/whosonfirst-data/whosonfirst-data/issues/1820#issuecomment-614176702

## Example

```
> go run -mod vendor cmd/count/main.go \
	-filter 'placetype:///?placetype=campus' \
	/usr/local/data/sfomuseum-data-whosonfirst/data/
	
542 records
```

## See also

* https://github.com/whosonfirst/go-whosonfirst-iterator
* https://github.com/whosonfirst/go-whosonfirst-index