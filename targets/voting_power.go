package targets

import (
	"cosmos-validator-mission-control/config"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"strconv"

	client "github.com/influxdata/influxdb1-client/v2"
)

// GetValidatorVotingPower to get voting power of a validator
func GetValidatorVotingPower(ops HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error in voting power: %v", err)
		return
	}

	var validatorResp ValidatorResp
	err = json.Unmarshal(resp.Body, &validatorResp)
	if err != nil {
		log.Printf("Error while unamrshelling ValidatorResp of voting power: %v", err)
		return
	}

	v := validatorResp.Validator.DelegatorShares
	vp := convertValue(v)
	if vp == "" {
		vp = "0"
	}
	log.Printf("VOTING POWER: %s", vp)

	_ = writeToInfluxDb(c, bp, "vcf_voting_power", map[string]string{}, map[string]interface{}{"power": vp + "muon"})

	votingPower, err := strconv.Atoi(vp)
	if err != nil {
		log.Println("Error wile converting string to int of voting power \t", err)
		return
	}

	if int64(votingPower) <= cfg.VotingPowerThreshold {
		_ = SendTelegramAlert(fmt.Sprintf("Your validator %s voting power has dropped below %d", cfg.ValidatorName, cfg.VotingPowerThreshold), cfg)
		_ = SendEmailAlert(fmt.Sprintf("Your validator %s voting power has dropped below %d", cfg.ValidatorName, cfg.VotingPowerThreshold), cfg)
	}

}

func convertValue(value string) string {
	bal, _ := strconv.ParseFloat(value, 64)

	a1 := bal / math.Pow(10, 6)
	amount := strconv.FormatFloat(a1, 'f', -1, 64)

	return amount
}
