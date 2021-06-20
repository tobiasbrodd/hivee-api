# hivee-api

API for Hivee

## Guide

To start **hivee-api**, you can either run it directly using **Go**, run it in **Docker** or use **Docker Compose**.

### Go

```
go build
go run main.go
```

### Docker

```
docker build -t hivee-api .
docker run -it --rm -p <8000>:8000 hivee-api
```

### Docker Compose

```
docker-compose up -d --build
```

### InfluxDB

**InfluxDB** needs to be running before starting **hivee-core**. Configuration for it should be in `config.yml`:
```
influx:
  token: <generated by InfluxDB>
  host: <host.docker.internal>
  port: <8086>
```

To start **InfluxDB**, run:

```
influxd
```

**InfluxDB** can also be started as a service on Linux:

```
sudo systemctl start influxdb
```
