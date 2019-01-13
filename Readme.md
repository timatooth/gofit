# Gofit
> Uses the Fitbit API to load data into InfluxDB for displaying data into Grafana for FitOps engineers.

[![Build Status](https://travis-ci.org/timatooth/gofit.svg?branch=master)](https://travis-ci.org/timatooth/gofit)

### Requirements
* Go 1.8+
* Fitbit API App ID/Secret (You need to request your own personal App keys in the Fitbit dashboard)
* InfluxDB 1.2
* Grafana 4.2+

#### Register App with FitBit API

Go to [dev.fitbit.com](https://dev.fitbit.com) and select **Register an App** and enter `http://localhost:4000/auth` under **Callback URL**. For all other URLs you can enter any URL, like `https://example.com`. Saving the app will produce a client ID and secret that you will need to set below before running gofit.

##### Todo
* Use the Fitbit subscriptions API to send metrics in near-realtime when a fitness tracker syncs with Fitbit Servers.
* Prometheus Fitbit exporter.

## Installing
Assuming you have set your $GOPATH https://golang.org/doc/code.html#GOPATH
Optionally add `$PATH=$PATH:$GOPATH/bin` to make `gofit` command available.

    # Install gofit and dependencies
    go get -v github.com/timatooth/gofit

## Running

    # Start InfluxDB and Grafana
    docker-compose up -d

    # Create the fitbit database in InfluxDB
    curl "http://localhost:8086/query" --data-urlencode "q=CREATE DATABASE fitbit"

    # Load in data with Gofit
    export FITBIT_CLIENT_ID=XXXXXX
    export FITBIT_CLIENT_SECRET=xxxxxxxxxxxxxxxxxxx
    $GOPATH/bin/gofit
    # You will be prompted to visit the Fitbit authorisation grant url.

### Screenshot
![Step Data](http://i.imgur.com/MdufcMC.png)

### Grafana

Go to [localhost:3000](http://localhost:3000), log in as admin/admin, click "Create your first data source" and enter:
* Name: fitbit
* Type: InfluxDB
* Url: http://influxdb:8086
* Database: fitbit

Dashboards json exports included for importing into Grafana inside `grafana/`.
