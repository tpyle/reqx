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

type FormData map[string]StringOrArray

func (d FormData) Serialize(w io.WriteCloser, c *context.RequestContext, respChannel chan error, contentTypeChannel chan string) {
	defer w.Close()
	defer close(respChannel)
	defer close(contentTypeChannel)
	contentTypeChannel <- "application/x-www-form-urlencoded"
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
	respChannel <- nil
}
