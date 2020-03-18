package targets

import (
	"chainflow-vitwit/config"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	client "github.com/influxdata/influxdb1-client/v2"
)

//SendSingleMissedBlockAlert Send missed block alert to telegram bot and mail
func SendSingleMissedBlockAlert(ops HTTPOptions, cfg *config.Config, c client.Client) error {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return err
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
		if cfg.MissedBlocksThreshold == 1 {
			_ = SendTelegramAlert(fmt.Sprintf("Validator missed a block at block height %s", cbh), cfg)
			_ = SendEmailAlert(fmt.Sprintf("Validator missed a block at block height %s", cbh), cfg)
			_ = writeToInfluxDb(c, bp, "vcf_continuous_missed_blocks", map[string]string{}, map[string]interface{}{"missed_blocks": cbh, "range": cbh})
			_ = writeToInfluxDb(c, bp, "vcf_missed_blocks", map[string]string{}, map[string]interface{}{"block_height": cbh, "current_height": cbh})
			_ = writeToInfluxDb(c, bp, "vcf_total_missed_blocks", map[string]string{}, map[string]interface{}{"block_height": cbh, "current_height": cbh})
		} else {
			_ = writeToInfluxDb(c, bp, "vcf_total_missed_blocks", map[string]string{}, map[string]interface{}{"block_height": cbh, "current_height": cbh})
		}

	}
	return nil
}

//GetMissedBlocks to get missed block and send alert
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
		blocks := GetContinuousMissedBlock(cfg, c)
		currentHeightFromDb := GetlatestCurrentHeightFromDB(cfg, c)
		blocksArray := strings.Split(blocks, ",")
		fmt.Println("blocks length ", int64(len(blocksArray)), currentHeightFromDb)
		// calling function to store single blocks
		err := SendSingleMissedBlockAlert(ops, cfg, c)
		if err != nil {
			log.Printf("Error while sending missed block alert: %v", err)

		}
		if cfg.MissedBlocksThreshold > 1 {
			if int64(len(blocksArray))-1 >= cfg.MissedBlocksThreshold {
				missedBlocks := strings.Split(blocks, ",")
				_ = SendTelegramAlert(fmt.Sprintf("Validator missed blocks from height %s to %s", missedBlocks[0], missedBlocks[len(missedBlocks)-2]), cfg)
				_ = SendEmailAlert(fmt.Sprintf("Validator missed blocks from height %s to %s", missedBlocks[0], missedBlocks[len(missedBlocks)-2]), cfg)
				_ = writeToInfluxDb(c, bp, "vcf_continuous_missed_blocks", map[string]string{}, map[string]interface{}{"missed_blocks": blocks, "range": missedBlocks[0] + " - " + missedBlocks[len(missedBlocks)-2]})
				_ = writeToInfluxDb(c, bp, "vcf_missed_blocks", map[string]string{}, map[string]interface{}{"block_height": "", "current_height": cbh})
				return
			}
			if len(blocksArray) == 1 {
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
			_ = writeToInfluxDb(c, bp, "vcf_missed_blocks", map[string]string{}, map[string]interface{}{"block_height": blocks, "current_height": cbh})
			return

		}
	}
	return
}

//GetContinuousMissedBlock returns the latest missed block
func GetContinuousMissedBlock(cfg *config.Config, c client.Client) string {
	var blocks string
	q := client.NewQuery("SELECT last(block_height) FROM vcf_missed_blocks", cfg.InfluxDB.Database, "")
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

//GetlatestCurrentHeightFromDB returns latest current height from db
func GetlatestCurrentHeightFromDB(cfg *config.Config, c client.Client) string {
	var currentHeight string
	q := client.NewQuery("SELECT last(current_height) FROM vcf_missed_blocks", cfg.InfluxDB.Database, "")
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
