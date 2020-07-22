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
$> ./bin/emit -h
Usage of ./bin/emit:
  -data-source string
    	A valid whosonfirst/go-whosonfirst-index URI. (default "directory://")
  -format-json
    	Format JSON output for each record.
  -json
    	Emit a JSON list.
  -null
    	Emit to /dev/null
  -oembed
    	Emit results as OEmbed records
  -oembed-author-name string
    	A default value for the OEmbed 'author_name' property. (default "SFO Museum")
  -oembed-author-uri-template string
    	A valid RFC 6570 URI template to use for the OEmbed 'author_url' property. (default "https://millsfield.sfomuseum.org/id/{wof_id}")
  -oembed-media-label string
    	A valid (WOF) media:properties.sizes property label to identify image data. (default "z")
  -oembed-media-uri-template fmt
    	A valid Go language fmt template for constucting a RFC 6570 URI template to use for the OEmbed 'url' property. (default "https://millsfield.sfomuseum.org/media/%s/%d_{secret}_{label}.{extension}")
  -oembed-provider-name string
    	A default value for the OEmbed 'provider_name' property. (default "SFO Museum")
  -oembed-provider-url string
    	A default value for the OEmbed 'provider_url' property. (default "https://millsfield.sfomuseum.org/")
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
$> bin/emit -data-source directory:// \
	/usr/local/data/sfomuseum-data-media/data/
```

#### [featurecollection://](https://github.com/whosonfirst/go-whosonfirst-index/fs)

Emit data from one or more files containing GeoJSON `FeatureCollection` (WOF) records.

```
$> bin/emit -data-source featurecollection:// \
	/path/to/featurecollection.geojson
```

Feature collection records may also be read from `STDIN`. For example:

```
$> cat /path/to/featurecollection.geojson \

   | bin/emit -data-source featurecollection:// STDIN
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

   | bin/emit -data-source filelist:// STDIN
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

   | bin/emit -data-source geojsonls:// STDIN
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

```
$> ./bin/emit /usr/local/data/sfomuseum-data-media/data/

{"bbox":[-122.387197,37.619087,-122.387197,37.619087],"geometry":{"coordinates":[-122.387197,37.619087],"type":"Point"},"id":1159341477,"properties":{"edtf:cessation":"uuuu","edtf:inception":"uuuu","geom:area":0,"geom:area_square_m":0,"geom:bbox":"-122.387197,37.619087,-122.387197,37.619087","geom:latitude":37.619087,"geom:longitude":-122.387197,"iso:country":"US","media:created":1443524024,"media:fingerprint":"857752f82858b46502479f803da8a52f1e168d5e","media:imagehash_avg":"a:f8ffbe8070f0c080","media:imagehash_diff":"d:40b0383aa3870f3d","media:medium":"image","media:mimetype":"image/jpeg","media:properties":{"colours":[{"closest":[{"hex":"#9c2542","name":"Big Dip O' Ruby","reference":"crayola"},{"hex":"#a52a2a","name":"brown","reference":"css4"}],"hex":"#8e362e","name":"#8e362e","reference":"vibrant"},{"closest":[{"hex":"#a5694f","name":"Sepia","reference":"crayola"},{"hex":"#a0522d","name":"sienna","reference":"css4"}],"hex":"#9c5b59","name":"#9c5b59","reference":"vibrant"},{"closest":[{"hex":"#cdc5c2","name":"Silver","reference":"crayola"},{"hex":"#c0c0c0","name":"silver","reference":"css4"}],"hex":"#c1c4c4","name":"#c1c4c4","reference":"vibrant"},{"closest":[{"hex":"#414a4c","name":"Outer Space","reference":"crayola"},{"hex":"#2f4f4f","name":"darkslategrey","reference":"css4"}],"hex":"#4b4a4a","name":"#4b4a4a","reference":"vibrant"}],"depicts":["1159160617"],"medium":"image","mimetype":"image/jpeg","sizes":{"b":{"extension":"jpg","height":682,"mimetype":"image/jpeg","secret":"PH2VfvpO7NRUAyns7QxdvOO2YfwYhyib2QK5FleNQ34LmzliMWYN","width":1024},"c":{"extension":"jpg","height":533,"mimetype":"image/jpeg","secret":"PH2VfvpO7NRUAyns7QxdvOO2YfwYhyib2QK5FleNQ34LmzliMWYN","width":800},"d":{"extension":"jpg","height":320,"mimetype":"image/jpeg","secret":"PH2VfvpO7NRUAyns7QxdvOO2YfwYhyib2QK5FleNQ34LmzliMWYN","width":320},"dd":{"extension":"jpg","height":533,"mimetype":"image/jpeg","secret":"PH2VfvpO7NRUAyns7QxdvOO2YfwYhyib2QK5FleNQ34LmzliMWYN","width":800},"n":{"extension":"jpg","height":213,"mimetype":"image/jpeg","secret":"PH2VfvpO7NRUAyns7QxdvOO2YfwYhyib2QK5FleNQ34LmzliMWYN","width":320},"o":{"extension":"jpg","height":2400,"mimetype":"image/jpeg","secret":"ACqgPxrMRzHdzGRJfkBqGVMtP2L9gTrn7mgfhMJesqhjXWJpmRK7","width":3600},"sq":{"extension":"jpg","height":320,"mimetype":"image/jpeg","secret":"PH2VfvpO7NRUAyns7QxdvOO2YfwYhyib2QK5FleNQ34LmzliMWYN","width":320},"z":{"extension":"jpg","height":426,"mimetype":"image/jpeg","secret":"PH2VfvpO7NRUAyns7QxdvOO2YfwYhyib2QK5FleNQ34LmzliMWYN","width":640}},"source":"user","status_id":0},"media:source":"sfomuseum","media:status_id":1,"mz:hierarchy_label":1,"mz:is_approximate":1,"mz:is_current":-1,"sfomuseum:placetype":"image","src:geom":"unknown","wof:belongsto":[102527513,102087579,1159341477,85688637,1159396315,1159396157,1159396321,102191575,85633793,85922583],"wof:breaches":[],"wof:country":"US","wof:created":1528920235,"wof:depicts":[1159396315,102527513,1360516119,1159396321,1159396157,1159160617],"wof:geomhash":"86e9b7d5fe1f6f1a6479fa62588a1dea","wof:hierarchy":[{"building_id":1159396321,"campus_id":102527513,"concourse_id":1159396315,"continent_id":102191575,"country_id":85633793,"county_id":102087579,"locality_id":85922583,"media_id":1159341477,"neighbourhood_id":-1,"region_id":85688637,"wing_id":1159396157}],"wof:id":1159341477,"wof:lastmodified":1577131152,"wof:name":"Installation view of \"The Nation’s Game: A History of the National Football League\"","wof:parent_id":1159396315,"wof:placetype":"media","wof:repo":"sfomuseum-data-media","wof:superseded_by":[],"wof:supersedes":[],"wof:tags":[]},"type":"Feature"}
...and so on
```

