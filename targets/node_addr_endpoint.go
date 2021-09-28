package targets

import (
	"cosmos-validator-mission-control/config"
	"fmt"
	"log"

	client "github.com/influxdata/influxdb1-client/v2"
)

// CheckGaiad function to get gaiad status and send
//alerts to telgram and email accounts
func CheckGaiad(ops HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return
	}

	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Println("Error while getting gaiad status..")
		return
	}

	if (resp.StatusCode != 200) && (resp.StatusCode != 202) {
		_ = SendTelegramAlert(fmt.Sprintf("%s Gaiad on your validator instance is not running: RPC is DOWN : \n%v", cfg.ValidatorName, string(resp.Body)), cfg)
		_ = SendEmailAlert(fmt.Sprintf("%s Gaiad on your validator instance is not running: RPC is DOWN : \n%v", cfg.ValidatorName, string(resp.Body)), cfg)
		_ = writeToInfluxDb(c, bp, "vcf_gaiad_status", map[string]string{}, map[string]interface{}{"status": 0})
		return
	}

	_ = writeToInfluxDb(c, bp, "vcf_gaiad_status", map[string]string{}, map[string]interface{}{"status": 1})
}
