package grpc

import (
	ctx "context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jhump/protoreflect/dynamic"
	"github.com/jhump/protoreflect/dynamic/grpcdynamic"
	"github.com/sirupsen/logrus"
	"github.com/tpyle/reqx/lib/requests/context"
	"google.golang.org/grpc"

	grpcUtils "github.com/tpyle/reqx/lib/util/grpc"
)

type GRPCServer struct {
	Hostname string `json:"hostname"`
	Port     int32  `json:"port"`
}

type GRPCOptions struct {
	Timeout  time.Duration `json:"timeout,omitempty"`
	Insecure bool          `json:"insecure,omitempty"`
}

type GRPCProtoInformation struct {
	ProtoFile           string   `json:"file"`
	IncludedDirectories []string `json:"include"`
}

type GRPCRequestSpec struct {
	Server           GRPCServer           `json:"server"`
	ProtoInformation GRPCProtoInformation `json:"proto"`
	Service          string               `json:"service"`
	Method           string               `json:"method"`
	Data             json.RawMessage      `json:"data"`
	Options          GRPCOptions          `json:"options"`
}

func (s GRPCRequestSpec) Send(c *context.RequestContext) error {
	logrus.Infof("Sending GRPC request: %+v", s)
	descriptor, err := grpcUtils.LoadProto(s.ProtoInformation.ProtoFile, s.ProtoInformation.IncludedDirectories)
	if err != nil {
		return err
	}

	method, err := grpcUtils.ResolveMethod(descriptor, s.Service, s.Method)
	if err != nil {
		return err
	}

	inputType := method.GetInputType()
	inputMessage := dynamic.NewMessage(inputType)
	inputMessage.UnmarshalJSON(s.Data)

	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", s.Server.Hostname, s.Server.Port), grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("Error dialing GRPC server: %w", err)
	}

	client := grpcdynamic.NewStub(conn)

	response, err := client.InvokeRpc(ctx.Background(), method, inputMessage)
	if err != nil {
		return fmt.Errorf("Error invoking GRPC method: %w", err)
	}

	jsonFormat, err := response.(*dynamic.Message).MarshalJSON()
	if err != nil {
		return fmt.Errorf("Error marshalling response to JSON: %w", err)
	}

	logrus.Infof("Response: %s", jsonFormat)

	return nil
}
