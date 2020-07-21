# go-whosonfirst-data

## Important

Work in progress. This package might be renamed.

## Tools

To build binary versions of these tools run the `cli` Makefile target. For example:

```
> make cli
go build -mod vendor -o bin/emit cmd/emit/main.go
```

## emit

```
> ./bin/emit -h
Usage of ./bin/emit:
  -format-json
    	Format JSON output for each record.
  -json
    	Emit a JSON list.
  -null
    	Emit to /dev/null
  -query value
    	One or more {PATH}={REGEXP} parameters for filtering records.
  -query-mode string
    	Specify how query filtering should be evaluated. Valid modes are: ALL, ANY (default "ALL")
  -stdout
    	Emit to STDOUT (default true)
  -uri string
    	A valid whosonfirst/go-whosonfirst-iterator URI. (default "directory:///")
```

For example:

```
$> bin/emit /usr/local/data/sfomuseum-data-media/data/ \

   -query 'properties.wof:belongs_to=\b102087579\b' \

   | wc -l

1122
```

## See also

* https://github.com/whosonfirst/go-whosonfirst-iterator