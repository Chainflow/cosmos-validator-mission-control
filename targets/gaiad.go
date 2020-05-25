package targets

import (
	"cosmos-validator-mission-control/config"
	"encoding/json"
	"log"

	client "github.com/influxdata/influxdb1-client/v2"
)

// GaiadVersion to get gaiad version by running command gaiad version
func GaiadVersion(ops HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return
	}

	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	var statusInfo GaiadStatusInfo
	err = json.Unmarshal(resp.Body, &statusInfo)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	version := statusInfo.Result.NodeInfo.Version

	_ = writeToInfluxDb(c, bp, "vcf_gaiad_version", map[string]string{}, map[string]interface{}{"v": version})
	log.Printf("Version: %s", version)
}
