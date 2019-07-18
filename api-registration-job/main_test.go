package main

import (
	"encoding/json"
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
				RegistrationURL: "test-url",
				app: &oDataWithBasicAuth{
					BasicUser:     "test-auth",
					BasicPassword: "test-pass",
				},
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
			got := tt.r.app.generateMetadata(tt.args.endpoint, tt.r)
			var inter interface{}
			err := json.Unmarshal(got, &inter)
			if err != nil {
				t.Errorf("Failed to parse %s", err)
			}
		})
	}
}
