package raw

import (
	"fmt"
	"io"

	"github.com/sirupsen/logrus"
	"github.com/tpyle/reqx/lib/requests/context"
)

type RawData []byte

func (d RawData) Serialize(w io.WriteCloser, c *context.RequestContext, respChannel chan error, contentTypeChannel chan string) {
	defer w.Close()
	defer close(respChannel)
	defer close(contentTypeChannel)

	contentTypeChannel <- "application/octet-stream"
	logrus.Debugf("Serializing RAW data: %+v", d)
	_, err := w.Write(d)
	if err != nil {
		logrus.WithError(err).Error("Failed to serialize RAW data")
		respChannel <- fmt.Errorf("failed to serialize RAW data: %w", err)
	}
	respChannel <- nil
}
