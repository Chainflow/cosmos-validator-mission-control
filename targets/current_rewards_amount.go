package targets

import (
	"cosmos-validator-mission-control/config"
	"encoding/json"
	"log"
	"strconv"

	client "github.com/influxdata/influxdb1-client/v2"
)

// GetCurrentRewardsAmount to get current rewards of a validator using operator address
func GetCurrentRewardsAmount(ops HTTPOptions, cfg *config.Config, c client.Client) {
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

	var rewardsResp CurrentRewardsAmount
	err = json.Unmarshal(resp.Body, &rewardsResp)
	if err != nil {
		log.Printf("Error while unmarshalling current rewards: %v", err)
		return
	}

	if len(rewardsResp.Result) > 0 {
		f, _ := strconv.ParseFloat(rewardsResp.Result[0].Amount, 64)
		num := int(f)
		amount := strconv.Itoa(num)
		addressBalance := convertToCommaSeparated(amount) + rewardsResp.Result[0].Denom
		_ = writeToInfluxDb(c, bp, "vcf_current_rewards_amount", map[string]string{}, map[string]interface{}{"amount": addressBalance})
		log.Printf("Current Rewards: %s", addressBalance)
	}
}
