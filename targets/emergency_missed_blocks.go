package targets

import (
	"cosmos-validator-mission-control/config"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	client "github.com/influxdata/influxdb1-client/v2"
)

// EmergencyContinuousMissedBlocks is to send alerts to pager duty if validator miss blocks continously
func EmergencyContinuousMissedBlocks(ops HTTPOptions, cfg *config.Config, c client.Client) {
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
		Endpoint:    cfg.ExternalRPC + "/block",
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
		if c.ValidatorAddress == cfg.ValidatorHexAddress {
			addrExists = true
		}
	}

	if !addrExists {
		blocks := GetEmergencyContinuousMissedBlocks(cfg, c)
		blocksArray := strings.Split(blocks, ",")

		currentHeightFromDb := GetlatestCurrentHeightFromDB(cfg, c)
		if cfg.EmergencyMissedBlocksThreshold >= 2 {
			if int64(len(blocksArray))-1 >= cfg.EmergencyMissedBlocksThreshold {
				// Send emergency missed block alerts to telgram as well as pagerduty

				missedBlocks := strings.Split(blocks, ",")
				_ = SendTelegramAlert(fmt.Sprintf("%s validator missed blocks from height %s to %s", cfg.ValidatorName, missedBlocks[0], missedBlocks[len(missedBlocks)-2]), cfg)
				_ = SendEmailAlert(fmt.Sprintf("%s validator missed blocks from height %s to %s", cfg.ValidatorName, missedBlocks[0], missedBlocks[len(missedBlocks)-2]), cfg)
				_ = SendEmergencyEmailAlert(fmt.Sprintf("%s validator missed blocks from height %s to %s", cfg.ValidatorName, missedBlocks[0], missedBlocks[len(missedBlocks)-2]), cfg)
				_ = writeToInfluxDb(c, bp, "vcf_emergency_missed_blocks", map[string]string{}, map[string]interface{}{"block_height": "", "current_height": cbh})
				return
			} else if len(blocksArray) == 1 {
				blocks = cbh + ","
			} else {
				rpcBlockHeight, _ := strconv.Atoi(cbh)
				dbBlockHeight, _ := strconv.Atoi(currentHeightFromDb)
				diff := rpcBlockHeight - dbBlockHeight
				if diff == 1 {
					blocks = blocks + cbh + ","
				} else if diff > 1 {
					blocks = ""
				}
			}
			_ = writeToInfluxDb(c, bp, "vcf_emergency_missed_blocks", map[string]string{}, map[string]interface{}{"block_height": blocks, "current_height": cbh})
		}
	}
	return
}

// GetlatestCurrentHeightFromDB returns latest current height from db
func GetlatestCurrentHeightFromMissedBlocks(cfg *config.Config, c client.Client) string {
	var currentHeight string
	q := client.NewQuery("SELECT last(current_height) FROM vcf_emergency_missed_blocks", cfg.InfluxDB.Database, "")
	if response, err := c.Query(q); err == nil && response.Error() == nil {
		for _, r := range response.Results {
			if len(r.Series) != 0 {
				for idx, col := range r.Series[0].Columns {
					if col == "last" {
						heightValue := r.Series[0].Values[0][idx]
						currentHeight = fmt.Sprintf("%v", heightValue)
						break
					}
				}
			}
		}
	}
	return currentHeight
}

// GetContinuousMissedBlock returns the latest missed block from the db
func GetEmergencyContinuousMissedBlocks(cfg *config.Config, c client.Client) string {
	var blocks string
	q := client.NewQuery("SELECT last(block_height) FROM vcf_emergency_missed_blocks", cfg.InfluxDB.Database, "")
	if response, err := c.Query(q); err == nil && response.Error() == nil {
		for _, r := range response.Results {
			if len(r.Series) != 0 {
				for idx, col := range r.Series[0].Columns {
					if col == "last" {
						heightValue := r.Series[0].Values[0][idx]
						blocks = fmt.Sprintf("%v", heightValue)
						break
					}
				}
			}
		}
	}
	return blocks
}
