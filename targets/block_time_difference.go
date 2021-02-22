package targets

import (
	"cosmos-validator-mission-control/config"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	client "github.com/influxdata/influxdb1-client/v2"
)

// GetBlockTimeDifference to calculate block time difference of prev block and current block
func GetBlockTimeDifference(ops HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	// Calling function to get validator latest
	// block height
	currentBlockHeight := GetValidatorBlock(cfg, c)
	if currentBlockHeight == "" {
		log.Println("Error while fetching current block height from db ", currentBlockHeight)
		return
	}

	ops.Endpoint = ops.Endpoint + "?height=" + currentBlockHeight
	currResp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	var currentBlockResp CurrentBlockWithHeight
	err = json.Unmarshal(currResp.Body, &currentBlockResp)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	currentBlockTime := currentBlockResp.Result.Block.Header.Time
	currentHeight, _ := strconv.Atoi(currentBlockHeight)

	prevHeight := currentHeight - 1
	ops.Endpoint = cfg.ValidatorRpcEndpoint + "/block"
	ops.Endpoint = ops.Endpoint + "?height=" + strconv.Itoa(prevHeight)

	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	var prevBlockResp CurrentBlockWithHeight
	err = json.Unmarshal(resp.Body, &prevBlockResp)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	prevBlockTime := prevBlockResp.Result.Block.Header.Time
	convertedCurrentTime, _ := time.Parse(time.RFC3339, currentBlockTime)
	conevrtedPrevBlockTime, _ := time.Parse(time.RFC3339, prevBlockTime)
	timeDiff := convertedCurrentTime.Sub(conevrtedPrevBlockTime)
	diffSeconds := fmt.Sprintf("%.2f", timeDiff.Seconds())

	_ = writeToInfluxDb(c, bp, "vcf_block_time_diff", map[string]string{}, map[string]interface{}{"time_diff": diffSeconds})
	log.Printf("time diff: %s", diffSeconds)

}
