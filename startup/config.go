package startup

import (
	"github.com/Yni9ht/qqbot/api"
	"github.com/Yni9ht/qqbot/vars"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"net/http"
)

func LoadConfig() error {
	configFileName := "config.yaml"
	b, err := ioutil.ReadFile(configFileName)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(b, vars.Config)
	if err != nil {
		return err
	}
	return nil
}

func LoadAPIClient() error {
	vars.OpenAPI = &api.OpenAPI{
		Token:   vars.Config.GetToken(),
		Client:  http.DefaultClient,
		Sandbox: true,
	}
	return nil
}
