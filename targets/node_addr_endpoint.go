package targets

import (
	"chainflow-vitwit/config"
	"fmt"
	"os/exec"

	client "github.com/influxdata/influxdb1-client/v2"
)

// CheckGaiad function to run the command get gaiad status and send
//alerts to telgram and email accounts
func CheckGaiad(_ HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return
	}

	cmd := exec.Command("bash", "-c", "</dev/tcp/0.0.0.0/26656 &>/dev/null")
	out, err := cmd.CombinedOutput()
	if err != nil {
		_ = SendTelegramAlert("Gaiad on your validator instance is not running", cfg)
		_ = SendEmailAlert("Gaiad on your validator instance is not running", cfg)
		_ = writeToInfluxDb(c, bp, "vcf_gaiad_status", map[string]string{}, map[string]interface{}{"status": 0})
		return
	}

	resp := string(out)
	if resp != "" {
		_ = SendTelegramAlert(fmt.Sprintf("Gaiad on your validator instance is not running: \n%v", resp), cfg)
		_ = SendEmailAlert(fmt.Sprintf("Gaiad on your validator instance is not running: \n%v", resp), cfg)
		_ = writeToInfluxDb(c, bp, "vcf_gaiad_status", map[string]string{}, map[string]interface{}{"status": 0})
		return
	}

	_ = writeToInfluxDb(c, bp, "vcf_gaiad_status", map[string]string{}, map[string]interface{}{"status": 1})
}
