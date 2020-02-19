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
		ProtocolVersion interface{}   `json:"protocol_version"`
		Id              string        `json:"id"`
		ListenAddr      string        `json:"listen_addr"`
		Network         string        `json:"network"`
		Version         string        `json:"version"`
		Channels        string        `json:"channels"`
		Moniker         string        `json:"moniker"`
		Other           interface{}   `json:"other"`
		SyncInfo        SyncInfo      `json:"sync_info"`
		ValidatorInfo   ValidatorInfo `json:"validator_info"`
	}

	GaiaCliStatus struct {
		NodeInfo GaiaCliStatusNodeInfo `json:"node_info"`
	}

	ValidatorUptime struct {
		Address string `json:"address"`
		Misses  string `json:"misses"`
		Period  string `json:"period"`
	}

	ValidatorDescription struct {
		Moniker  string `json:"moniker"`
		Identity string `json:"identity"`
		Website  string `json:"website"`
		Details  string `json:"details"`
	}

	ValidatorCommission struct {
		Rate          string `json:"rate"`
		MaxRate       string `json:"maxRate"`
		MaxChangeRate string `json:"maxChangeRate"`
		UpdateTime    string `json:"updateTime"`
	}

	ValidatorDetails struct {
		OperatorAddress         string               `json:"operatorAddress"`
		ConsensusPubKey         string               `json:"consensusPubkey"`
		Jailed                  bool                 `json:"jailed"`
		Tombstoned              bool                 `json:"tombstoned"`
		Status                  string               `json:"status"`
		Tokens                  string               `json:"tokens"`
		TokensSelfBonded        string               `json:"tokensSelfBonded"`
		DelegatorShares         string               `json:"delegatorShares"`
		Description             ValidatorDescription `json:"description"`
		UnbondingHeight         string               `json:"unbondingHeight"`
		Commission              ValidatorCommission  `json:"commission"`
		UnbondingCommissiontime string               `json:"unbondingCompletionTime"`
	}

	Validator struct {
		Address    string           `json:"address"`
		Weight     string           `json:"weight"`
		WeightRank int              `json:"weight_rank"`
		Uptime     ValidatorUptime  `json:"uptime"`
		Details    ValidatorDetails `json:"details"`
	}

	ValidatorResp struct {
		Validator Validator `json:"validator"`
	}

	AccountBalance struct {
		Denom  string `json:"denom"`
		Amount string `json:"amount"`
	}

	Account struct {
		Address       string           `json:"address"`
		Balance       []AccountBalance `json:"balance"`
		PubKey        string           `json:"pubKey"`
		AccountNumber string           `json:"accountNumber"`
		Sequence      string           `json:"sequence"`
		Vested        bool             `json:"vested"`
		VestingInfo   interface{}      `json:"vestingInfo"`
	}

	AccountResp struct {
		Account Account `json:"account"`
	}
)
