package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	client "github.com/influxdata/influxdb1-client/v2"
	"github.com/timatooth/gofit/fitbitapi"
)

func loadInfluxData(api *fitbitapi.Api) {
	var (
		InfluxHostname     = os.Getenv("INFLUX_HOSTNAME")
		InfluxDatabaseName = os.Getenv("INFLUX_DB")
		InfluxUsername     = os.Getenv("INFLUX_USERNAME")
		InfluxPassword     = os.Getenv("INFLUX_PASSWORD")
	)
	if InfluxHostname == "" {
		InfluxHostname = "http://localhost:8086"
	}
	if InfluxDatabaseName == "" {
		InfluxDatabaseName = "fitbit"
	}

	fmt.Println("Loading step data into influxdb...")
	activitySteps := api.GetActivitySteps()

	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     InfluxHostname,
		Username: InfluxUsername,
		Password: InfluxPassword,
	})

	if err != nil {
		log.Fatal(err)
	}

	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  InfluxDatabaseName,
		Precision: "s",
	})
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range activitySteps.Steps {
		t1, e := time.Parse("2006-01-02", v.Time)
		if e != nil {
			log.Fatal(e)
		}
		tags := map[string]string{"steps": "steps-total"}
		fields := map[string]interface{}{
			"steps": v.Value,
		}
		pt, err := client.NewPoint("activity_steps", tags, fields, t1)
		if err != nil {
			log.Fatal(err)
		}
		bp.AddPoint(pt)
	}

	// Write the batch
	if err := c.Write(bp); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Done loading steps")

	fmt.Println("Loading resting heartrate data")
	activityHeart := api.GetRestingHeartrate()

	bp, err2 := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  InfluxDatabaseName,
		Precision: "s",
	})
	if err2 != nil {
		log.Fatal(err)
	}

	for _, v := range activityHeart.HeartData {
		t1, e := time.Parse("2006-01-02", v.Date)
		if e != nil {
			log.Fatal(e)
		}
		tags := map[string]string{"heart": "resting-heart"}
		fields := map[string]interface{}{
			"resting": v.Value.RestingHeartRate,
		}
		pt, err := client.NewPoint("heart", tags, fields, t1)
		if err != nil {
			log.Fatal(err)
		}
		bp.AddPoint(pt)
	}

	// Write the batch
	if err := c.Write(bp); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Done")

	fmt.Println("Loading 30 days of 1s intraday heartrate data...")
	//Get Heart Rate Intraday Time Series
	now := time.Now()
	for i := 0; i < 30; i++ {
		dateString := now.AddDate(0, 0, -i).Format("2006-01-02")
		fmt.Printf("Loading: %s\n", dateString)
		series := api.GetHeartrateTimeSeries(dateString)

		bp, _ = client.NewBatchPoints(client.BatchPointsConfig{
			Database:  InfluxDatabaseName,
			Precision: "s",
		})

		for _, point := range series.GetNormalisedSeries("Pacific/Auckland") {
			tags := map[string]string{"heart": "intraday-heart"}
			fields := map[string]interface{}{
				"rate": point.Value,
			}
			pt, err := client.NewPoint("heart-intraday", tags, fields, point.Timestamp)
			if err != nil {
				log.Fatal(err)
			}
			bp.AddPoint(pt)
		}

		// Write the batch
		if err := c.Write(bp); err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println("Done")
}

func main() {
	mux := http.NewServeMux()
	api := fitbitapi.New(os.Getenv("FITBIT_CLIENT_ID"), os.Getenv("FITBIT_CLIENT_SECRET"), "http://localhost:4000/auth")

	mux.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query()["code"][0]
		api.LoadAccessToken(code)
		fmt.Fprintf(w, api.GetProfile())

		loadInfluxData(&api)
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html")
		fmt.Fprintf(w, "Visit: <a href=%q>%q</a>", api.AuthorizeUri, api.AuthorizeUri)
	})

	fmt.Println("Visit: " + api.AuthorizeUri)
	log.Fatal(http.ListenAndServe(":4000", mux))

}
