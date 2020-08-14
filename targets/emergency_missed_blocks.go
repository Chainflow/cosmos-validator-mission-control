package targets

import (
	"cosmos-validator-mission-control/config"
	"fmt"
	"strconv"
	"strings"

	client "github.com/influxdata/influxdb1-client/v2"
)

// SendEmeregencyAlerts is to send alerts to pager duty if validator miss blocks continously
func SendEmeregencyAlerts(cfg *config.Config, c client.Client, cbh string) error {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return err
	}

	blocks := GetEmergencyContinuousMissedBlocks(cfg, c)
	blocksArray := strings.Split(blocks, ",")

	currentHeightFromDb := GetlatestCurrentHeightFromMissedBlocks(cfg, c)
	if cfg.EmergencyMissedBlocksThreshold >= 2 {
		if int64(len(blocksArray))-1 >= cfg.EmergencyMissedBlocksThreshold {
			// Send emergency missed block alerts to telgram as well as pagerduty

			missedBlocks := strings.Split(blocks, ",")
			_ = SendTelegramAlert(fmt.Sprintf("%s validator missed blocks from height %s to %s", cfg.ValidatorName, missedBlocks[0], missedBlocks[len(missedBlocks)-2]), cfg)
			_ = SendEmailAlert(fmt.Sprintf("%s validator missed blocks from height %s to %s", cfg.ValidatorName, missedBlocks[0], missedBlocks[len(missedBlocks)-2]), cfg)
			_ = SendEmergencyEmailAlert(fmt.Sprintf("%s validator missed blocks from height %s to %s", cfg.ValidatorName, missedBlocks[0], missedBlocks[len(missedBlocks)-2]), cfg)
			_ = writeToInfluxDb(c, bp, "vcf_emergency_missed_blocks", map[string]string{}, map[string]interface{}{"block_height": "", "current_height": cbh})
			return nil
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

	return nil
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
