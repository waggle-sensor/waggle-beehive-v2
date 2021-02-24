# Beehive InfluxDB Data API

This service provides an implementation of the SDR API in front of InfluxDB.

## Dev Notes

We are currently evaluating InfluxDB as our backend alongside rolling our own query system on top of PostgesSQL or Cassandra as a fallback.

## Querying Data

The `/api/v1/query` provides a simple query system. Here's an example:

```txt
curl -X POST http://service-ip:10000/api/v1/query -d '
{
    "start": "2021-02-24T20:00:00Z",
    "end": "2021-02-24T21:00:00Z",
    "filter": {
        "name": "sys.*"
    }
}
'
```

This will return all measurements with a name starting with `sys` between 2021-02-24 20:00:00 and 2021-02-24 21:00:00. The filter fields can match any of the measurement metadata. This includes common things like `node` and `plugin`.

Here are some more examples:

```txt
curl -X POST http://service-ip:10000/api/v1/query -d '
{
    "start": "-1h",
    "filter": {
        "node": "0000000000000001"
    }
}
'
```

This queries all measurements from node 0000000000000001 in the last hour.

```txt
curl -X POST http://service-ip:10000/api/v1/query -d '
{
    "start": "-1h",
    "filter": {
        "node": "0000000000000001",
        "plugin": "metsense:.*"
    }
}
'
```

This queries all measurements from any metsense plugin version running on node 0000000000000001 in the last hour.

## Query Results

Queries are provided as line-delimited JSON. For example:

```txt
{"timestamp":"2021-02-24T21:14:33.094407742Z","name":"env.temperature.gen","value":1.8749457256338125,"meta":{"node":"0000000000000001","plugin":"metsense:1.0.2"}}
{"timestamp":"2021-02-24T21:14:34.097759221Z","name":"env.temperature.gen","value":3.4616782879021497,"meta":{"node":"0000000000000001","plugin":"metsense:1.0.2"}}
{"timestamp":"2021-02-24T21:14:35.099309444Z","name":"env.temperature.gen","value":3.935407701067743,"meta":{"node":"0000000000000001","plugin":"metsense:1.0.2"}}
{"timestamp":"2021-02-24T21:14:36.102012155Z","name":"env.temperature.gen","value":0.660707909927028,"meta":{"node":"0000000000000001","plugin":"metsense:1.0.2"}}
{"timestamp":"2021-02-24T21:14:37.104884504Z","name":"env.temperature.gen","value":0.5932408953781276,"meta":{"node":"0000000000000001","plugin":"metsense:1.0.2"}}
```

Each record contains:

* `timestamp`: Timestamp as nanoseconds since UNIX epoch time.
* `name`: Name of measurement.
* `value`: Value of measurement.
* `meta`: Metadata fields about measurement.
