package multipart

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path"

	"github.com/sirupsen/logrus"
	"github.com/tpyle/reqx/lib/requests/context"
)

type MultipartFormField struct {
	Name     string `json:"name"`
	Value    string `json:"value,omitempty"`
	FileName string `json:"filename,omitempty"`
}

type MultipartFormData []MultipartFormField

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

func (d MultipartFormData) Serialize(w io.WriteCloser, c *context.RequestContext, respChannel chan error, contentTypeChannel chan string) {
	defer w.Close()
	defer close(respChannel)
	defer close(contentTypeChannel)

	writer := multipart.NewWriter(w)
	defer writer.Close()
	contentTypeChannel <- writer.FormDataContentType()

	for _, field := range d {
		if field.FileName == "" {
			err := writer.WriteField(field.Name, field.Value)
			if err != nil {
				logrus.WithError(err).Error("Failed to write field")
				respChannel <- err
				return
			}
		} else {
			// Open the file
			filePath := path.Join(c.FileLocation, field.FileName)
			file, err := os.Open(filePath)
			if err != nil {
				logrus.WithError(err).Error("Failed to open file")
				respChannel <- err
				return
			}
			defer file.Close()

			fileMimeType, err := GetMimeType(file)
			if err != nil {
				logrus.WithError(err).Error("Failed to get file mime type")
				respChannel <- err
				return
			}
			// Create a form file
			header := make(textproto.MIMEHeader)
			contentDisposition := fmt.Sprintf(`form-data; name="%s"; filename="%s"`, field.Name, field.FileName)
			header.Set("Content-Disposition", contentDisposition)
			header.Set("Content-Type", fileMimeType)
			part, err := writer.CreatePart(header)

			if err != nil {
				logrus.WithError(err).Error("Failed to create part")
				respChannel <- err
				return
			}

			// Copy the file data to the form file
			_, err = io.Copy(part, file)
			if err != nil {
				logrus.WithError(err).Error("Failed to copy file data")
				respChannel <- err
				return
			}
		}
	}
}
