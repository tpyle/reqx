package main

import (
	"log"
	"time"

	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/jhump/protoreflect/dynamic"
	"github.com/tpyle/reqx/lib/requests"
	"github.com/tpyle/reqx/lib/requests/context"
	reqxHttpContext "github.com/tpyle/reqx/lib/requests/context/http"
	"github.com/tpyle/reqx/lib/requests/grpc"
)

func test() {
	parser := protoparse.Parser{
		ImportPaths: []string{"./trifecta-schemas/site"},
	}
	fds, err := parser.ParseFiles("./site_service.proto")
	if err != nil {
		log.Fatal(err)
	}
	for _, fd := range fds {
		for _, service := range fd.GetServices() {
			log.Println(service.GetName())
			for _, method := range service.GetMethods() {
				log.Println(method.GetFullyQualifiedName())
				log.Println(method.GetInputType().GetFullyQualifiedName())
				log.Println(method.GetOutputType().GetFullyQualifiedName())
				if method.GetName() == "getById" {
					log.Println("Found it!")
					message := method.GetInputType()
					log.Printf("Message: %v", message)
					messageInstance := dynamic.NewMessage(message)
					log.Printf("Message Instance: %v", messageInstance)
					messageInstance.UnmarshalJSON([]byte(`{"id": 1}`))
					log.Printf("Message Instance: %v", messageInstance)
					jsonContent, err := messageInstance.MarshalJSON()
					if err != nil {
						log.Fatal(err)
					}
					log.Printf("JSON Content: %s", jsonContent)
				}
			}
		}
	}
}

func main() {
	// test()
	// return
	// var spec requests.RequestSpec = reqxHttp.HTTPRequestSpec{
	// 	URL: reqxHttp.HTTPRequestURL{
	// 		Protocol: "http",
	// 		Hostname: "localhost",
	// 		Port:     8080,
	// 		Path:     "/",
	// 	},
	// 	Method: "POST",
	// 	Format: reqxHttp.JSON,
	// 	Data: json.HTTPRequestJSONData{
	// 		"test": "test",
	// 		"test2": map[string]interface{}{
	// 			"test3": "test3",
	// 		},
	// 	},
	// }

	// var id uint32 = 1
	reqx := requests.ReqX{
		Metadata: requests.Metadata{
			FriendlyName:     "Test Request",
			Order:            1,
			CustomProperties: nil,
		},
		Request: requests.Request{
			RequestType: requests.GRPC,
			Spec: grpc.GRPCRequestSpec{
				Server: grpc.GRPCServer{
					Hostname: "localhost",
					Port:     10000,
				},
				ProtoInformation: grpc.GRPCProtoInformation{
					ProtoFile:           "site_service.proto",
					IncludedDirectories: []string{"./trifecta-schemas/site"},
				},
				Service: "zebra.site.SiteService",
				Method:  "getById",
				Data:    []byte(`{"id": 1}`),
			},
		},
		Assertions: []requests.Assertion{},
	}

	// reqx := requests.ReqX{
	// 	Metadata: requests.Metadata{
	// 		FriendlyName:     "Test Request",
	// 		Order:            1,
	// 		CustomProperties: nil,
	// 	},
	// 	Request: requests.Request{
	// 		RequestType: requests.HTTP,
	// 		Spec: reqxHttp.HTTPRequestSpec{
	// 			URL: reqxHttp.HTTPRequestURL{
	// 				Protocol: "https",
	// 				Hostname: "localhost",
	// 				Port:     8443,
	// 				Path:     "/",
	// 			},
	// 			Method: "POST",
	// 			Format: reqxHttp.MULTIPART,
	// 			Data: reqxMultipart.MultipartFormData{
	// 				{
	// 					Name:     "test",
	// 					Value:    "test",
	// 					FileName: "",
	// 				},
	// 				{
	// 					Name:     "test2",
	// 					Value:    "test2",
	// 					FileName: "",
	// 				},
	// 				{
	// 					Name:     "test3",
	// 					Value:    "",
	// 					FileName: "test3.txt",
	// 				},
	// 			},
	// 			Headers: map[string]string{
	// 				"test": "test",
	// 			},
	// 			Options: reqxHttp.HTTPRequestOptions{
	// 				SkipTLSVerify: false,
	// 			},
	// 			// Format: reqxHttp.FORM,
	// 			// Data: form.FormData{
	// 			// 	"name":  []string{"test", "test2"},
	// 			// 	"value": []string{"test"},
	// 			// },
	// 			// Format: reqxHttp.JSON,
	// 			// Data: json.HTTPRequestJSONData{
	// 			// 	"test": "test",
	// 			// 	"test2": map[string]interface{}{
	// 			// 		"test3": "test3",
	// 			// 	},
	// 			// },
	// 		},
	// 	},
	// 	Assertions: []requests.Assertion{},
	// }

	requestContext := context.RequestContext{
		HTTPContext: reqxHttpContext.HTTPRequestContext{
			Timeout: time.Second * 10,
		},
		FileLocation: "./",
	}
	err := reqx.Request.Spec.Send(&requestContext)
	if err != nil {
		log.Fatal(err)
	}

	// SendJSONRequest()
	// SendGraphQLRequest()
	// SendFormRequest()
	// SendMultiPartRequest()

	// protoFile := `
	//     syntax = "proto3";
	//     package test;
	//     message TestMessage {
	//         string name = 1;
	//     }
	// `

	// // This is your JSON data
	// jsonData := `{"name":"test"}`

	// // Parse the proto file
	// protodesc.NewFile(protoreflect.FileDescriptor{}, nil)
	// fd, err := protodesc.FromFileDescriptorProto(protoFile)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // Get the message descriptor for the message type
	// md := fd.Messages().ByName(protoreflect.Name("TestMessage"))

	// // Create a new dynamic message
	// dynMsg := dynamicpb.NewMessage(md)

	// // Unmarshal the JSON data into the dynamic message
	// err = protojson.Unmarshal([]byte(jsonData), dynMsg)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // Now dynMsg is a protobuf message populated with the data from jsonData
	// log.Println(dynMsg)
}
