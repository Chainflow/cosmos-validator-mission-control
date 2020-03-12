package targets

import (
	"chainflow-vitwit/config"
	"encoding/json"
	"fmt"
	"log"

	client "github.com/influxdata/influxdb1-client/v2"
)

// Function to get get missed block and send alert
func GetMissedBlocks(ops HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return
	}

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
		Endpoint:    cfg.ExternalRPC + "block",
		QueryParams: QueryParams{"height": cbh},
		Method:      "GET",
	})
	if err != nil {
		log.Printf("Error getting details of current block: %v", err)
		return
	}

	var b CurrentBlockWithHeight
	err = json.Unmarshal(resp.Body, &b)
	if err != nil {
		log.Printf("Error: %v", err)
		return
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
		_ = writeToInfluxDb(c, bp, "vcf_missed_blocks", map[string]string{}, map[string]interface{}{"block_height": cbh})
	}
}
