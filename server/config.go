package server

import (
	"encoding/json"
	"io/ioutil"
)

var (
	defaultCfg = config{
		Port:               "4000",
		UploadPath:         "./upload",
		ElasticSearchURL:   "",
		ElasticSearchIndex: "",
	}
)

type config struct {
	Port               string `json:"port"`
	UploadPath         string `json:"upload_path"`
	ElasticSearchURL   string `json:"elastic_search_url"`
	ElasticSearchIndex string `json:"elastic_search_index"`
}

func loadConfig() *config {
	config := &defaultCfg
	file, err := ioutil.ReadFile("./config.json")
	if err == nil {
		if err := json.Unmarshal(file, &config); err != nil {
			panic(err)
		}
	}

	if err != nil {
		panic("unable to read config file : " + err.Error())
	}

	return config
}
