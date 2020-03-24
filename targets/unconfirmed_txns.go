package targets

import (
	"chainflow-vitwit/config"
	"encoding/json"
	"log"

	client "github.com/influxdata/influxdb1-client/v2"
)

//GetUnconfimedTxns to get the no of txns
func GetUnconfimedTxns(ops HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	var unconfirmedTxns UnconfirmedTxns
	err = json.Unmarshal(resp.Body, &unconfirmedTxns)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	totalUnconfirmedTxns := unconfirmedTxns.Result.Total

	_ = writeToInfluxDb(c, bp, "vcf_unconfirmed_txns", map[string]string{}, map[string]interface{}{"unconfirmed_txns": totalUnconfirmedTxns})
	log.Println("No of unconfirmed txns: ", totalUnconfirmedTxns)
}
