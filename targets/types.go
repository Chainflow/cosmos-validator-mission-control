package targets

import (
	"chainflow-vitwit/config"
	client "github.com/influxdata/influxdb1-client/v2"
)

type (
	QueryParams map[string]string

	HTTPOptions struct {
		Endpoint    string
		QueryParams QueryParams
		Body        []byte
		Method      string
	}

	Target struct {
		ExecutionType string
		HTTPOptions   HTTPOptions
		Name          string
		Func          func(m HTTPOptions, cfg *config.Config, c client.Client)
	}

	Targets struct {
		List []Target
	}

	PingResp struct {
		StatusCode int
		Body       []byte
	}

	Peer struct {
		RemoteIP         string      `json:"remote_ip"`
		ConnectionStatus interface{} `json:"connection_status"`
		IsOutbound       bool        `json:"is_outbound"`
		NodeInfo         interface{} `json:"node_info"`
	}

	NetInfoResult struct {
		Listening bool          `json:"listening"`
		Listeners []interface{} `json:"listeners"`
		NumPeers  string        `json:"n_peers"`
		Peers     []Peer        `json:"peers"`
	}

	NetInfo struct {
		JSONRpc string        `json:"jsonrpc"`
		Id      string        `json:"id"`
		Result  NetInfoResult `json:"result"`
	}

	SyncInfo struct {
		LatestBlockHash   string `json:"latest_block_hash"`
		LatestAppHash     string `json:"latest_app_hash"`
		LatestBlockHeight string `json:"latest_block_height"`
		LatestBlockTime   string `json:"latest_block_time"`
		CatchingUp        bool   `json:"catching_up"`
	}

	ValidatorInfo struct {
		Address     string      `json:"address"`
		PubKey      interface{} `json:"pub_key"`
		VotingPower string      `json:"voting_power"`
	}

	GaiaCliStatusNodeInfo struct {
		ProtocolVersion interface{} `json:"protocol_version"`
		Id              string      `json:"id"`
		ListenAddr      string      `json:"listen_addr"`
		Network         string      `json:"network"`
		Version         string      `json:"version"`
		Channels        string      `json:"channels"`
		Moniker         string      `json:"moniker"`
		Other           interface{} `json:"other"`
	}

	GaiaCliStatus struct {
		NodeInfo      GaiaCliStatusNodeInfo `json:"node_info"`
		SyncInfo      SyncInfo              `json:"sync_info"`
		ValidatorInfo ValidatorInfo         `json:"validator_info"`
	}

	ValidatorDescription struct {
		Moniker  string `json:"moniker"`
		Identity string `json:"identity"`
		Website  string `json:"website"`
		Details  string `json:"details"`
	}

	ValidatorCommissionRates struct {
		Rate          string `json:"rate"`
		MaxRate       string `json:"max_rate"`
		MaxChangeRate string `json:"max_change_rate"`
	}

	ValidatorCommission struct {
		CommissionRates ValidatorCommissionRates `json:"commission_rates"`
		UpdateTime      string                   `json:"update_time"`
	}

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

	ValidatorResp struct {
		Height string          `json:"height"`
		Result ValidatorResult `json:"result"`
	}

	AccountBalance struct {
		Denom  string `json:"denom"`
		Amount string `json:"amount"`
	}

	AccountResp struct {
		Height string           `json:"height"`
		Result []AccountBalance `json:"result"`
	}

	CurrentBlockPrecommit struct {
		Type             int64       `json:"type"`
		Height           string      `json:"height"`
		Round            string      `json:"round"`
		BlockId          interface{} `json:"block_id"`
		Timestamp        string      `json:"timestamp"`
		ValidatorAddress string      `json:"validator_address"`
		ValidatorIndex   string      `json:"validator_index"`
		Signature        string      `json:"signature"`
	}

	CurrentBlockLastCommit struct {
		BlockId    interface{}             `json:"block_id"`
		Precommits []CurrentBlockPrecommit `json:"precommits"`
	}

	CurrentBlock struct {
		Header     interface{}            `json:"header"`
		Data       interface{}            `json:"data"`
		Evidence   interface{}            `json:"evidence"`
		LastCommit CurrentBlockLastCommit `json:"last_commit"`
	}

	CurrentBlockWithHeightResult struct {
		BlockMeta interface{}  `json:"block_meta"`
		Block     CurrentBlock `json:"block"`
	}

	CurrentBlockWithHeight struct {
		JSONRPC string                       `json:"jsonrpc"`
		Id      string                       `json:"id"`
		Result  CurrentBlockWithHeightResult `json:"result"`
	}
)
