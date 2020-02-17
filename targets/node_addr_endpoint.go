package targets

import (
	"chainflow-vitwit/config"
)

func GetNodeAddrEndpointData(ops HTTPOptions, cfg *config.Config) {
	resp, err := HitHTTPTarget(ops)
	if err != nil {
		_ = SendTelegramAlert("Gaiad is not running", cfg)
		return
	}

	if resp.StatusCode == 200 {
		return
	}
	_ = SendTelegramAlert("Gaiad is not running", cfg)
}
