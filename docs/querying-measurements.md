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
{"timestamp":"2021-09-14T19:09:38.329031006Z","name":"env.temperature","value":46.91,"meta":{"host":"000048b02d15bc87.ws-nxcore","node":"000048b02d15bc87","plugin":"plugin-iio:0.3.0","sensor":"bme280","vsn":"W019"}}
{"timestamp":"2021-09-14T19:10:08.374442041Z","name":"env.temperature","value":46.86,"meta":{"host":"000048b02d15bc87.ws-nxcore","node":"000048b02d15bc87","plugin":"plugin-iio:0.3.0","sensor":"bme280","vsn":"W019"}}
{"timestamp":"2021-09-14T19:10:38.431108286Z","name":"env.temperature","value":46.83,"meta":{"host":"000048b02d15bc87.ws-nxcore","node":"000048b02d15bc87","plugin":"plugin-iio:0.3.0","sensor":"bme280","vsn":"W019"}}
{"timestamp":"2021-09-14T19:11:08.49221237Z","name":"env.temperature","value":46.81,"meta":{"host":"000048b02d15bc87.ws-nxcore","node":"000048b02d15bc87","plugin":"plugin-iio:0.3.0","sensor":"bme280","vsn":"W019"}}
{"timestamp":"2021-09-14T19:11:38.525777857Z","name":"env.temperature","value":46.79,"meta":{"host":"000048b02d15bc87.ws-nxcore","node":"000048b02d15bc87","plugin":"plugin-iio:0.3.0","sensor":"bme280","vsn":"W019"}}
{"timestamp":"2021-09-14T19:12:08.575869104Z","name":"env.temperature","value":46.78,"meta":{"host":"000048b02d15bc87.ws-nxcore","node":"000048b02d15bc87","plugin":"plugin-iio:0.3.0","sensor":"bme280","vsn":"W019"}}
...
```

Each record contains the following fields

* `timestamp`: Timestamp of when measurement was taken.
* `name`: Name of measurement.
* `value`: Value of measurement.
* `meta`: Metadata fields about measurement.

_Note: Records are only ordered by timestamp **within each group of unique name and metadata**. This should generally not be an issue as many aggregations and visualizations process each group independently._

## Example Query

We'll start by doing a query to get all environment data in the last hour.

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

If we wanted look specifically at temperature values, we can update our `name` filter to be `env.temperature`:

```sh
curl -H 'Content-Type: application/json' https://data.sagecontinuum.org/api/v1/query -d '
{
    "start": "-1h",
    "filter": {
        "name": "env.temperature"
    }
}
'
```

Next, if we want to drill down to the specific node `W019`, we can add filter for the `vsn` field:

```sh
curl -H 'Content-Type: application/json' https://data.sagecontinuum.org/api/v1/query -d '
{
    "start": "-1h",
    "filter": {
        "name": "env.temperature",
        "vsn": "W019"
    }
}
'
```

Finally, if we want to get this data for a specific time range, we can provide absolute `start` and `end` timestamps:

```sh
curl -H 'Content-Type: application/json' https://data.sagecontinuum.org/api/v1/query -d '
{
    "start": "2021-09-14T00:00:00Z",
    "end": "2021-09-15T00:00:00Z",
    "filter": {
        "name": "env.temperature",
        "vsn": "W019"
    }
}
'
```

## More Examples

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
