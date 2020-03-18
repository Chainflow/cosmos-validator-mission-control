package server

type (
	// map of strings for query params
	QueryParams map[string]string

	// struct for http options
	HTTPOptions struct {
		Endpoint    string
		QueryParams QueryParams
		Body        []byte
		Method      string
	}

	// Ping response
	PingResp struct {
		StatusCode int
		Body       []byte
	}

	// Validator description
	ValidatorDescription struct {
		Moniker  string `json:"moniker"`
		Identity string `json:"identity"`
		Website  string `json:"website"`
		Details  string `json:"details"`
	}

	// Validator commission rates
	ValidatorCommissionRates struct {
		Rate          string `json:"rate"`
		MaxRate       string `json:"max_rate"`
		MaxChangeRate string `json:"max_change_rate"`
	}

	// Validator commission
	ValidatorCommission struct {
		CommissionRates ValidatorCommissionRates `json:"commission_rates"`
		UpdateTime      string                   `json:"update_time"`
	}

	// Validator result
	ValidatorResult struct {
		OperatorAddress   string               `json:"operator_address"`
		ConsensusPubKey   string               `json:"consensus_pubkey"`
		Jailed            bool                 `json:"jailed"`
		Status            int                  `json:"status"`
		Tokens            string               `json:"tokens"`
		DelegatorShares   string               `json:"delegator_shares"`
		Description       ValidatorDescription `json:"description"`
		UnbondingHeight   string               `json:"unbonding_height"`
		UnbondingTime     string               `json:"unbonding_time"`
		Commission        ValidatorCommission  `json:"commission"`
		MinSelfDelegation string               `json:"min_self_delegation"`
	}

	// Validator response
	ValidatorResp struct {
		Height string          `json:"height"`
		Result ValidatorResult `json:"result"`
	}

	// Precommits of current block
	CurrentBlockPrecommit struct {
		Type             int64       `json:"type"`
		Height           string      `json:"height"`
		Round            string      `json:"round"`
		BlockID          interface{} `json:"block_id"`
		Timestamp        string      `json:"timestamp"`
		ValidatorAddress string      `json:"validator_address"`
		ValidatorIndex   string      `json:"validator_index"`
		Signature        string      `json:"signature"`
	}

	// Last commit of current block
	CurrentBlockLastCommit struct {
		BlockID    interface{}             `json:"block_id"`
		Precommits []CurrentBlockPrecommit `json:"precommits"`
	}

	// Current block
	CurrentBlock struct {
		Header struct {
			Height string `json:"height"`
			Time   string `json:"time`
		} `json:"header"`
		Data       interface{}            `json:"data"`
		Evidence   interface{}            `json:"evidence"`
		LastCommit CurrentBlockLastCommit `json:"last_commit"`
	}

	// Current block details
	CurrentBlockWithHeightResult struct {
		BlockMeta interface{}  `json:"block_meta"`
		Block     CurrentBlock `json:"block"`
	}

	// Current block height response
	CurrentBlockWithHeight struct {
		JSONRPC string                       `json:"jsonrpc"`
		ID      string                       `json:"id"`
		Result  CurrentBlockWithHeightResult `json:"result"`
	}

	// Latest block of a network
	NetworkLatestBlock struct {
		Result struct {
			SyncInfo struct {
				LatestBlockHeight string `json:"latest_block_height"`
			} `json:"sync_info"`
		} `json:"result"`
	}
)
