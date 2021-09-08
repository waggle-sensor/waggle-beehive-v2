# Querying Measurements

The following instructions will assume Beehive is deployed to the domain `sagecontinuum.org`, but all the steps should apply any other domain.

## Data Model

Beehive adopts a data model of "timeseries tagged with simple key-value metadata" similar to [InfluxDB](https://www.influxdata.com/products/influxdb/) and [Prometheus](https://prometheus.io). Abstractly, this means every measurement comes with:

* A timestamp of when measurement was taken.
* A measurement name. (ex. `env.temperature`)
* A measurement value. (ex. `23.1`)
* Simple key-value metadata pairs. (ex. `node=1234, host=rpi, plugin=metsense:1.0.3, camera=bottom`)

## Query API

Beehive's query API provides the following features:

* Time range selection
* Filtering by metadata
* Limiting number of results

The query request used by the API is a JSON body with the following structure

```json
{
    "start": "absolute or relative timestamp",
    "end": "absolute or relative timestamp",
    "tail": 100,
    "filter": {
        "tag1": "match pattern 1",
        "tag2": "match pattern 2",
        "...": "...",
    }
}
```

The `start`, `end`, `filter` and `tail` fields are all optional and can be included as needed.

The `start` and `end` fields specify an absolute or relative time range to query.

Absolute timestamps must be in a `YYYY-MM-DDTHH:MM:SSZ` format.

Relative timestamps must be in a `±ns`, `±nm` or `±nh` format where `s`, `m` and `h` indicate units of seconds, minutes and hours and `n` is the number. For example, `-4h` indicates 4 hours in the past.

The `tail` field limits results to the _most recent_ `n` records _for each_ unique combination of measurement name and meta fields.

## Query Response Format

Query responses are provided as newline separated JSON records. For example:

```json
{"timestamp":"2021-09-01T18:31:51.944475139Z","name":"env.temperature","value":51.18,"meta":{"node":"000048b02d05a0a4","plugin":"plugin-iio:0.2.0","sensor":"bme280"}}
{"timestamp":"2021-09-01T18:32:21.994670762Z","name":"env.temperature","value":51.11,"meta":{"node":"000048b02d05a0a4","plugin":"plugin-iio:0.2.0","sensor":"bme280"}}
{"timestamp":"2021-09-01T18:32:52.040360967Z","name":"env.temperature","value":51.03,"meta":{"node":"000048b02d05a0a4","plugin":"plugin-iio:0.2.0","sensor":"bme280"}}
{"timestamp":"2021-09-01T18:33:22.084612518Z","name":"env.temperature","value":50.93,"meta":{"node":"000048b02d05a0a4","plugin":"plugin-iio:0.2.0","sensor":"bme280"}}
{"timestamp":"2021-09-01T18:33:52.128641281Z","name":"env.temperature","value":50.87,"meta":{"node":"000048b02d05a0a4","plugin":"plugin-iio:0.2.0","sensor":"bme280"}}
{"timestamp":"2021-09-01T18:34:22.172853123Z","name":"env.temperature","value":50.79,"meta":{"node":"000048b02d05a0a4","plugin":"plugin-iio:0.2.0","sensor":"bme280"}}
{"timestamp":"2021-09-01T18:34:52.224589881Z","name":"env.temperature","value":50.69,"meta":{"node":"000048b02d05a0a4","plugin":"plugin-iio:0.2.0","sensor":"bme280"}}
{"timestamp":"2021-09-01T18:35:22.272575797Z","name":"env.temperature","value":50.62,"meta":{"node":"000048b02d05a0a4","plugin":"plugin-iio:0.2.0","sensor":"bme280"}}
{"timestamp":"2021-09-01T18:35:52.321517702Z","name":"env.temperature","value":50.56,"meta":{"node":"000048b02d05a0a4","plugin":"plugin-iio:0.2.0","sensor":"bme280"}}
{"timestamp":"2021-09-01T18:36:22.340393743Z","name":"env.temperature","value":50.51,"meta":{"node":"000048b02d05a0a4","plugin":"plugin-iio:0.2.0","sensor":"bme280"}}
```

Each record contains the following fields

* `timestamp`: Timestamp of when measurement was taken.
* `name`: Name of measurement.
* `value`: Value of measurement.
* `meta`: Metadata fields about measurement.


## Example Queries

The following query will get all environment data in the last hour.

```sh
curl -H 'Content-Type: application/json' https://data.sagecontinuum.org/api/v1/query -d '
{
    "start": "-1h",
    "filter": {
        "name": "env.*"
    }
}
'
```

The following query will return all measurements with a name starting with `sys` in the five minutes.

```sh
curl -H 'Content-Type: application/json' https://data.sagecontinuum.org/api/v1/query -d '
{
    "start": "-5m",
    "filter": {
        "name": "sys.*"
    }
}
'
```

The following query will return all environmental related measurements between 10:00 and 12:00 on 2021-01-01.

```sh
curl -H 'Content-Type: application/json' https://data.sagecontinuum.org/api/v1/query -d '
{
    "start": "2021-01-01T10:00:00Z",
    "end": "2021-01-01T12:00:00Z",
    "filter": {
        "name": "env.*"
    }
}
'
```

The following query will find all measurements from all IIO plugins v0.2.x in the last 24 hours.

```sh
curl -H 'Content-Type: application/json' https://data.sagecontinuum.org/api/v1/query -d '
{
    "start": "-24h",
    "filter": {
        "plugin": "plugin-iio:0.2.*"
    }
}
'
```

The following query will get the latest uptime measurements from all devices in the last 3 days.

```sh
curl -H 'Content-Type: application/json' https://data.sagecontinuum.org/api/v1/query -d '
{
    "start": "-3d",
    "tail": 1,
    "filter": {
        "name": "sys.uptime"
    }
}
'
```
