# Gofit
> Uses the Fitbit API to load data into InfluxDB for displaying data into Grafana for FitOps engineers.

[![Build Status](https://travis-ci.org/timatooth/gofit.svg?branch=master)](https://travis-ci.org/timatooth/gofit)

### Requirements
* Go 1.8+
* Fitbit API App ID/Secret (You need to request your own personal App keys in the Fitbit dashboard)
* InfluxDB 1.2
* Grafana 4.2+

##### Todo
* Use the Fitbit subscriptions API to send metrics in near-realtime when a fitness tracker syncs with Fitbit Servers.
* Prometheus Fitbit exporter.

## Installing
Assuming you have set your $GOPATH https://golang.org/doc/code.html#GOPATH
Optionally add `$PATH=$PATH:$GOPATH/bin` to make `gofit` command available.

    # Install gofit and dependencies
    go get -v github.com/timatooth/gofit

## Running

    export FITBIT_CLIENT_ID=XXXXXX
    export FITBIT_CLIENT_SECRET=xxxxxxxxxxxxxxxxxxx
    $GOPATH/bin/gofit
    # You will be prompted to visit the Fitbit authorisation grant url.
    # Assumes you have InfluxDB already running

### Screenshot
![Step Data](http://i.imgur.com/MdufcMC.png)

### Grafana
Dashboards json exports included for importing into Grafana inside `grafana/`.
