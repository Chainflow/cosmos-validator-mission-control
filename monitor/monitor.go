package monitor

import (
	"bytes"
	"chainflow-vitwit/config"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

func InitNodeMonitors(cfg *config.Config) *Monitors {
	return &Monitors{List: []NodeMonitor{
		{
			ExecutionType: "http",
			Name:          "Is Gaiad Running",
			HTTPOptions: HTTPOptions{
				Endpoint: cfg.NodeURL,
				Method:   http.MethodGet,
			},
			Func: CheckGaiad,
		},
		{
			ExecutionType: "http",
			Name:          "Number of peers",
			HTTPOptions: HTTPOptions{
				Endpoint: cfg.NodeURL + "net_info?",
				Method:   http.MethodGet,
			},
			Func: NumPeers,
		},
		{
			ExecutionType: "http",
			Name:          "Peer Addresses",
			HTTPOptions: HTTPOptions{
				Endpoint: cfg.NodeURL + "net_info?",
				Method:   http.MethodGet,
			},
			Func: PeerAddresses,
		},
		{
			ExecutionType: "cmd",
			Name:          "Is Validator Active",
			Func:          IsValidatorActive,
		},
		{
			ExecutionType: "cmd",
			Name:          "Current Block Height",
			Func:          CurrentBlockHeight,
		},
		{
			ExecutionType: "cmd",
			Name:          "Caught Up?",
			Func:          CaughtUp,
		},
		{
			ExecutionType: "http",
			Name:          "Operator Address",
			HTTPOptions: HTTPOptions{
				Endpoint: cfg.LCDEndpoint + "v1/validator/" + cfg.OperatorAddress,
				Method:   http.MethodGet,
			},
			Func: OperatorAddress,
		},
		{
			ExecutionType: "http",
			Name:          "Address",
			HTTPOptions: HTTPOptions{
				Endpoint: cfg.LCDEndpoint + "v1/validator/" + cfg.OperatorAddress,
				Method:   http.MethodGet,
			},
			Func: Address,
		},
		{
			ExecutionType: "http",
			Name:          "Fee",
			HTTPOptions: HTTPOptions{
				Endpoint: cfg.LCDEndpoint + "v1/validator/" + cfg.OperatorAddress,
				Method:   http.MethodGet,
			},
			Func: Fee,
		},
		{
			ExecutionType: "http",
			Name:          "Max Rate",
			HTTPOptions: HTTPOptions{
				Endpoint: cfg.LCDEndpoint + "v1/validator/" + "v1/validator/" + cfg.OperatorAddress,
				Method:   http.MethodGet,
			},
			Func: MaxRate,
		},
		{
			ExecutionType: "http",
			Name:          "Max Change Rate",
			HTTPOptions: HTTPOptions{
				Endpoint: cfg.LCDEndpoint + "v1/validator/" + cfg.OperatorAddress,
				Method:   http.MethodGet,
			},
			Func: MaxChangeRate,
		},
		{
			ExecutionType: "http",
			Name:          "Address Balance",
			HTTPOptions: HTTPOptions{
				Endpoint: cfg.LCDEndpoint + "v1/account/" + cfg.AccountAddress,
				Method:   http.MethodGet,
			},
			Func: AddressBalance,
		},
		{
			ExecutionType: "cmd",
			Name:          "Voting Power",
			Func:          VotingPower,
		},
		{
			ExecutionType: "http",
			Name:          "Validator Details",
			HTTPOptions: HTTPOptions{
				Endpoint: cfg.LCDEndpoint + "v1/validator/" + "v1/validator/" + cfg.OperatorAddress,
				Method:   http.MethodGet,
			},
			Func: ValidatorDesc,
		},
		{
			ExecutionType: "cmd",
			Name:          "Gaiad Version",
			Func:          GaiadVersion,
		},
	}}
}

func addQueryParameters(req *http.Request, queryParams QueryParams) {
	params := url.Values{}
	for key, value := range queryParams {
		params.Add(key, value)
	}
	req.URL.RawQuery = params.Encode()
}

func newHttpRequest(m HTTPOptions) (*http.Request, error) {
	// make new request
	req, err := http.NewRequest(m.Method, m.Endpoint, bytes.NewBuffer(m.Body))
	if err != nil {
		return nil, err
	}

	// Add any query parameters to the URL.
	if len(m.QueryParams) != 0 {
		addQueryParameters(req, m.QueryParams)
	}

	return req, nil
}

func makeResponse(res *http.Response) (*PingResp, error) {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return &PingResp{}, err
	}

	response := &PingResp{
		StatusCode: res.StatusCode,
		Body:       body,
	}
	_ = res.Body.Close()
	return response, nil
}

func RunMonitor(m HTTPOptions) (*PingResp, error) {
	req, err := newHttpRequest(m)
	if err != nil {
		return nil, err
	}

	client := http.Client{Timeout: time.Duration(5 * time.Second)}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	res, err := makeResponse(resp)
	if err != nil {
		return nil, err
	}

	return res, nil
}
