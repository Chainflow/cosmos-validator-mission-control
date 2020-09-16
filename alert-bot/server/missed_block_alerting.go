package server

import (
	"cosmos-validator-mission-control/alert-bot/config"
	"encoding/json"
	"fmt"
	"log"

	client "github.com/influxdata/influxdb1-client/v2"
)

// SendSingleMissedBlockAlert to send missed block alerting
func SendSingleMissedBlockAlert(ops HTTPOptions, cfg *config.Config, c client.Client) {
	log.Println("Calling missed block alerting")

	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	var networkLatestBlock NetworkLatestBlock
	err = json.Unmarshal(resp.Body, &networkLatestBlock)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	cbh := networkLatestBlock.Result.SyncInfo.LatestBlockHeight

	resp, err = HitHTTPTarget(HTTPOptions{
		Endpoint:    cfg.ExternalRPC + "/block",
		QueryParams: QueryParams{"height": cbh},
		Method:      "GET",
	})
	if err != nil {
		log.Printf("Error getting details of current block: %v", err)
		return
	}

	var b BlockResponse
	err = json.Unmarshal(resp.Body, &b)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	addrExists := false

	for _, c := range b.Result.Block.LastCommit.Precommits {
		if c.ValidatorAddress == cfg.ValidatorHexAddress {
			addrExists = true
		}
	}

	if !addrExists {
		_ = SendTelegramAlert(fmt.Sprintf("%s validator missed a block at block height %s", cfg.ValidatorName, cbh), cfg)
		_ = SendEmailAlert(fmt.Sprintf("%s validator missed a block at block height %s", cfg.ValidatorName, cbh), cfg)
		log.Println("Sent missed block alerting")
	}

	// Calling function to check validator jailed status
	err = CheckValidatorJailed(cfg)
	if err != nil {
		log.Printf("Error while sending jailed alerting: %v", err)
		return
	}

	return
}
