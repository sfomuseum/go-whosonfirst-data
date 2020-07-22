# go-whosonfirst-data

Tools for working with Who's On First (WOF) style documents.

## Important

Work in progress. This package might be renamed.

## Tools

To build binary versions of these tools run the `cli` Makefile target. For example:

```
> make cli
go build -mod vendor -o bin/emit cmd/emit/main.go
```

## emit

A command-line tool for parsing and emitting individual WOF records from a WOF data source.

```
> ./bin/emit -h
Usage of ./bin/emit:
  -data-source string
    	A valid whosonfirst/go-whosonfirst-index data source URI. (default "directory://")
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
```

### Data sources

The `emit` tool uses the [go-whosonfirst-index](https://github.com/whosonfirst/go-whosonfirst-index) package to read data from a variety of sources. The following data sources are supported by default:
 
#### [directory://](https://github.com/whosonfirst/go-whosonfirst-index/fs)

Emit data from one or more directories containing GeoJSON (WOF) records.

For example:

```
$> bin/emit -data-source directory:// /usr/local/data/sfomuseum-data-media/data/
```

#### [featurecollection://](https://github.com/whosonfirst/go-whosonfirst-index/fs)

Emit data from one or more files containing GeoJSON `FeatureCollection` (WOF) records.

```
$> bin/emit -data-source featurecollection:// /path/to/featurecollection.geojson
```

Feature collection records may also be read from `STDIN`. For example:

```
$> cat /path/to/featurecollection.geojson \

   bin/emit -data-source featurecollection:// STDIN
```   

#### [file://](https://github.com/whosonfirst/go-whosonfirst-index/fs)

Emit data from one or more files.

For example:

```
$> bin/emit -data-source file:// \

	/path/to/feature1.geojson \
	/path/to/feature2.geojson	
```

#### [filelist://](https://github.com/whosonfirst/go-whosonfirst-index/fs)

Emit data from a list of files (to read and emit).

For example:

```
$> bin/emit -data-source filelist:// /path/to/files.txt
```

File lists may also be read from `STDIN`. For example:

```
$> cat /path/to/files.txt \

   bin/emit -data-source filelist:// STDIN
```   

#### [geojsonls://](https://github.com/whosonfirst/go-whosonfirst-index/fs)

Emit data encoded in one or more line-delimited GeoJSON files.

For example:

```
$> bin/emit -data-source geojsonls:// /path/to/features.jsonl
```

Line-delimited GeoJSON records may also be read from `STDIN`. For example:

```
$> cat /path/to/features.jsonl \

   bin/emit -data-source geojsonls:// STDIN
```   

#### [git://](https://github.com/whosonfirst/go-whosonfirst-index-git)

Emit data from one or more Git repositories:

For example:

```
$> bin/emit -data-source git:// \
	https://github.com/sfomuseum-data/sfomuseum-data-flights-2020-07.git \
	https://github.com/sfomuseum-data/sfomuseum-data-flights-2020-06.git	
```

#### [repo://](https://github.com/whosonfirst/go-whosonfirst-index/fs)

Emit data from one or more directories where the relevent WOF data is expected to be found in a `data` subdirectory.

For example:

```
$> bin/emit -data-source repo:// \
	/usr/local/data/sfomuseum-data-media \
	/usr/local/data/sfomuseum-data-media-collection	
```

### JSON

By default all records are emitted as line-delimited JSON records. A side-effect of this is that the default WOF formatting is lost. In order to preserve the original formatting pass in the `-format-json` flag.

### Inline queries

You can also specify inline queries by passing a `-query` parameter which is a string in the format of:

```
{PATH}={REGULAR EXPRESSION}
```

Paths follow the dot notation syntax used by the [tidwall/gjson](https://github.com/tidwall/gjson) package and regular expressions are any valid [Go language regular expression](https://golang.org/pkg/regexp/). Successful path lookups will be treated as a list of candidates and each candidate's string value will be tested against the regular expression's [MatchString](https://golang.org/pkg/regexp/#Regexp.MatchString) method.

For example:

```
$> bin/emit /usr/local/data/sfomuseum-data-media/data/ \

   -query 'properties.wof:belongs_to=\b102087579\b' \

   | wc -l

1122
```

You can pass multiple `-query` parameters. The default query mode is to ensure that all queries match but you can also specify that only one or more queries need to match by passing the `-query-mode ANY` flag.

For example, this is how you would query the `sfomuseum-data-flights-2020-07` Git repository filtering for records involving either Boeing 737-8 or Airbus A321 aircraft. The results are emitted as a JSON list and piped to the `jq` tool which prints their `sfomuseum:flight_id` property:

```
> ./bin/emit \
	-json \
	-query 'properties.icao:aircraft=B738' \
	-query 'properties.icao:aircraft=A321' \
	-query-mode ANY \
	-data-source git:// \
	https://github.com/sfomuseum-data/sfomuseum-data-flights-2020-07.git \

   | jq '.[]["properties"]["sfomuseum:flight_id"]'
   
"20200701-A-DAL-696"
"20200701-A-DAL-807"
"20200701-A-DAL-958"
"20200701-A-JBU-115"
"20200701-A-JBU-1415"
"20200701-A-JBU-1833"
"20200701-A-JBU-415"
"20200701-A-JBU-577"
"20200701-A-JBU-915"
"20200701-A-SCX-395"
"20200701-A-SWA-1654"
"20200701-A-SWA-1817"
"20200701-A-SWA-2065"
"20200701-A-SWA-300"
"20200701-A-SWA-930"
...and so on
"20200721-D-SCX-396"
"20200721-D-SWA-1655"
"20200721-D-SWA-1693"
"20200721-D-SWA-3244"
"20200721-D-SWA-946"
"20200721-D-UAL-1273"
"20200721-D-UAL-1578"
"20200721-D-UAL-352"
"20200721-D-UAL-355"
"20200721-D-UAL-367"
"20200721-D-UAL-673"
```

### OEmbed

It is also possible to emit OpenAccess records as [OEmbed](https://oembed.com/) documents of type "photo".

For example:

```
$> ./bin/emit \
	-format-json \
	-oembed \
	/usr/local/data/sfomuseum-data-media/data/

{
  "version": "1.0",
  "type": "photo",
  "width": 640,
  "height": 453,
  "title": "Installation view of \"Before the 21st Century: An Ode to Boats, Cars, Motorcycles, Planes, and Trains\"",
  "url": "https://millsfield.sfomuseum.org/media/137/702/095/5/1377020955_GkHONnz4lqxYWQ9me6mBLNmZdthfrTKv_z.jpg",
  "author_name": "SFO Museum",
  "author_url": "https://millsfield.sfomuseum.org/id/1377020955",
  "provider_name": "SFO Museum",
  "provider_url": "https://millsfield.sfomuseum.org/",
  "object_uri": "wof://id/1377020955"
}
{
  "version": "1.0",
  "type": "photo",
  "width": 640,
  "height": 380,
  "title": "Installation view of \"Airline Identity: Marks, Brands and Logos\"",
  "url": "https://millsfield.sfomuseum.org/media/137/704/368/7/1377043687_agKvxo3EzdgRyNUahkBXadodNcV0Vvgx_z.jpg",
  "author_name": "SFO Museum",
  "author_url": "https://millsfield.sfomuseum.org/id/1377043687",
  "provider_name": "SFO Museum",
  "provider_url": "https://millsfield.sfomuseum.org/",
  "object_uri": "wof://id/1377043687"
}
... and so on
```

...if not media information can be found then the code will create a base-64 encoded data URL for the feature's geometry.

```
$> ./bin/emit \
	-format-json \
	-oembed \
	-query 'properties.wof:concordances.iata:code=SFO' \
	/usr/local/data/sfomuseum-data-whosonfirst/data/

{
  "version": "1.0",
  "type": "photo",
  "width": 800,
  "height": 640,
  "title": "San Francisco International Airport",
  "url": "data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iODAwLjAwMDAwMCIgaGVpZ2h0PSI2NDAuMDAwMDAwIiB2aWV3Qm94PSIwIDAgODAwIDY0MCIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj48cGF0aCBkPSJNNDguMDYzMTI4IDMwOC4yMDk3ODg ...(truncated for the sake of brevity)... IFoiLz48L3N2Zz4=",
  "author_name": "SFO Museum",
  "author_url": "https://millsfield.sfomuseum.org/id/102527513",
  "provider_name": "SFO Museum",
  "provider_url": "https://millsfield.sfomuseum.org/",
  "object_uri": "wof://id/102527513"
}
```

If you put the value of the `url` property in to an HTML `<img />` you'd see this:

![](docs/images/oembed-sfo-svg.png)

## See also

* https://github.com/whosonfirst/go-whosonfirst-iterator