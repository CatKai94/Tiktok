package tools

import (
	"bufio"
	"github.com/goccy/go-json"
	"os"
)

type Config struct {
	AppName string `json:"app_name"`
	AppMode string `json:"app_mode"`
	AppHost string `jsom:"app_host"`
	AppPort string `json:"app_port"`
}

var _cfg *Config = nil

func ParseConfig(path string) (*Config, error) {
	file, err := os.Open(path)

	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	decoder := json.NewDecoder(reader)
	decoder.Decode(_cfg)

	if err = decoder.Decode(&_cfg); err != nil { //参数时interfece类型，参数需要传入地址
		return nil, err
	}
	return _cfg, nil
}
