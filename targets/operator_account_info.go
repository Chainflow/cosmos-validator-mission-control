package targets

import (
	"chainflow-vitwit/config"
	"encoding/json"
	client "github.com/influxdata/influxdb1-client/v2"
	"log"
)

func GetAccountInfo(ops HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return
	}

	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	var accResp AccountResp
	err = json.Unmarshal(resp.Body, &accResp)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	addressBalance := accResp.Account.Balance[0].Amount + accResp.Account.Balance[0].Denom
	_ = writeToInfluxDb(c, bp, "vcf_account_balance", map[string]string{}, map[string]interface{}{"balance": addressBalance})
	log.Printf("Address Balance: %s", addressBalance)
}
