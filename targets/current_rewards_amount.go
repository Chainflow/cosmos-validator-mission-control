package targets

import (
	"cosmos-validator-mission-control/config"
	"encoding/json"
	"log"
	"net/http"
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

	var rewardsResp Rewards
	err = json.Unmarshal(resp.Body, &rewardsResp)
	if err != nil {
		log.Printf("Error while unmarshalling current rewards: %v", err)
		return
	}

	var oustandingRewards float64
	var denom string

	if len(rewardsResp.Rewards.Rewards) > 0 {
		f, _ := strconv.ParseFloat(rewardsResp.Rewards.Rewards[0].Amount, 64)
		oustandingRewards = f
		denom = rewardsResp.Rewards.Rewards[0].Denom
	}

	commission := GetValCommission(ops, cfg, c)
	if oustandingRewards != 0 && commission != 0 {
		rewards := oustandingRewards - commission
		num := int(rewards)
		amount := strconv.Itoa(num)
		r := convertToCommaSeparated(amount) + denom
		_ = writeToInfluxDb(c, bp, "vcf_current_rewards_amount", map[string]string{}, map[string]interface{}{"amount": r})
		log.Printf("Current Rewards: %s", r)
	}

	log.Printf("Commission : %f and Outstanding Rewards : %f", commission, oustandingRewards)
}

// GetValCommission which return the commission of a validator
func GetValCommission(ops HTTPOptions, cfg *config.Config, c client.Client) float64 {
	ops = HTTPOptions{
		Endpoint: cfg.LCDEndpoint + "/cosmos/distribution/v1beta1/validators/" + cfg.ValOperatorAddress + "/commission",
		Method:   http.MethodGet,
	}

	var commission float64

	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error: %v", err)
		return commission
	}

	var result Commission
	err = json.Unmarshal(resp.Body, &result)
	if err != nil {
		log.Printf("Error while unmarshalling commission: %v", err)
		return commission
	}

	if len(result.Commission.Commission) > 0 {
		f, _ := strconv.ParseFloat(result.Commission.Commission[0].Amount, 64)
		commission = f

	}

	return commission
}
