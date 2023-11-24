package json

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/sirupsen/logrus"
	"github.com/tpyle/reqx/lib/requests/context"
)

type HTTPRequestJSONData map[string]interface{}

func (d HTTPRequestJSONData) Serialize(w io.WriteCloser, c *context.RequestContext, respChannel chan error) {
	defer w.Close()
	defer close(respChannel)
	logrus.Debugf("Serializing JSON data: %+v", d)
	err := json.NewEncoder(w).Encode(d)
	if err != nil {
		logrus.WithError(err).Error("Failed to serialize JSON data")
		respChannel <- fmt.Errorf("failed to serialize JSON data: %w", err)
	}
}

func (d HTTPRequestJSONData) GetContentType() string {
	return "application/json"
}
