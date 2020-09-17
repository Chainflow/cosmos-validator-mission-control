package server

import "time"

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

	// BlockResponse response of a block information
	BlockResponse struct {
		Jsonrpc string `json:"jsonrpc"`
		Result  struct {
			BlockID interface{} `json:"block_id"`
			Block   struct {
				Header interface{} `json:"header"`
				Data   struct {
					Txs interface{} `json:"txs"`
				} `json:"data"`
				Evidence struct {
					Evidence interface{} `json:"evidence"`
				} `json:"evidence"`
				LastCommit struct {
					Height     string      `json:"height"`
					Round      string      `json:"round"`
					BlockID    interface{} `json:"block_id"`
					Signatures []struct {
						BlockIDFlag      int       `json:"block_id_flag"`
						ValidatorAddress string    `json:"validator_address"`
						Timestamp        time.Time `json:"timestamp"`
						Signature        string    `json:"signature"`
					} `json:"signatures"`
				} `json:"last_commit"`
			} `json:"block"`
		} `json:"result"`
	}

	// ValidatorResp defines validator result on a particular height
	ValidatorResp struct {
		Height string `json:"height"`
		Result struct {
			OperatorAddress string `json:"operator_address"`
			ConsensusPubkey string `json:"consensus_pubkey"`
			Jailed          bool   `json:"jailed"`
			Status          int    `json:"status"`
			Tokens          string `json:"tokens"`
			DelegatorShares string `json:"delegator_shares"`
			Description     struct {
				Moniker         string `json:"moniker"`
				Identity        string `json:"identity"`
				Website         string `json:"website"`
				SecurityContact string `json:"security_contact"`
				Details         string `json:"details"`
			} `json:"description"`
			UnbondingHeight string    `json:"unbonding_height"`
			UnbondingTime   time.Time `json:"unbonding_time"`
			Commission      struct {
				CommissionRates struct {
					Rate          string `json:"rate"`
					MaxRate       string `json:"max_rate"`
					MaxChangeRate string `json:"max_change_rate"`
				} `json:"commission_rates"`
				UpdateTime time.Time `json:"update_time"`
			} `json:"commission"`
			MinSelfDelegation string `json:"min_self_delegation"`
		} `json:"result"`
	}
)
