package main

import (
	"cosmos-validator-mission-control/alert-bot/config"
	"cosmos-validator-mission-control/alert-bot/server"
	"fmt"
	"log"
	"sync"
	"time"

	client "github.com/influxdata/influxdb1-client/v2"
)

func main() {
	cfg, err := config.ReadConfigFromFile()
	if err != nil {
		log.Fatal(err)
	}

	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     fmt.Sprintf("http://localhost:%s", cfg.InfluxDB.Port),
		Username: cfg.InfluxDB.Username,
		Password: cfg.InfluxDB.Password,
	})

	var wg sync.WaitGroup
	wg.Add(1)

	// Calling go routine to send alerts for missed blocks
	go func() {
		for {
			if err := server.GetMissedBlocks(cfg, c); err != nil {
				fmt.Println("Error while sending missed block alerts", err)
			}
			time.Sleep(4 * time.Second)
		}
	}()

	// Calling go routine to send alert about validator status
	go func() {
		for {
			if err := server.ValidatorStatusAlert(cfg); err != nil {
				fmt.Println("Error while sending jailed alerts", err)
			}
			time.Sleep(60 * time.Second)
		}
	}()

	wg.Wait()
}
