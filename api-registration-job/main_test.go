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
				ApplicationName: "test-app",
				ProviderName:    "test-provider",
				ProductName:     "test-product",
				SystemURL:       "https://test-hostname.com",
				BasicUser:       "test-auth",
				BasicPassword:   "test-pass",
				RegistrationURL: "test-url",
			},
			args: args{
				endpoint: endpointInfo{
					Path:        "/test-api",
					Name:        "test-api",
					Description: "test-description",
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
