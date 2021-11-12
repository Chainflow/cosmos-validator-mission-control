package server

import (
	"cosmos-validator-mission-control/alert-bot/config"
	"encoding/json"
	"fmt"
	"log"
)

// SendSingleMissedBlockAlert to send missed block alerting
func SendSingleMissedBlockAlert(cfg *config.Config) error {
	log.Println("Calling missed block alerting")
	ops := HTTPOptions{
		Endpoint: cfg.ExternalRPC + "/status",
		Method:   "GET",
	}

	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error: %v", err)
		return err
	}

	var networkLatestBlock NetworkLatestBlock
	err = json.Unmarshal(resp.Body, &networkLatestBlock)
	if err != nil {
		log.Printf("Error while unmarshelling n/w status res : %v", err)
		return err
	}

	cbh := networkLatestBlock.Result.SyncInfo.LatestBlockHeight

	resp, err = HitHTTPTarget(HTTPOptions{
		Endpoint:    cfg.ExternalRPC + "/block",
		QueryParams: QueryParams{"height": cbh},
		Method:      "GET",
	})
	if err != nil {
		log.Printf("Error getting details of current block: %v", err)
		return err
	}

	var b CurrentBlock
	err = json.Unmarshal(resp.Body, &b)
	if err != nil {
		log.Printf("Error while unmarshelling block res : %v", err)
		return err
	}

	if &b != nil {

		addrExists := false

		for _, c := range b.Result.Block.LastCommit.Signatures {
			if c.ValidatorAddress == cfg.ValidatorHexAddress {
				addrExists = true
			}
		}

		if !addrExists {
			_ = SendTelegramAlert(fmt.Sprintf("%s validator missed a block at height %s", cfg.ValidatorName, cbh), cfg)
			_ = SendEmailAlert(fmt.Sprintf("%s validator missed a block at height %s", cfg.ValidatorName, cbh), cfg)
		}

		// Calling function to check validator jailed status
		err = CheckValidatorJailed(cfg)
		if err != nil {
			log.Printf("Error while sending jailed alerting: %v", err)
			return err
		}
	}

	return nil
}
