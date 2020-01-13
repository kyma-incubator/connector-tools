package compass

import (
	"context"
	"encoding/json"
	"golang.org/x/oauth2"
	"net/http"
	"net/http/httptest"
	"testing"
)



func TestCreateHttpClient(t *testing.T) {

	clientSecret := "secret123"
	clientID := "iamyourclient"

	authServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		err := r.ParseForm()

		if err != nil {
			t.Fatalf("no valid OAuth2 request received, body invalid: %s", err.Error())
		}

		targetScopes := scope_app_read + " " + scope_app_write
		if r.Form.Get("scope") != targetScopes {
			t.Errorf("invalid scope request, received %q shhould have received %q",
				r.Form.Get("scope"), targetScopes)
		}

		uname, pwd, ok := r.BasicAuth()
		if !ok {
			t.Error("no valid OAuth2 request received, authorization header invalid")
		}

		if uname != clientID || pwd != clientSecret {
			t.Error("invalid Client ID and or Client Secret provided")
		}

		jsonResponseMap := map[string]interface{}{
			"access_token": "token",
			"token_type":   "Bearer",
			"expires_in":   10,
		}

		responseBytes, _ := json.Marshal(&jsonResponseMap)

		w.Header().Set("Content-Type", "application/json")
		w.Write(responseBytes)

	}))
	defer authServer.Close()

	apiServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			t.Fatalf("no authorization header provided")
		}

		if authHeader != "Bearer token" {
			t.Errorf("auth header should be \"Bearer token\" but was %q", authHeader)
		}

		w.WriteHeader(200)
	}))
	defer apiServer.Close()

	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, authServer.Client())

	httpClient, err := CreateHttpClient(ctx, clientID, clientSecret, authServer.URL)

	if err != nil {
		t.Fatalf("no valid OAuth2 http client created: %s", err.Error())
	}

	_, err = httpClient.Get(apiServer.URL)

	if err != nil {
		t.Fatalf("request to mock server failed with: %s", err.Error())
	}
}