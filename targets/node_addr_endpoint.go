package targets

import (
	"chainflow-vitwit/config"
	"fmt"
	"os/exec"
)

func CheckGaiad(_ HTTPOptions, cfg *config.Config) {
	cmd := exec.Command("bash", "-c", "</dev/tcp/0.0.0.0/26656 &>/dev/null")
	out, err := cmd.CombinedOutput()
	if err != nil {
		_ = SendTelegramAlert("Gaiad is not running", cfg)
		_ = SendEmailAlert("Gaiad is not running", cfg)
		GaiadRunningGauge.Set(float64(0))
		return
	}

	resp := string(out)
	if resp != "" {
		_ = SendTelegramAlert(fmt.Sprintf("Gaiad is not running: \n%v", resp), cfg)
		_ = SendEmailAlert(fmt.Sprintf("Gaiad is not running: \n%v", resp), cfg)
		GaiadRunningGauge.Set(float64(0))
		return
	}
	GaiadRunningGauge.Set(float64(1))
}
