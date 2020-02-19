package targets

import (
	"chainflow-vitwit/config"
	"encoding/json"
	client "github.com/influxdata/influxdb1-client/v2"
	"log"
)

func GetAccountInfo(ops HTTPOptions, cfg *config.Config, c client.Client) {
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

	log.Printf("Address Balance: %s", addressBalance)
}
