package requests

import (
	"encoding/json"

	"github.com/mitchellh/mapstructure"

	"github.com/tpyle/reqx/lib/requests/context"
	"github.com/tpyle/reqx/lib/requests/http"
)

type RequestType string

const (
	HTTP RequestType = "http"
)

type RequestSpec interface {
	Send(*context.RequestContext) error
}

type Request struct {
	RequestType RequestType `json:"type"`
	Spec        RequestSpec `json:"spec"`
}

func (r *Request) UnmarshalJSON(b []byte) error {
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		panic(err)
	}

	r.RequestType = RequestType(m["type"].(string))

	var req RequestSpec = nil
	switch r.RequestType {
	case HTTP:
		req = http.HTTPRequestSpec{}
	}
	r.Spec = req

	if err := mapstructure.Decode(m["spec"], &r.Spec); err != nil {
		panic(err)
	}
	return nil
}

type Assertion struct {
	Operator string `json:"operator"`
	Field    string `json:"field"`
	Operand  string `json:"operand"`
	JsonPath string `json:"jsonPath"`
}

type Metadata struct {
	FriendlyName     string            `json:"friendlyName"`
	Order            int32             `json:"order"`
	CustomProperties map[string]string `json:"customProperties"`
}

type ReqX struct {
	Metadata   Metadata    `json:"metadata"`
	Request    Request     `json:"request"`
	Assertions []Assertion `json:"assertions"`
}
