package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/tpyle/reqx/lib/requests"
	"github.com/tpyle/reqx/lib/requests/context"
	reqxHttpContext "github.com/tpyle/reqx/lib/requests/context/http"
	reqxHttp "github.com/tpyle/reqx/lib/requests/http"
	"github.com/tpyle/reqx/lib/requests/http/form"
)

func SendJSONRequest() {
	netClient := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest("POST", "http://localhost:8080", bytes.NewBuffer([]byte(`{"name":"test"}`)))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := netClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Response Body:", string(body))
}

func SendGraphQLRequest() {
	netClient := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest("POST", "http://localhost:8080", bytes.NewBuffer([]byte(`{"query":"query { hello }"}`)))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/graphql")
	resp, err := netClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Response Body:", string(body))
}

func SendFormRequest() {
	netClient := &http.Client{
		Timeout: time.Second * 10,
	}
	formData := url.Values{
		"name": {"test"},
	}
	req, err := http.NewRequest("POST", "http://localhost:8080", strings.NewReader(formData.Encode()))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := netClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Response Body:", string(body))
}

func GetMimeType(file *os.File) (string, error) {
	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)
	_, err := file.Read(buffer)
	if err != nil {
		return "", err
	}

	// Reset the read pointer if necessary.
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return "", err
	}

	return http.DetectContentType(buffer), nil
}

func SendMultiPartRequest() {
	netClient := &http.Client{
		Timeout: time.Second * 10,
	}

	pipeReader, pipeWriter := io.Pipe()
	writer := multipart.NewWriter(pipeWriter)

	go func() {
		defer pipeWriter.Close()
		defer writer.Close()

		// Open the file
		file, err := os.Open("./main.go")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		fileMimeType, err := GetMimeType(file)
		if err != nil {
			log.Fatal(err)
		}
		// Create a form file
		header := make(textproto.MIMEHeader)
		header.Set("Content-Disposition", `form-data; name="fileField"; filename="main.go"`)
		header.Set("Content-Type", fileMimeType)
		part, err := writer.CreatePart(header)

		if err != nil {
			log.Fatal(err)
		}

		// Copy the file data to the form file
		_, err = io.Copy(part, file)
		if err != nil {
			log.Fatal(err)
		}

		// Write other fields
		_ = writer.WriteField("name", "test")
	}()

	req, err := http.NewRequest("POST", "http://localhost:8080", pipeReader)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := netClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Response Body:", string(respBody))
}

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
				Format: reqxHttp.FORM,
				Data: form.HTTPRequestFormData{
					"name":  "test",
					"value": "test",
				},
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
