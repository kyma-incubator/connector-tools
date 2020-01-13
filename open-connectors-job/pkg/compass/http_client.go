package compass

import (
	"context"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"net/http"
	"time"
)

const (
	scope_app_read  = "application:read"
	scope_app_write = "application:write"
)

//Creates new HTTP Client with OAuth2 Client Credentials grant type. The client is based on the default
//and augmented with a timeout value
func CreateHttpClientTimeout(ctx context.Context, clientID string, clientSecret string,
	tokenUrl string, timeoutMills int) (*http.Client, error) {


	if log.GetLevel() == log.TraceLevel {
		log.Tracef("creating new compass http client with clientID %q, %d digit client secret, " +
			"tokenUrl %q and timeout %d", clientID, len(clientSecret), tokenUrl, timeoutMills)
	} else {
		log.Debug("creating new compass http client with timeout")
	}

	httpClient := &http.Client{
		Timeout: time.Duration(timeoutMills) * time.Millisecond,
	}
	ctx = context.WithValue(ctx, oauth2.HTTPClient, httpClient)

	return CreateHttpClient(ctx, clientID, clientSecret, tokenUrl)

}


//Creates new HTTP Client with OAuth2 Client Credentials grant type. The client is based on the
//http.Client supplied in the context under oauth2.HTTPClient (or default if nothing is specified)
func CreateHttpClient(ctx context.Context, clientID string, clientSecret string,
	tokenUrl string) (*http.Client, error) {

	config := clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     tokenUrl,
		Scopes:       []string{scope_app_read, scope_app_write},
	}

	client := config.Client(ctx)

	if log.GetLevel() == log.TraceLevel {
		log.Tracef("creating new compass http client with clientID %q, %d digit client secret, " +
			"tokenUrl %q and custom httpClient %t", clientID, len(clientSecret), tokenUrl,
			ctx.Value(oauth2.HTTPClient) == true)
	} else {
		log.Debug("creating new compass http client")
	}

	return client, nil
}
