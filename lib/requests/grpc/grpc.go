package grpc

import (
	ctx "context"
	"fmt"
	"time"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/jhump/protoreflect/dynamic"
	"github.com/jhump/protoreflect/dynamic/grpcdynamic"
	"github.com/sirupsen/logrus"
	"github.com/tpyle/reqx/lib/requests/context"
	"google.golang.org/grpc"
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
	Server           GRPCServer             `json:"server"`
	ProtoInformation GRPCProtoInformation   `json:"proto"`
	Service          string                 `json:"service"`
	Method           string                 `json:"method"`
	Data             map[string]interface{} `json:"data"`
	Options          GRPCOptions            `json:"options"`
}

func loadProto(fileName string, includePaths []string) (*desc.FileDescriptor, error) {
	parser := protoparse.Parser{}
	parser.ImportPaths = includePaths

	descs, err := parser.ParseFiles(fileName)
	if err != nil {
		return nil, fmt.Errorf("Error parsing proto file: %w", err)
	}
	return descs[0], nil
}

func (s GRPCRequestSpec) Send(c *context.RequestContext) error {
	descriptor, err := loadProto(s.ProtoInformation.ProtoFile, s.ProtoInformation.IncludedDirectories)
	if err != nil {
		return err
	}

	service := descriptor.FindService(s.Service)
	if service == nil {
		return fmt.Errorf("Service %s not found in proto file", s.Service)
	}

	method := service.FindMethodByName(s.Method)
	if method == nil {
		return fmt.Errorf("Method %s not found in service %s", s.Method, s.Service)
	}

	inputType := method.GetInputType()

	inputMessage := dynamic.NewMessage(inputType)

	for k, v := range s.Data {
		if err := inputMessage.TrySetFieldByName(k, v); err != nil {
			return fmt.Errorf("Error setting field %s: %w", k, err)
		}
	}

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
