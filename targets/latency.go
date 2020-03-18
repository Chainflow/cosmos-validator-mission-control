package targets

import (
	"chainflow-vitwit/config"
	"fmt"
	"log"
	"strings"

	client "github.com/influxdata/influxdb1-client/v2"
	ping "github.com/sparrc/go-ping"
)

// GetLatency to calculate latency
func GetLatency(_ HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return
	}

	q := client.NewQuery(fmt.Sprintf("SELECT * FROM vcf_peer_addresses"), cfg.InfluxDB.Database, "")
	if response, err := c.Query(q); err == nil && response.Error() == nil {
		var addresses []string
		for _, r := range response.Results {
			if len(r.Series) != 0 {
				noOfValues := len(r.Series[0].Values)
				if noOfValues != 0 {
					n := noOfValues - 1
					addressValues := fmt.Sprintf("%v", r.Series[0].Values[n][1])
					addresses = strings.Split(addressValues, ", ")
				}
			}
		}
		for _, addr := range addresses {
			log.Printf("peer address %s", addr)
			pinger, err := ping.NewPinger(addr)
			if err != nil {
				log.Printf("Error while getting ping response of %s", addr)
				return
			}
			pinger.Count = 10
			pinger.Run() // blocks until finished
			stats := pinger.Statistics()

			fields := map[string]interface{}{
				"Address":    stats.Addr,
				"PacketLoss": stats.PacketLoss,
				"PacketSent": stats.PacketsSent,
			}

			log.Printf("Writing address latency in db %s", addr)
			_ = writeToInfluxDb(c, bp, "vcf_validator_latency", map[string]string{}, fields)
		}
	}
}
