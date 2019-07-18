package main

import "testing"

func Test_getAPIURL(t *testing.T) {
	app := restWithAPIKey{
		apikey: "sdd",
		source: "source",
	}

	tests := []struct {
		name      string
		systemURL string
		path      string
	}{
		{
			name:      "no prefix suffix",
			systemURL: "https://system.com/x1.svc",
			path:      "users",
		},
		{
			name:      "both prefix suffix",
			systemURL: "https://system.com/x1.svc/",
			path:      "/users",
		},
		{
			name:      "only prefix",
			systemURL: "https://system.com/x1.svc",
			path:      "/users",
		},
		{
			name:      "only suffix",
			systemURL: "https://system.com/x1.svc/",
			path:      "users",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expected := "https://system.com/x1.svc/users"

			apiURL := app.getAPIUrl(tt.systemURL, tt.path)
			if apiURL != expected {
				t.Errorf("failed, expected %s, got %s", expected, apiURL)
			}
		})
	}
}
