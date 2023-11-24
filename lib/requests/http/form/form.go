package form

import (
	"encoding/json"
	"io"
	"net/url"

	"github.com/sirupsen/logrus"
	"github.com/tpyle/reqx/lib/requests/context"
)

type StringOrArray []string

func (soa *StringOrArray) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err == nil {
		*soa = []string{s}
		return nil
	}
	var sa []string
	if err := json.Unmarshal(b, &sa); err == nil {
		*soa = sa
		return nil
	}
	return nil
}

type HTTPRequestFormData map[string]StringOrArray

func (d HTTPRequestFormData) Serialize(w io.WriteCloser, c *context.RequestContext, respChannel chan error) {
	defer w.Close()
	defer close(respChannel)
	logrus.Debugf("Serializing form data: %+v", d)
	formData := url.Values{}
	for k, v := range d {
		for _, v2 := range v {
			formData.Add(k, v2)
		}
	}
	_, err := w.Write([]byte(formData.Encode()))
	if err != nil {
		logrus.WithError(err).Error("Failed to serialize form data")
		respChannel <- err
	}
}

func (d HTTPRequestFormData) GetContentType() string {
	return "application/x-www-form-urlencoded"
}
