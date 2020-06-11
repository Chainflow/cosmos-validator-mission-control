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

	var applicationInfo ApplicationInfo
	err = json.Unmarshal(resp.Body, &applicationInfo)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	gaiaVersion := applicationInfo.ApplicationVersion.Version

	_ = writeToInfluxDb(c, bp, "vcf_gaiad_version", map[string]string{}, map[string]interface{}{"v": gaiaVersion})
	log.Printf("Version: %s", gaiaVersion)
}
