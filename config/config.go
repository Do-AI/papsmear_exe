package config

import (
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type papSmearConfig struct {
	Credential struct {
		ClientID     string `yaml:"client-id"`
		ClientSecret string `yaml:"client-secret"`
		OrgNo        int    `yaml:"orgNo"`
		AgentNo      int    `yaml:"agentNo"`
		Username     string `yaml:"username"`
	} `yaml:"credential"`
	Db struct {
		Path   string `yaml:"path"`
		Schema string `yaml:"schema"`
	} `yaml:"db"`
	Log struct {
		Path string `yaml:"path"`
	} `yaml:"log"`
	Goroutine int `yaml:"goroutine"`
	Svc       struct {
		URL    string `yaml:"url"`
		APIURL string `yaml:"api-url"`
	} `yaml:"svc"`
	Slack struct {
		Token   string `yaml:"token"`
		Channel string `yaml:"channel"`
	} `yaml:"slack"`
	Folder     string   `yaml:"folder"`
	Extensions []string `yaml:"extensions"`
	Patch      struct {
		Level     []int32 `yaml:"level"`
		Size      int64   `yaml:"size"`
		SizeFront int64   `yaml:"size_front"`
		Overlap   int64   `yaml:"overlap"`
		Quality   int64   `yaml:"quality"`
	} `yaml:"patch"`
}

var Config papSmearConfig

func readConfig() {
	fileName, _ := filepath.Abs("config/config.yaml")
	yamlFile, err := ioutil.ReadFile(fileName)
	err = yaml.Unmarshal(yamlFile, &Config)

	if err != nil {
		panic(err)
	}
}

func init() {
	readConfig()
}
