package uploader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"privateterraformregistry/internal/modules"
)

type Uploader interface {
	Upload(r *http.Request, m modules.Module) error
}

func New(c *Config) Uploader {
	return &uploader{
		config: c,
	}
}

type Config struct {
	MaxUploadSize int64
	DataDir       string
}

type uploader struct {
	config *Config
}

func (uploader *uploader) Upload(r *http.Request, m modules.Module) error {
	r.ParseMultipartForm(uploader.config.MaxUploadSize)
	file, _, err := r.FormFile("module")

	if err != nil {
		return err
	}
	defer file.Close()

	var filePath = fmt.Sprintf("%s/%s.%s.%s.%s.tar.gz", uploader.config.DataDir, m.Namespace, m.System, m.Name, m.Version)
	newFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)

	if err != nil {
		return err
	}

	defer newFile.Close()
	_, err = io.Copy(newFile, file)

	return err
}
