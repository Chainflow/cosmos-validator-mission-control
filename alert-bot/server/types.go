package server

import (
	"cosmos-validator-mission-control/alert-bot/config"

	client "github.com/influxdata/influxdb1-client/v2"
)

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
		// Type             int64       `json:"type"`
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
		JSONRPC string `json:"jsonrpc"`
		// ID      string      `json:"id"`
		Result BlockResult `json:"result"`
	}

	// NetworkLatestBlock stores latest block height info
	NetworkLatestBlock struct {
		Result struct {
			SyncInfo struct {
				LatestBlockHeight string `json:"latest_block_height"`
			} `json:"sync_info"`
		} `json:"result"`
	}

	// Target is a structure which holds all the parameters of a target
	//this could be used to write endpoints for each functionality
	Target struct {
		ExecutionType string
		HTTPOptions   HTTPOptions
		Name          string
		Func          func(m HTTPOptions, cfg *config.Config, c client.Client)
		ScraperRate   string
	}

	// Targets list of all the targets
	Targets struct {
		List []Target
	}

	ValidatorRpcStatus struct {
		Jsonrpc string `json:"jsonrpc"`
		Result  struct {
			NodeInfo interface{} `json:"node_info"`
			SyncInfo struct {
				LatestBlockHash   string `json:"latest_block_hash"`
				LatestAppHash     string `json:"latest_app_hash"`
				LatestBlockHeight string `json:"latest_block_height"`
				LatestBlockTime   string `json:"latest_block_time"`
				CatchingUp        bool   `json:"catching_up"`
			} `json:"sync_info"`
			ValidatorInfo struct {
				Address string `json:"address"`
				PubKey  struct {
					Type  string `json:"type"`
					Value string `json:"value"`
				} `json:"pub_key"`
				VotingPower string `json:"voting_power"`
			} `json:"validator_info"`
		} `json:"result"`
	}

	// ValidatorsHeight struct which represents the details of validator
	ValidatorsHeight struct {
		Jsonrpc string `json:"jsonrpc"`
		Result  struct {
			BlockHeight string `json:"block_height"`
			Validators  []struct {
				Address string `json:"address"`
				PubKey  struct {
					Type  string `json:"type"`
					Value string `json:"value"`
				} `json:"pub_key"`
				VotingPower      string `json:"voting_power"`
				ProposerPriority string `json:"proposer_priority"`
			} `json:"validators"`
		} `json:"result"`
	}

	// Peer is a structure which holds the info about a peer address
	Peer struct {
		RemoteIP         string      `json:"remote_ip"`
		ConnectionStatus interface{} `json:"connection_status"`
		IsOutbound       bool        `json:"is_outbound"`
		NodeInfo         struct {
			Moniker string `json:"moniker"`
			Network string `json:"network"`
		} `json:"node_info"`
	}

	// NetInfoResult struct
	NetInfoResult struct {
		Listening bool          `json:"listening"`
		Listeners []interface{} `json:"listeners"`
		NumPeers  string        `json:"n_peers"`
		Peers     []Peer        `json:"peers"`
	}

	// NetInfo is a structre which holds the details of address
	NetInfo struct {
		JSONRpc string        `json:"jsonrpc"`
		Result  NetInfoResult `json:"result"`
	}
)
