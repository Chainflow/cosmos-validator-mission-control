package server

type (
	// QueryParams to map the query params of an url
	QueryParams map[string]string

	// HTTPOptions of a target
	HTTPOptions struct {
		Endpoint    string
		QueryParams QueryParams
		Body        []byte
		Method      string
	}

	// PingResp struct
	PingResp struct {
		StatusCode int
		Body       []byte
	}

	// NetworkLatestBlock stores latest block height info
	NetworkLatestBlock struct {
		Result struct {
			SyncInfo struct {
				LatestBlockHeight string `json:"latest_block_height"`
			} `json:"sync_info"`
		} `json:"result"`
	}

	// CurrentBlock is a struct which holds latest current block details
	CurrentBlock struct {
		Result struct {
			Block struct {
				LastCommit struct {
					Signatures []struct {
						ValidatorAddress string `json:"validator_address"`
						Signature        string `json:"signature"`
					} `json:"signatures"`
				} `json:"last_commit"`
			} `json:"block"`
		} `json:"result"`
	}

	// Validator is a struct which holds validator staking info
	Validator struct {
		Validator struct {
			Jailed bool   `json:"jailed"`
			Status string `json:"status"`
		} `json:"validator"`
	}

	// ValidatorRpcStatus is a struct which holds the status response
	ValidatorRpcStatus struct {
		Jsonrpc string `json:"jsonrpc"`
		Result  struct {
			SyncInfo struct {
				LatestBlockHash   string `json:"latest_block_hash"`
				LatestAppHash     string `json:"latest_app_hash"`
				LatestBlockHeight string `json:"latest_block_height"`
				LatestBlockTime   string `json:"latest_block_time"`
				CatchingUp        bool   `json:"catching_up"`
			} `json:"sync_info"`
			ValidatorInfo struct {
				VotingPower string `json:"voting_power"`
			} `json:"validator_info"`
		} `json:"result"`
	}
)