Or:

```
$> ./bin/emit --format-json /usr/local/data/sfomuseum-data-media/data/

{
  "id": 1377012109,
  "type": "Feature",
  "properties": {
    "edtf:cessation": "2017-01-30",
    "edtf:inception": "2016-05-26",
    "geom:area": 0,
    ...
  },
  "geometry": "..."
}
and so on...
```

If you want to emit records as a valid JSON list then enable the `-json` flag.

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

Image data for OEmbed records is expected to be found in one or more properties in the `properties.media:properties` dictionary of the WOF record. Specifically:

* `properties.media:properties.sizes.{STRING_LABEL}` - A dictionary containing dimensions and other details for constructing a URL for an image identified by a string label.
* `properties.media:properties.uri_template` - A valid RFC6570 URI template for constructing an image URL using the details derived from `properties.media:properties.sizes.{STRING_LABEL}`.

Here is [an abbreviated example](https://raw.githubusercontent.com/sfomuseum-data/sfomuseum-data-media/master/data/115/933/962/7/1159339627.geojson) of a WOF record with `media:` properties.

```
{
  "id": 1159339627,
  "type": "Feature",
  "properties": {
    ...  		
    "media:created": 1508957796,
    "media:fingerprint": "fd6e55e1ea940673e8dd7edfdacf0c2d546b8d6a",
    "media:imagehash_avg": "a:c2e7e7effbff0000",
    "media:imagehash_diff": "d:868e8e9e92c2eaba",
    "media:medium": "image",
    "media:mimetype": "image/jpeg",
    "media:properties": {
      "medium": "image",
      "mimetype": "image/jpeg",
      "sizes": {
        ...      	       
        "z": {
          "extension": "jpg",
          "height": 480,
          "mimetype": "image/jpeg",
          "secret": "UaqY5CItyrimU82DjYTYxy6XfRZXO0tD1YfBHWYhLnxGK1id8sdf",
          "width": 320
        },
	"uri_template": "https://millsfield.sfomuseum.org/media/115/933/962/7/1159339627_{secret}_{label}.{extension}"
      }    			
    }
}
```

If a WOF record does contain a `media:properties.uri_template` property then the value of the `-oembed-media-uri-template` flag will be used to construct a URI template. For example:

```
	wof_id := 1159339627
	wof_tree := "115/933/962/7"
	media_template_uri = fmt.Sprintf(opts.MediaURITemplate, wof_tree, wof_id)
```

_Important: Everything involving media in WOF documents remains a work in progress and is still subject to change._

If no relevant media information can be found in a WOF feature then the code will render the feature's geometry as an SVG image and assign a base-64 encoded data URL of the representation to the OEmbed record's `url` property. For example:

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

* https://github.com/whosonfirst/go-whosonfirst-index