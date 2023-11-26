package grpc

import (
	"fmt"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
)

func ResolveMethod(descriptor *desc.FileDescriptor, service string, method string) (*desc.MethodDescriptor, error) {
	if descriptor == nil {
		return nil, fmt.Errorf("No File Descriptor provided")
	}
	serviceDescriptor := descriptor.FindService(service)
	if serviceDescriptor == nil {
		for _, serviceItem := range descriptor.GetServices() {
			if serviceItem.GetName() == service {
				serviceDescriptor = serviceItem
				break
			}
		}
	}
	if serviceDescriptor == nil {
		return nil, fmt.Errorf("Service %s not found", service)
	}
	methodDescriptor := serviceDescriptor.FindMethodByName(method)
	if methodDescriptor == nil {
		return nil, fmt.Errorf("Method %s not found", method)
	}
	return methodDescriptor, nil
}

func LoadProto(fileName string, includePaths []string) (*desc.FileDescriptor, error) {
	parser := protoparse.Parser{}
	parser.ImportPaths = includePaths

	descs, err := parser.ParseFiles(fileName)
	if err != nil {
		return nil, fmt.Errorf("Error parsing proto file: %w", err)
	}
	return descs[0], nil
}
