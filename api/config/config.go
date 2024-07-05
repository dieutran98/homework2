package config

import (
	"embed"
	"io"

	"github.com/pkg/errors"
)

//go:embed *.json
var f embed.FS

func GetConfigFile() (io.ReadCloser, error) {
	file, err := f.Open("env.json")
	if err != nil {
		return nil, errors.Wrap(err, "failed to read file")
	}
	return file, err
}
