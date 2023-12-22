package requests

import (
	"encoding/json"
	"fmt"
	"os"

	"sigs.k8s.io/yaml"

	"github.com/tpyle/reqx/lib/requests/context"
	"github.com/tpyle/reqx/lib/requests/grpc"
	"github.com/tpyle/reqx/lib/requests/http"
)

type RequestType string

const (
	HTTP RequestType = "HTTP"
	GRPC RequestType = "GRPC"
)

type RequestSpec interface {
	Send(*context.RequestContext) error
}

type Request struct {
	RequestType RequestType `json:"type"`
	Spec        RequestSpec `json:"spec"`
}

func (r *Request) UnmarshalJSON(b []byte) error {
	var temp struct {
		RequestType RequestType     `json:"type"`
		Spec        json.RawMessage `json:"spec"`
	}
	if err := json.Unmarshal(b, &temp); err != nil {
		return fmt.Errorf("Unable to unmarshal request: %w", err)
	}

	r.RequestType = temp.RequestType

	switch r.RequestType {
	case HTTP:
		var req http.HTTPRequestSpec
		if err := json.Unmarshal(temp.Spec, &req); err != nil {
			return fmt.Errorf("Unable to unmarshal HTTP request: %w", err)
		}
		r.Spec = req
	case GRPC:
		var req grpc.GRPCRequestSpec
		if err := json.Unmarshal(temp.Spec, &req); err != nil {
			return fmt.Errorf("Unable to unmarshal GRPC request: %w", err)
		}
		r.Spec = req
	default:
		return fmt.Errorf("Unknown request type: %s", r.RequestType)
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

func LoadFromFile(filename string) (*ReqX, error) {
	body, err := os.ReadFile(filename)

	var reqx ReqX
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(body, &reqx); err != nil {
		return nil, err
	}
	return &reqx, nil
}
