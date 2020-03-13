package targets

import (
	"chainflow-vitwit/config"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	client "github.com/influxdata/influxdb1-client/v2"
)

// Get voting power of a validator
func GetValidatorVotingPower(ops HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	// Calling function to get current block height
	currentHeight := GetValidatorBlock(cfg, c)

	ops.Endpoint = ops.Endpoint + "?height=" + currentHeight

	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	var validatorHeightResp ValidatorsHeight
	err = json.Unmarshal(resp.Body, &validatorHeightResp)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	for _, val := range validatorHeightResp.Result.Validators {
		if val.Address == cfg.ValidatorAddress {
			var vp string
			fmt.Printf("VOTING POWER: %s\n", val.VotingPower)
			if val.VotingPower != "" {
				vp = val.VotingPower
			} else {
				vp = "0"
			}
			_ = writeToInfluxDb(c, bp, "vcf_voting_power", map[string]string{}, map[string]interface{}{"power": vp + "muon"})
			log.Println("Voting Power \n", vp)

			votingPower, err := strconv.Atoi(vp)
			if err != nil {
				log.Println("Error wile converting string to int of voting power \t", err)
			}

			if int64(votingPower) <= cfg.VotingPowerThreshold {
				_ = SendTelegramAlert(fmt.Sprintf("Your validator's voting power has dropped below %d", cfg.VotingPowerThreshold), cfg)
				_ = SendEmailAlert(fmt.Sprintf("Your validator's voting power has dropped below %d", cfg.VotingPowerThreshold), cfg)
			}
		}
	}

}
