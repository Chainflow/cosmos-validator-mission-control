package targets

import (
	"cosmos-validator-mission-control/config"
	"encoding/json"
	"fmt"
	"log"

	client "github.com/influxdata/influxdb1-client/v2"
)

// GetLatestProposedBlockAndTime to get latest proposed block height and time
func GetLatestProposedBlockAndTime(ops HTTPOptions, cfg *config.Config, c client.Client) {
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

	var latestBlockTime string
	var proposerAddress string
	var blockHeight string

	if cfg.DaemonName == "akashd" {
		var blockResp AkashBlockInfo
		err = json.Unmarshal(resp.Body, &blockResp)
		if err != nil {
			log.Printf("Error: %v", err)
			return
		}

		latestBlockTime = blockResp.Block.Header.Time
		proposerAddress = blockResp.Block.Header.ProposerAddress
		blockHeight = blockResp.Block.Header.Height

	} else {
		var blockResp LastProposedBlockAndTime
		err = json.Unmarshal(resp.Body, &blockResp)
		if err != nil {
			log.Printf("Error: %v", err)
			return
		}

		latestBlockTime = blockResp.BlockMeta.Header.Time
		proposerAddress = blockResp.BlockMeta.Header.ProposerAddress
		blockHeight = blockResp.BlockMeta.Header.Height
	}

	log.Println("latest block details :: ", latestBlockTime, proposerAddress, blockHeight)

	blockTime := GetUserDateFormat(latestBlockTime)
	fmt.Println("last proposed block time", blockTime)

	if cfg.ValidatorHexAddress == proposerAddress {
		fields := map[string]interface{}{
			"height":     blockHeight,
			"block_time": blockTime,
		}

		_ = writeToInfluxDb(c, bp, "vcf_last_proposed_block", map[string]string{}, fields)
	}
}
