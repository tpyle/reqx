package form

import (
	"io"
	"net/url"

	"github.com/sirupsen/logrus"
	"github.com/tpyle/reqx/lib/requests/context"
)

type HTTPRequestFormData map[string]string

func (d HTTPRequestFormData) Serialize(w io.WriteCloser, c *context.RequestContext, respChannel chan error) {
	defer w.Close()
	defer close(respChannel)
	logrus.Debugf("Serializing form data: %+v", d)
	formData := url.Values{}
	for k, v := range d {
		formData.Set(k, v)
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
