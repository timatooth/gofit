# Gofit

Uses the Fitbit API to load data into a TSDB like Prometheus for displaying
data into grafana.

## Requirements
* Go 1.8
* Fitbit API App ID/Secret (You need to request your own personal App keys in Fitbit dashboard)
* Prometheus (soon)
* Grafana (soon)

## Running

    export FITBIT_CLIENT_ID=XXXXXX
    export FITBIT_CLIENT_SECRET=xxxxxxxxxxxxxxxxxxx
    go build
    ./gofit
    # You will be prompted to visit the Fitbit authorisation grant url.

### Screenshot
![Step Data](http://i.imgur.com/MdufcMC.png)