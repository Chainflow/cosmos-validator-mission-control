package targets

import (
	"chainflow-vitwit/config"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	client "github.com/influxdata/influxdb1-client/v2"
)

// GetNetworkLatestBlock to get latest block height of a network
func GetNetworkLatestBlock(ops HTTPOptions, cfg *config.Config, c client.Client) {
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

	var networkBlock NetworkLatestBlock
	err = json.Unmarshal(resp.Body, &networkBlock)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	networkBlockHeight, err := strconv.Atoi(networkBlock.Result.SyncInfo.LatestBlockHeight)
	if err != nil {
		log.Println("Error while converting network height from string to int ", err)
	}
	_ = writeToInfluxDb(c, bp, "vcf_network_latest_block", map[string]string{}, map[string]interface{}{"block_height": networkBlockHeight})
	log.Printf("Network height: %d", networkBlockHeight)

	// Calling function to get validator latest
	// block height
	validatorHeight := GetValidatorBlock(cfg, c)

	vaidatorBlockHeight, _ := strconv.Atoi(validatorHeight)
	heightDiff := networkBlockHeight - vaidatorBlockHeight

	_ = writeToInfluxDb(c, bp, "vcf_height_difference", map[string]string{}, map[string]interface{}{"difference": heightDiff})
	log.Printf("Network height: %d", networkBlockHeight)
}

// GetValidatorBlock returns validator current block height
func GetValidatorBlock(cfg *config.Config, c client.Client) string {
	var validatorHeight string
	q := client.NewQuery("SELECT last(height) FROM vcf_current_block_height", cfg.InfluxDB.Database, "")
	if response, err := c.Query(q); err == nil && response.Error() == nil {
		for _, r := range response.Results {
			if len(r.Series) != 0 {
				for idx, col := range r.Series[0].Columns {
					if col == "last" {
						heightValue := r.Series[0].Values[0][idx]
						validatorHeight = fmt.Sprintf("%v", heightValue)
						break
					}
				}
			}
		}
	}

	return validatorHeight
}
