package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func main() {
	addr := flag.String("addr", ":10000", "service addr")
	influxdbURL := flag.String("influxdb.url", getenv("INFLUXDB_URL", "http://localhost:8086"), "influxdb url")
	influxdbToken := flag.String("influxdb.token", getenv("INFLUXDB_TOKEN", ""), "influxdb token")
	influxdbBucket := flag.String("influxdb.bucket", getenv("INFLUXDB_BUCKET", ""), "influxdb bucket")
	influxdbTimeout := flag.Duration("influxdb.timeout", mustParseDuration(getenv("INFLUXDB_TIMEOUT", "15m")), "influxdb client timeout")
	flag.Parse()

	log.Printf("connecting to influxdb at %s", *influxdbURL)
	client := influxdb2.NewClient(*influxdbURL, *influxdbToken)
	defer client.Close()

	// TODO figure out reasonable timeout on potentially large result sets
	client.Options().HTTPClient().Timeout = *influxdbTimeout

	// NOTE temporarily redirecting to sage docs. can change to something better later.
	http.Handle("/", http.RedirectHandler("https://docs.sagecontinuum.org/docs/tutorials/accessing-data", http.StatusTemporaryRedirect))

	http.Handle("/api/v1/query", &Service{
		Backend: &InfluxBackend{
			Client: client,
			Org:    "waggle",
			Bucket: *influxdbBucket,
		},
	})

	// NOTE optional endpoint to expose testing bucket
	// http.Handle("/api/testing/query", &Service{
	// 	Backend: &InfluxBackend{
	// 		Client: client,
	// 		Org:    "waggle",
	// 		Bucket: "testing",
	// 	},
	// })

	log.Printf("service listening on %s", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal(err)
	}
}

func getenv(key string, fallback string) string {
	if s, ok := os.LookupEnv(key); ok {
		return s
	}
	return fallback
}

func mustParseDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		panic(err)
	}
	return d
}
