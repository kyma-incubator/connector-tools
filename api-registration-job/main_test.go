package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"testing"
)

func Test_registrationApp_generateMetadata(t *testing.T) {
	type args struct {
		endpoint endpointInfo
	}
	tests := []struct {
		name string
		r    registrationApp
		args args
	}{
		// TODO: Add test cases.
		{
			name: "render correctly",
			r: registrationApp{
				provider:     "test-provider",
				product:      "test-product",
				hostname:     "https://test-hostname.com",
				authUsername: "test-auth",
				authPassword: "test-pass",
				apiURL:       "test-url",
			},
			args: args{
				endpoint: endpointInfo{
					API:          "/test-api",
					Name:         "test-api",
					Description:  "test-description",
					HelpDoc:      "test-doc",
					Scenario:     "test-scenario",
					ScenarioName: "test-scenaario-name",
				},
			},
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.r.generateMetadata(tt.args.endpoint)
			var inter interface{}
			err := json.Unmarshal(got, &inter)
			if err != nil {
				t.Errorf("Failed to parse %s", err)
			}
		})
	}
}

func TestProtocol(t *testing.T) {
	given := "hostname.com"
	u, _ := url.Parse(given)
	if u.Scheme == "" {
		u.Scheme = "https"
	}

	fmt.Println(u.String())

}
