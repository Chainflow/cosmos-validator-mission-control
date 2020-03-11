package targets

import (
	"chainflow-vitwit/config"

	client "github.com/influxdata/influxdb1-client/v2"
)

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

	// struct for the target
	Target struct {
		ExecutionType string
		HTTPOptions   HTTPOptions
		Name          string
		Func          func(m HTTPOptions, cfg *config.Config, c client.Client)
	}

	// Array of targets
	Targets struct {
		List []Target
	}

	// Ping response
	PingResp struct {
		StatusCode int
		Body       []byte
	}

	// Peer response
	Peer struct {
		RemoteIP         string      `json:"remote_ip"`
		ConnectionStatus interface{} `json:"connection_status"`
		IsOutbound       bool        `json:"is_outbound"`
		NodeInfo         struct {
			Moniker string `json:"moniker"`
			Network string `json:"network"`
		} `json:"node_info"`
	}

	// Net info result
	NetInfoResult struct {
		Listening bool          `json:"listening"`
		Listeners []interface{} `json:"listeners"`
		NumPeers  string        `json:"n_peers"`
		Peers     []Peer        `json:"peers"`
	}

	// Net info response
	NetInfo struct {
		JSONRpc string        `json:"jsonrpc"`
		Id      string        `json:"id"`
		Result  NetInfoResult `json:"result"`
	}

	// Sync info response
	SyncInfo struct {
		LatestBlockHash   string `json:"latest_block_hash"`
		LatestAppHash     string `json:"latest_app_hash"`
		LatestBlockHeight string `json:"latest_block_height"`
		LatestBlockTime   string `json:"latest_block_time"`
		CatchingUp        bool   `json:"catching_up"`
	}

	// Validator info response
	ValidatorInfo struct {
		Address     string      `json:"address"`
		PubKey      interface{} `json:"pub_key"`
		VotingPower string      `json:"voting_power"`
	}

	// Node info
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

	// gaia cli status
	GaiaCliStatus struct {
		NodeInfo      GaiaCliStatusNodeInfo `json:"node_info"`
		SyncInfo      SyncInfo              `json:"sync_info"`
		ValidatorInfo ValidatorInfo         `json:"validator_info"`
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

	// Account Balance
	AccountBalance struct {
		Denom  string `json:"denom"`
		Amount string `json:"amount"`
	}

	// Account response
	AccountResp struct {
		Height string           `json:"height"`
		Result []AccountBalance `json:"result"`
	}

	// Precommits of current block
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

	// Last commit of current block
	CurrentBlockLastCommit struct {
		BlockId    interface{}             `json:"block_id"`
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
		Id      string                       `json:"id"`
		Result  CurrentBlockWithHeightResult `json:"result"`
	}

	// Proposal result
	ProposalResultContent struct {
		Type  string `json:"type"`
		Value struct {
			Title       string `json:"title"`
			Description string `json:"description"`
		} `json:"value"`
	}

	// Proposal result
	ProposalResult struct {
		Content          ProposalResultContent `json:"content"`
		Id               string                `json:"id"`
		ProposalStatus   string                `json:"proposal_status"`
		FinalTallyResult interface{}           `json:"final_tally_result"`
		SubmitTime       string                `json:"submit_time"`
		DepositEndTime   string                `json:"deposit_end_time"`
		TotalDeposit     []interface{}         `json:"total_deposit"`
		VotingStartTime  string                `json:"voting_start_time"`
		VotingEndTime    string                `json:"voting_end_time"`
	}

	// Proposal response
	Proposals struct {
		Height string           `json:"height"`
		Result []ProposalResult `json:"result"`
	}

	// Self delegation balance
	SelfDelegationBalance struct {
		Balance string `json:"balance"`
	}

	// Self delegation response
	SelfDelegation struct {
		Height string                `json:"height"`
		Result SelfDelegationBalance `json:"result"`
	}

	// Current rewards amount response
	CurrentRewardsAmount struct {
		Height string           `json:"height"`
		Result []AccountBalance `json:"result"`
	}

	// Last proposed block response
	LastProposedBlockAndTime struct {
		BlockMeta struct {
			BlockID interface{} `json:"block_id"`
			Header  struct {
				Version struct {
					Block string `json:"block"`
					App   string `json:"app"`
				} `json:"version"`
				ChainID            string      `json:"chain_id"`
				Height             string      `json:"height"`
				Time               string      `json:"time"`
				NumTxs             string      `json:"num_txs"`
				TotalTxs           string      `json:"total_txs"`
				LastBlockID        interface{} `json:"last_block_id"`
				LastCommitHash     string      `json:"last_commit_hash"`
				DataHash           string      `json:"data_hash"`
				ValidatorsHash     string      `json:"validators_hash"`
				NextValidatorsHash string      `json:"next_validators_hash"`
				ConsensusHash      string      `json:"consensus_hash"`
				AppHash            string      `json:"app_hash"`
				LastResultsHash    string      `json:"last_results_hash"`
				EvidenceHash       string      `json:"evidence_hash"`
				ProposerAddress    string      `json:"proposer_address"`
			} `json:"header"`
		} `json:"block_meta"`
		Block interface{} `json:"block"`
	}

	// Voters of a proposal
	ProposalVoters struct {
		Height string `json:"height"`
		Result []struct {
			ProposalID string `json:"proposal_id"`
			Voter      string `json:"voter"`
			Option     string `json:"option"`
		} `json:"result"`
	}

	// Latest block of a network
	NetworkLatestBlock struct {
		Result struct {
			SyncInfo struct {
				LatestBlockHeight string `json:"latest_block_height"`
			} `json:"sync_info"`
		} `json:"result"`
	}

	// Validator block response
	ValidatorsHeight struct {
		Jsonrpc string `json:"jsonrpc"`
		ID      string `json:"id"`
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

	// Response of
	// proposal depositors
	Depositors struct {
		Height string `json:"height"`
		Result []struct {
			ProposalID string `json:"proposal_id"`
			Depositor  string `json:"depositor"`
			Amount     []struct {
				Denom  string `json:"denom"`
				Amount string `json:"amount"`
			} `json:"amount"`
		} `json:"result"`
	}
)
