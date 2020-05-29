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

	// ValidatorDescription struct
	ValidatorMetaInfo struct {
		Moniker  string `json:"moniker"`
		Identity string `json:"identity"`
		Website  string `json:"website"`
		Details  string `json:"details"`
	}

	// ValidatorCommissionRates struct
	ValidatorCommissionRates struct {
		Rate          string `json:"rate"`
		MaxRate       string `json:"max_rate"`
		MaxChangeRate string `json:"max_change_rate"`
	}

	// ValidatorCommission struct
	ValidatorCommission struct {
		CommissionRates ValidatorCommissionRates `json:"commission_rates"`
		UpdateTime      string                   `json:"update_time"`
	}

	// ValidatorResult struct
	ValidatorDetails struct {
		OperatorAddress   string              `json:"operator_address"`
		ConsensusPubKey   string              `json:"consensus_pubkey"`
		Jailed            bool                `json:"jailed"`
		Status            int                 `json:"status"`
		Tokens            string              `json:"tokens"`
		DelegatorShares   string              `json:"delegator_shares"`
		Description       ValidatorMetaInfo   `json:"description"`
		UnbondingHeight   string              `json:"unbonding_height"`
		UnbondingTime     string              `json:"unbonding_time"`
		Commission        ValidatorCommission `json:"commission"`
		MinSelfDelegation string              `json:"min_self_delegation"`
	}

	// ValidatorResp defines validator result on a particular height
	ValidatorResp struct {
		Height string           `json:"height"`
		Result ValidatorDetails `json:"result"`
	}

	// CommitInfo struct
	CommitInfo struct {
		Type             int64       `json:"type"`
		Height           string      `json:"height"`
		Round            string      `json:"round"`
		BlockID          interface{} `json:"block_id"`
		Timestamp        string      `json:"timestamp"`
		ValidatorAddress string      `json:"validator_address"`
		ValidatorIndex   string      `json:"validator_index"`
		Signature        string      `json:"signature"`
	}

	// LastCommitInfo stores block precommits
	LastCommitInfo struct {
		BlockID    interface{}  `json:"block_id"`
		Precommits []CommitInfo `json:"precommits"`
	}

	// BlockInfo stores latest block details
	BlockInfo struct {
		Header struct {
			Height string `json:"height"`
			Time   string `json:"time`
		} `json:"header"`
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

	// AkashBlockResponse
	AkashBlockResponse struct {
		Jsonrpc string `json:"jsonrpc"`
		Result  struct {
			BlockID interface{} `json:"block_id"`
			Block   struct {
				Header struct {
					ChainID string `json:"chain_id"`
					Height  string `json:"height"`
					Time    string `json:"time"`
				} `json:"header"`
				Data       interface{} `json:"data"`
				Evidence   interface{} `json:"evidence"`
				LastCommit struct {
					Height     string `json:"height"`
					Signatures []struct {
						BlockIDFlag      int    `json:"block_id_flag"`
						ValidatorAddress string `json:"validator_address"`
						Timestamp        string `json:"timestamp"`
						Signature        string `json:"signature"`
					} `json:"signatures"`
				} `json:"last_commit"`
			} `json:"block"`
		} `json:"result"`
	}
)
