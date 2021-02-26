package targets

import (
	"cosmos-validator-mission-control/config"
	"encoding/json"
	"log"
	"strconv"

	client "github.com/influxdata/influxdb1-client/v2"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func convertToCommaSeparated(amt string) string {
	a, err := strconv.Atoi(amt)
	if err != nil {
		return amt
	}
	p := message.NewPrinter(language.English)
	return p.Sprintf("%d", a)
}

// GetAccountInfo to get account balance information using account address
func GetAccountInfo(ops HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return
	}

	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error: %v", err)
		_ = writeToInfluxDb(c, bp, "vcf_account_balance", map[string]string{}, map[string]interface{}{"balance": "NA"})
		return
	}

	var accResp AccountResp
	err = json.Unmarshal(resp.Body, &accResp)
	if err != nil {
		log.Printf("Error while unmarshalling bank balances res: %v", err)
		_ = writeToInfluxDb(c, bp, "vcf_account_balance", map[string]string{}, map[string]interface{}{"balance": "NA"})
		return
	}

	if len(accResp.Balances) > 0 {
		addressBalance := convertToCommaSeparated(accResp.Balances[0].Amount) + accResp.Balances[0].Denom
		_ = writeToInfluxDb(c, bp, "vcf_account_balance", map[string]string{}, map[string]interface{}{"balance": addressBalance})
		log.Printf("Address Balance: %s", addressBalance)
	}
}
