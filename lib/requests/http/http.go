package http

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"

	"github.com/tpyle/reqx/lib/requests/context"
	reqxForm "github.com/tpyle/reqx/lib/requests/http/form"
	reqxJson "github.com/tpyle/reqx/lib/requests/http/json"
	reqxMultipart "github.com/tpyle/reqx/lib/requests/http/multipart"
)

type HTTPRequestFormat string

const (
	JSON      HTTPRequestFormat = "json"
	FORM      HTTPRequestFormat = "form"
	MULTIPART HTTPRequestFormat = "multipart"
)

type HTTPRequestData interface {
	Serialize(io.WriteCloser, *context.RequestContext, chan error, chan string)
}

type HTTPRequestURL struct {
	Protocol string            `json:"protocol"`
	Hostname string            `json:"hostname"`
	Port     int32             `json:"port"`
	Path     string            `json:"path"`
	Query    map[string]string `json:"query"`
}

type HTTPRequestOptions struct {
	Timeout       time.Duration `json:"timeout,omitempty"`
	SkipTLSVerify bool          `json:"insecure,omitempty"`
}

type HTTPRequestSpec struct {
	Method  string             `json:"method"`
	URL     HTTPRequestURL     `json:"url"`
	Format  HTTPRequestFormat  `json:"format"`
	Data    HTTPRequestData    `json:"data"`
	Options HTTPRequestOptions `json:"options"`
	Headers map[string]string  `json:"headers"`
}

func (s *HTTPRequestSpec) UnmarshalJSON(b []byte) error {
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		panic(err)
	}

	s.Format = HTTPRequestFormat(m["format"].(string))

	switch s.Format {
	case JSON:
		s.Data = reqxJson.JSONData{}
	case FORM:
		s.Data = reqxForm.FormData{}
	case MULTIPART:
		s.Data = reqxMultipart.MultipartFormData{}
	}

	if err := mapstructure.Decode(m["data"], &s.Data); err != nil {
		panic(err)
	}
	return nil
}

func (s HTTPRequestSpec) Send(c *context.RequestContext) error {
	resolvedTimeout := c.HTTPContext.Timeout
	if s.Options.Timeout != 0 {
		resolvedTimeout = s.Options.Timeout
	}
	if resolvedTimeout == 0 {
		// Default to 60 seconds
		resolvedTimeout = time.Second * 60
	}
	netClient := http.Client{
		Timeout: resolvedTimeout,
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: c.HTTPContext.SkipTLSVerify || s.Options.SkipTLSVerify,
			},
		},
	}

	var url string
	if s.URL.Port == 0 {
		url = fmt.Sprintf("%s://%s%s", s.URL.Protocol, s.URL.Hostname, s.URL.Path)
	} else {
		url = fmt.Sprintf("%s://%s:%d%s", s.URL.Protocol, s.URL.Hostname, s.URL.Port, s.URL.Path)
	}
	logrus.Infof("Sending request to %s", url)

	pipeReader, pipeWriter := io.Pipe()
	errChan := make(chan error, 1)
	contentTypeChan := make(chan string)

	go s.Data.Serialize(pipeWriter, c, errChan, contentTypeChan)
	logrus.Debugf("Request: %v", s)

	req, err := http.NewRequest(s.Method, url, pipeReader)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	select {
	case ct := <-contentTypeChan:
		// Set the Content-Type header
		req.Header.Set("Content-Type", ct)
	case err := <-errChan:
		// If an error occurred in the Serialize method, return it
		return fmt.Errorf("error serializing request data: %w", err)
	case <-time.After(c.HTTPContext.Timeout):
		// If the Serialize method takes too long, return an error
		return fmt.Errorf("timeout serializing request data")
	}

	for k, v := range s.Headers {
		if strings.ToLower(k) == "content-type" && s.Format == MULTIPART {
			// Skip setting the Content-Type header if we're sending multipart data, as it needs to be set by the Serialize method
			continue
		}
		req.Header.Set(k, v)
	}
	resp, err := netClient.Do(req)
	logrus.Debugf("Received response: %v", resp)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()
	select {
	case err = <-errChan:
		if err != nil {
			return err
		}
	case <-time.After(c.HTTPContext.Timeout):
		return fmt.Errorf("timeout serializing request data")
	}

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}

	return nil
}
