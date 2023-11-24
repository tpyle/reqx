package http

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

type HTTPRequestSpec struct {
	Method string            `json:"method"`
	URL    HTTPRequestURL    `json:"url"`
	Format HTTPRequestFormat `json:"format"`
	Data   HTTPRequestData   `json:"data"`
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
	netClient := http.Client{
		Timeout: c.HTTPContext.Timeout,
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: c.HTTPContext.SkipTLSVerify,
			},
		},
	}
	url := fmt.Sprintf("%s://%s:%d%s", s.URL.Protocol, s.URL.Hostname, s.URL.Port, s.URL.Path)
	logrus.Infof("Sending request to %s", url)

	pipeReader, pipeWriter := io.Pipe()
	errChan := make(chan error)
	contentTypeChan := make(chan string)

	go s.Data.Serialize(pipeWriter, c, errChan, contentTypeChan)
	logrus.Debugf("Request: %v", s)

	req, err := http.NewRequest(s.Method, url, pipeReader)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", <-contentTypeChan)
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
