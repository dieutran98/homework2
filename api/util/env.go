package util

import (
	"caching/config"
	"encoding/json"

	"github.com/pkg/errors"
)

type redis struct {
	Addr     string `json:"addr"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

type env struct {
	Redis redis `json:"redis"`
}

var myEnv *env

func GetEnv() (env, error) {
	if myEnv == nil {
		file, err := config.GetConfigFile()
		if err != nil {
			return env{}, errors.Wrap(err, "failed get config")
		}
		defer file.Close()
		if err := json.NewDecoder(file).Decode(&myEnv); err != nil {
			return env{}, errors.Wrap(err, "failed decode config")
		}
	}

	return *myEnv, nil
}
