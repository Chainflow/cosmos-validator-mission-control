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

	// ValidatorResp defines validator result on a particular height
	ValidatorResp struct {
		OperatorAddress string `json:"operator_address"`
		Jailed          bool   `json:"jailed"`
	}

	// CommitInfo struct
	CommitInfo struct {
		ValidatorAddress string `json:"validator_address"`
		Signature        string `json:"signature"`
	}

	// LastCommitInfo stores block precommits
	LastCommitInfo struct {
		BlockID    interface{}  `json:"block_id"`
		Precommits []CommitInfo `json:"precommits"`
	}

	// BlockInfo stores latest block details
	BlockInfo struct {
		Header     interface{}    `json:"header"`
		Data       interface{}    `json:"data"`
		Evidence   interface{}    `json:"evidence"`
		LastCommit LastCommitInfo `json:"last_commit"`
	}

	// BlockResult stores block meta information
	BlockResult struct {
		BlockMeta interface{} `json:"block_meta"`
		Block     BlockInfo   `json:"block"`
	}

	// BlockResponse response of a block information
	BlockResponse struct {
		JSONRPC string      `json:"jsonrpc"`
		Result  BlockResult `json:"result"`
	}

	// NetworkLatestBlock stores latest block height info
	NetworkLatestBlock struct {
		Result struct {
			SyncInfo struct {
				LatestBlockHeight string `json:"latest_block_height"`
			} `json:"sync_info"`
		} `json:"result"`
	}
)
