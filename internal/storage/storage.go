package storage

import (
	"context"
	"fmt"
	"strconv"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	log "github.com/sirupsen/logrus"
	"github.com/tobiasbrodd/hivee-api/internal/coretypes"
)

type Storage struct {
	Influx       *influxdb2.Client
	Organization string
}

func (storage *Storage) ReadMeasureHistory(start string, stop string, measurement string, location string, every string) []coretypes.Measure {
	log.Infof("Reading measurement %s: %s\n", measurement, location)
	if len(start) == 0 {
		start = "-30d"
	}
	if len(stop) == 0 {
		stop = "now()"
	}
	if len(measurement) == 0 {
		measurement = "temperature"
	}
	if len(location) == 0 {
		location = "Indoor"
	}
	if len(every) == 0 {
		every = "1d"
	}

	bucket := "hivee"
	query := fmt.Sprintf(`from(bucket: "%s")
	|> range(start: "%s", stop: "%s")
	|> filter(fn: (r) => r["_measurement"] == "%s")
	|> filter(fn: (r) => r["_field"] == "value")
	|> filter(fn: (r) => r["location"] == "%s")
	|> aggregateWindow(every: "%s", fn: mean, createEmpty: false)
	|> yield(name: "mean")`, bucket, start, stop, measurement, location, every)

	reader := storage.getReader()
	result, err := (*reader).Query(context.Background(), query)

	var history []coretypes.Measure
	if err == nil {
		for result.Next() {
			value, _ := strconv.ParseFloat(fmt.Sprintf("%v", result.Record().Value()), 64)
			timestamp := result.Record().Time().Unix()
			history = append(history, coretypes.Measure{Value: value, Timestamp: timestamp})
		}
		if result.Err() != nil {
			log.Errorf("Query parsing error: %s\n", result.Err().Error())
		}
	} else {
		panic(err)
	}

	return history
}

func (storage Storage) getReader() *api.QueryAPI {
	reader := (*storage.Influx).QueryAPI(storage.Organization)

	return &reader
}

func New(authToken string, host string, port int, org string) *Storage {
	client := influxdb2.NewClient(fmt.Sprintf("http://%s:%d", host, port), authToken)
	storage := &Storage{Influx: &client, Organization: org}

	return storage
}

func (storage Storage) Close() {
	(*storage.Influx).Close()
}
