package targets

import (
	"chainflow-vitwit/config"
	"encoding/json"
	"log"

	client "github.com/influxdata/influxdb1-client/v2"
)

func GetLatProposedBlockAndTime(ops HTTPOptions, cfg *config.Config, c client.Client) {
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

	var blockResp LastProposedBlockAndTime
	err = json.Unmarshal(resp.Body, &blockResp)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	if cfg.ValidatorAddress == blockResp.BlockMeta.Header.ProposerAddress {
		fields := map[string]interface{}{
			"height":     blockResp.BlockMeta.Header.Height,
			"block_time": blockResp.BlockMeta.Header.Time,
		}

		_ = writeToInfluxDb(c, bp, "vcf_last_proposed_block", map[string]string{}, fields)
	}
}
