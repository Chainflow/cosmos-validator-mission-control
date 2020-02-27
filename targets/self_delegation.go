package targets

import (
	"chainflow-vitwit/config"
	"encoding/json"
	"log"

	client "github.com/influxdata/influxdb1-client/v2"
)

func GetSelfDelegation(ops HTTPOptions, cfg *config.Config, c client.Client) {
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

	var delegationResp SelfDelegation
	err = json.Unmarshal(resp.Body, &delegationResp)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	addressBalance := convertToCommaSeparated(delegationResp.Result.Balance) + "uatoms"
	_ = writeToInfluxDb(c, bp, "vcf_self_delegation_balance", map[string]string{}, map[string]interface{}{"balance": addressBalance})
	log.Printf("Address Balance: %s", addressBalance)	
}