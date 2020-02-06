package config

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestNewFromFile(t *testing.T) {

	targetConfig := &JobConfig{
		LogLevel: "Info",
		Compass: CompassConfig{
			ClientID:     "ClientID",
			ClientSecret: "ClientSecret",
			DirectorURL:  "https://DirectorURL",
			TenantId:     "TenantId",
			TokenUrl:     "https://TokenUrl",
			TimeoutMills: 2000,
		},
		OpenConnectors: OpenConnectorsConfig{
			Hostname:           "Hostname",
			OrganizationSecret: "OrganizationSecret",
			UserSecret:         "UserSecret",
			TimeoutMills:       2000,
			Tags: []string{
				"test",
				"test2",
			},
		},
	}

	targetConfigJson, _ := json.Marshal(targetConfig)

	if err := ioutil.WriteFile("config.json", targetConfigJson, 777); err != nil {
		t.Fatalf("error writing config file: %s", err.Error())
	}

	configFile, err := os.Open("config.json")
	if err != nil {
		t.Fatalf("error opening config file: %s", err.Error())
	}

	readConfig, err := NewFromFile(context.Background(), configFile)
	if err != nil {
		t.Fatalf("error reading config file: %s", err.Error())
	}

	if !reflect.DeepEqual(*readConfig, *targetConfig)  {
		t.Errorf("expected read config to be equal to target config.\nread config:\n%+v\n"+
			"target config:\n%+v\n", readConfig, targetConfig)
	}

	os.Remove(configFile.Name())
	configFile.Close()
}
