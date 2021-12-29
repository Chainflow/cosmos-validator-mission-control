package server

import (
	"cosmos-validator-mission-control/alert-bot/config"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func GetNetworkBlock(cfg *config.Config) error {
	ops := HTTPOptions{
		Endpoint: cfg.ExternalRPC + "/status?",
		Method:   http.MethodGet,
	}

	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error: %v", err)
		return err
	}

	var networkBlock NetworkLatestBlock
	err = json.Unmarshal(resp.Body, &networkBlock)
	if err != nil {
		log.Printf("Error: %v", err)
		return err
	}

	if &networkBlock != nil {

		networkBlockHeight, err := strconv.Atoi(networkBlock.Result.SyncInfo.LatestBlockHeight)
		if err != nil {
			log.Println("Error while converting network height from string to int ", err)
		}

		// Calling function to get validator latest
		// block height
		validatorHeight, err := GetValBlockHeight(cfg)
		if validatorHeight == "" {
			log.Println("Error while fetching validator block height from db ", validatorHeight)
			return err
		}

		vaidatorBlockHeight, _ := strconv.Atoi(validatorHeight)
		heightDiff := networkBlockHeight - vaidatorBlockHeight

		log.Printf("Network height: %d and Validator Height: %d", networkBlockHeight, vaidatorBlockHeight)

		// Send alert
		if int64(heightDiff) >= cfg.BlockDiffThreshold {
			_ = SendTelegramAlert(fmt.Sprintf("%s Block difference between network and validator has exceeded %d", cfg.ValidatorName, cfg.BlockDiffThreshold), cfg)
			_ = SendEmailAlert(fmt.Sprintf("%s Block difference between network and validator has exceeded %d", cfg.ValidatorName, cfg.BlockDiffThreshold), cfg)

			log.Println("Sent alert of block height difference")
		}

	}

	return nil

}

func GetValBlockHeight(cfg *config.Config) (string, error) {
	ops := HTTPOptions{
		Endpoint: cfg.ValidatorRpcEndpoint + "/status?",
		Method:   http.MethodGet,
	}

	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error: %v", err)
		return "", err
	}

	var status ValidatorRpcStatus
	err = json.Unmarshal(resp.Body, &status)
	if err != nil {
		log.Printf("Error: %v", err)
		return "", nil
	}

	currentBlockHeight := status.Result.SyncInfo.LatestBlockHeight
	return currentBlockHeight, err
}
