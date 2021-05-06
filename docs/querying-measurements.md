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
{"timestamp":"2021-02-24T21:14:33.094407Z","name":"env.temperature.gen","value":1.8749457256338125,"meta":{"node":"0000000000000001","plugin":"metsense:1.0.2"}}
{"timestamp":"2021-02-24T21:14:34.097759Z","name":"env.temperature.gen","value":3.4616782879021497,"meta":{"node":"0000000000000001","plugin":"metsense:1.0.2"}}
{"timestamp":"2021-02-24T21:14:35.099309Z","name":"env.temperature.gen","value":3.935407701067743,"meta":{"node":"0000000000000001","plugin":"metsense:1.0.2"}}
{"timestamp":"2021-02-24T21:14:36.102012Z","name":"env.temperature.gen","value":0.660707909927028,"meta":{"node":"0000000000000001","plugin":"metsense:1.0.2"}}
{"timestamp":"2021-02-24T21:14:37.104884Z","name":"env.temperature.gen","value":0.5932408953781276,"meta":{"node":"0000000000000001","plugin":"metsense:1.0.2"}}
```

Each record contains the following fields

* `timestamp`: Timestamp of when measurement was taken.
* `name`: Name of measurement.
* `value`: Value of measurement.
* `meta`: Metadata fields about measurement.


## Example Queries

The following query will return all measurements with a name starting with `sys` in the last hour.

```sh
curl -H 'Content-Type: application/json' https://sdr.sagecontinuum.org/api/v1/query -d '
{
    "start": "-1h",
    "filter": {
        "name": "sys.*"
    }
}
'
```

The following query will return all environmental related measurements between 10:00 and 12:00 on 2021-01-01.

```sh
curl -H 'Content-Type: application/json' https://sdr.sagecontinuum.org/api/v1/query -d '
{
    "start": "2021-01-01T10:00:00Z",
    "end": "2021-01-01T12:00:00Z",
    "filter": {
        "name": "env.*"
    }
}
'
```

The following query will find all temperature related measurements from metsense v1.x plugins in the last 24 hours.

```sh
curl -H 'Content-Type: application/json' https://sdr.sagecontinuum.org/api/v1/query -d '
{
    "start": "-24h",
    "filter": {
        "plugin": "metsense:1.*",
        "name": "env.temperature.*"
    }
}
'
```

The following query will get the latest uptime measurements from all devices in the last 7 days.

```sh
curl -H 'Content-Type: application/json' https://sdr.sagecontinuum.org/api/v1/query -d '
{
    "start": "-7d",
    "tail": 1,
    "filter": {
        "name": "sys.uptime"
    }
}
'
```
