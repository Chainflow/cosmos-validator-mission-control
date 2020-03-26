package server

import (
	"chainflow-vitwit/alert-bot/config"
	"encoding/json"
	"fmt"
	"log"
)

// SendSingleMissedBlockAlert to send missed block alerting
func SendSingleMissedBlockAlert(cfg *config.Config) error {
	log.Println("Calling missed block alerting")
	ops := HTTPOptions{
		Endpoint: cfg.ExternalRPC + "status",
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
		log.Printf("Error: %v", err)
		return err
	}

	cbh := networkLatestBlock.Result.SyncInfo.LatestBlockHeight

	resp, err = HitHTTPTarget(HTTPOptions{
		Endpoint:    cfg.ExternalRPC + "block",
		QueryParams: QueryParams{"height": cbh},
		Method:      "GET",
	})
	if err != nil {
		log.Printf("Error getting details of current block: %v", err)
		return err
	}

	var b CurrentBlockWithHeight
	err = json.Unmarshal(resp.Body, &b)
	if err != nil {
		log.Printf("Error: %v", err)
		return err
	}

	addrExists := false

	for _, c := range b.Result.Block.LastCommit.Precommits {
		if c.ValidatorAddress == cfg.ValidatorAddress {
			addrExists = true
		}
	}

	if !addrExists {
		_ = SendTelegramAlert(fmt.Sprintf("Validator missed a block at block height %s", cbh), cfg)
		_ = SendEmailAlert(fmt.Sprintf("Validator missed a block at block height %s", cbh), cfg)
		log.Println("Sent missed block alerting")
	}

	// Calling function to check validator jailed status
	err = JailedTxAlerting(cfg)
	if err != nil {
		log.Printf("Error while sending jailed alerting: %v", err)
		return err
	}

	return nil
}
