package main

import (
	"log"
	"time"

	"github.com/tpyle/reqx/lib/requests"
	"github.com/tpyle/reqx/lib/requests/context"
	reqxHttpContext "github.com/tpyle/reqx/lib/requests/context/http"
	reqxHttp "github.com/tpyle/reqx/lib/requests/http"
	reqxMultipart "github.com/tpyle/reqx/lib/requests/http/multipart"
)

func main() {
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

	reqx := requests.ReqX{
		Metadata: requests.Metadata{
			FriendlyName:     "Test Request",
			Order:            1,
			CustomProperties: nil,
		},
		Request: requests.Request{
			RequestType: requests.HTTP,
			Spec: reqxHttp.HTTPRequestSpec{
				URL: reqxHttp.HTTPRequestURL{
					Protocol: "http",
					Hostname: "localhost",
					Port:     8080,
					Path:     "/",
				},
				Method: "POST",
				Format: reqxHttp.MULTIPART,
				Data: reqxMultipart.MultipartFormData{
					{
						Name:     "test",
						Value:    "test",
						FileName: "",
					},
					{
						Name:     "test2",
						Value:    "test2",
						FileName: "",
					},
					{
						Name:     "test3",
						Value:    "",
						FileName: "test3.txt",
					},
				},
				// Format: reqxHttp.FORM,
				// Data: form.FormData{
				// 	"name":  []string{"test", "test2"},
				// 	"value": []string{"test"},
				// },
				// Format: reqxHttp.JSON,
				// Data: json.HTTPRequestJSONData{
				// 	"test": "test",
				// 	"test2": map[string]interface{}{
				// 		"test3": "test3",
				// 	},
				// },
			},
		},
		Assertions: []requests.Assertion{},
	}

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
