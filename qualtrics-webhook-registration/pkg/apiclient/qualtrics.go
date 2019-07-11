package apiclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/kyma-incubator/connector-tools/qualtrics-webhook-registration/pkg/util"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	qualtricsApiPath      = "/API/v3/eventsubscriptions/"
	qualtricsApiKeyHeader = "X-API-TOKEN"
)

//Pointer to a Qualtrics apiclient
type Qualtrics struct {
	APIKey string
	URL    string
	Client *http.Client
}

type QualtricsSubscription struct {
	ID             string `json:"id,omitempty"`
	Topics         string `json:"topics"`
	PublicationURL string `json:"publicationUrl"`
	SharedKey      string `json:"sharedKey,omitempty"`
}

type subscriptionListResponse struct {
	Result subscriptionListElements `json:"result"`
}

type subscriptionListElements struct {
	Elements []QualtricsSubscription `json:"elements"`
}

type subscriptionCreateResponse struct {
	Result subscriptionCreateResponseResult `json:"result"`
}

type subscriptionCreateResponseResult struct {
	ID string `json:"id"`
}



func NewQualtricsSubscription(apikey string, url string, timeout time.Duration) (*Qualtrics, error) {
	return NewQualtricsSubscriptionWithClient(apikey, url, &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			MaxIdleConns:    10,
			MaxConnsPerHost: 20,
		},
	})
}

func NewQualtricsSubscriptionWithClient(apikey string, url string, client *http.Client) (*Qualtrics, error) {
	if apikey == "" {
		return nil, fmt.Errorf("apikey must not be empty")
	}

	if url == "" {
		return nil, fmt.Errorf("url must not be empty")
	}

	if client == nil {
		return nil, fmt.Errorf("client must not be empty")
	}

	return &Qualtrics{
		APIKey: apikey,
		URL:    url,
		Client: client,
	}, nil
}

func (i *Qualtrics) DeleteSubscription(subscriptionID string, ctx *util.RequestContext) error {

	log.WithFields(ctx.GetLoggerFields()).Debugf("Deleting Subscription %q", subscriptionID)

	url, err := url.Parse(fmt.Sprintf("%s%s%s",i.URL, qualtricsApiPath, subscriptionID))

	if err != nil {
		log.WithFields(ctx.GetLoggerFields()).Errorf("error assembling delete url: %s", err.Error())
		return fmt.Errorf("error assembling delete url: %s", err.Error())
	}

	req, err := http.NewRequest(http.MethodDelete, url.String(), nil)

	if err != nil {
		log.WithFields(ctx.GetLoggerFields()).Errorf("error assembling delete request: %s", err.Error())
		return fmt.Errorf("error assembling delete request: %s", err.Error())
	}

	req.Header.Set(qualtricsApiKeyHeader, i.APIKey)
	ctx.IncludeTraceHeaders(req.Header)

	resp, err := i.Client.Do(req)

	if err != nil {
		log.WithFields(ctx.GetLoggerFields()).Errorf("error deleting subscription: %s", err.Error())
		return  fmt.Errorf("error deleting subscription: %s", err.Error())
	}

	if resp.StatusCode != http.StatusOK {

		if log.GetLevel() != log.TraceLevel {
			log.WithFields(ctx.GetLoggerFields()).Errorf("error deleting subscription: %d (%s)",
				resp.StatusCode, resp.Status)
		} else {
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			defer resp.Body.Close()

			if err != nil {
				bodyBytes = []byte(err.Error())
			}

			log.WithFields(ctx.GetLoggerFields()).Tracef("error deleting subscription: %d (%s): %s",
				resp.StatusCode, resp.Status, string(bodyBytes))
		}
		return fmt.Errorf("error deleting subscription: %d (%s)", resp.StatusCode, resp.Status)
	}
 	return nil
}

func (i *Qualtrics) GetSubscriptionList(ctx *util.RequestContext) ([]QualtricsSubscription, error) {

	log.WithFields(ctx.GetLoggerFields()).Debug("Reading Qualtrics Subscriptions")

	url, err := url.Parse(fmt.Sprintf("%s%s",i.URL, qualtricsApiPath))

	if err != nil {
		log.WithFields(ctx.GetLoggerFields()).Errorf("error assembling get subscription list url: %s", err.Error())
		return nil, fmt.Errorf("error assembling get subscription list url: %s", err.Error())
	}

	req, err := http.NewRequest(http.MethodGet, url.String(), nil)

	if err != nil {
		log.WithFields(ctx.GetLoggerFields()).Errorf("error assembling get subscription list request: %s",
			err.Error())
		return nil, fmt.Errorf("error assembling get subscription list request: %s", err.Error())
	}

	req.Header.Set(qualtricsApiKeyHeader, i.APIKey)

	ctx.IncludeTraceHeaders(req.Header)

	resp, err := i.Client.Do(req)

	if err != nil {
		log.WithFields(ctx.GetLoggerFields()).Errorf("error getting subscription list: %s", err.Error())
		return nil, fmt.Errorf("error getting subscription list: %s", err.Error())
	}

	if resp.StatusCode != http.StatusOK {

		if log.GetLevel() != log.TraceLevel {
			log.WithFields(ctx.GetLoggerFields()).Errorf("error getting subscription list: %d (%s)",
				resp.StatusCode, resp.Status)
		} else {
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			defer resp.Body.Close()

			if err != nil {
				bodyBytes = []byte(err.Error())
			}

			log.WithFields(ctx.GetLoggerFields()).Tracef("error getting subscription list: %d (%s): %s",
				resp.StatusCode, resp.Status, string(bodyBytes))
		}
		return nil, fmt.Errorf("error getting subscription list: %d (%s)",
			resp.StatusCode, resp.Status)
	}

	var respJson subscriptionListResponse

	dec := json.NewDecoder(resp.Body)
	defer resp.Body.Close()

	err = dec.Decode(&respJson)

	if err != nil {
		log.WithFields(ctx.GetLoggerFields()).Errorf("error getting parsing subscription list: %s", err.Error())
		return nil, fmt.Errorf("error getting parsing subscription list: %s", err.Error())
	}

	return respJson.Result.Elements, nil
}

func (i *Qualtrics) CreateSubscription(subscription *QualtricsSubscription, ctx *util.RequestContext) (string, error) {

	if log.GetLevel() != log.TraceLevel {
		log.WithFields(ctx.GetLoggerFields()).Debugf("Creating Subscription for topic %s", subscription.Topics)
	} else {
		log.WithFields(ctx.GetLoggerFields()).Tracef("Creating Subscription for topic %s: %+v",
			subscription.Topics, subscription)
	}

	url, err := url.Parse(fmt.Sprintf("%s%s",i.URL, qualtricsApiPath))

	if err != nil {
		log.WithFields(ctx.GetLoggerFields()).Errorf("error assembling create subscription url: %s", err.Error())
		return "", fmt.Errorf("error assembling create subscription url: %s", err.Error())
	}

	subscriptionByte, err := json.Marshal(subscription)

	if err != nil {
		log.WithFields(ctx.GetLoggerFields()).Errorf("error assembling create subscription body: %s", err.Error())
		return "", fmt.Errorf("error assembling create subscription body: %s", err.Error())
	}

	req, err := http.NewRequest(http.MethodPost, url.String(), bytes.NewReader(subscriptionByte))

	if err != nil {
		log.WithFields(ctx.GetLoggerFields()).Errorf("error assembling create subscription request: %s",
			err.Error())
		return "", fmt.Errorf("error assembling create subscription request: %s",
			err.Error())
	}

	req.Header.Set(qualtricsApiKeyHeader, i.APIKey)
	req.Header.Set("Content-Type", "application/json")

	ctx.IncludeTraceHeaders(req.Header)

	resp, err := i.Client.Do(req)

	if err != nil {
		log.WithFields(ctx.GetLoggerFields()).Errorf("error creating subscription: %s", err.Error())
		return "", fmt.Errorf("error creating subscription: %s", err.Error())
	}

	if resp.StatusCode != http.StatusOK {

		if log.GetLevel() != log.TraceLevel {
			log.WithFields(ctx.GetLoggerFields()).Errorf("error creating subscription: %d (%s)",
				resp.StatusCode, resp.Status)
		} else {
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			defer resp.Body.Close()

			if err != nil {
				bodyBytes = []byte(err.Error())
			}

			log.WithFields(ctx.GetLoggerFields()).Tracef("error creating subscription: %d (%s): %s",
				resp.StatusCode, resp.Status, string(bodyBytes))
		}
		return "", fmt.Errorf("error creating subscription: %d (%s)", resp.StatusCode, resp.Status)
	}

	var respJson subscriptionCreateResponse

	dec := json.NewDecoder(resp.Body)
	defer resp.Body.Close()

	err = dec.Decode(&respJson)

	if err != nil {
		log.WithFields(ctx.GetLoggerFields()).Errorf("error getting parsing create subscription response: %s",
			err.Error())
		return "", fmt.Errorf("error getting parsing create subscription response: %s",
			err.Error())
	}

	return respJson.Result.ID, nil
}

func (i *Qualtrics) UpdateSubscription(subscription *QualtricsSubscription, ctx *util.RequestContext) (string, error) {

	if log.GetLevel() != log.TraceLevel {
		log.WithFields(ctx.GetLoggerFields()).Debugf("Updating Subscription for topic %s with id %s",
			subscription.Topics, subscription.ID)
	} else {
		log.WithFields(ctx.GetLoggerFields()).Tracef("Updating Subscription for topic %s with id %s: %+v",
			subscription.Topics, subscription.ID, subscription)
	}

	//Update is create followed by delete as no matching endpoint is offered

	id, err := i.CreateSubscription(subscription, ctx)

	if err != nil {
		log.WithFields(ctx.GetLoggerFields()).Errorf("error updating subscription whilst " +
			"creating new one: %s", err.Error())
		return "", fmt.Errorf("error updating subscription whilst " +
			"creating new one: %s", err.Error())
	}


	err = i.DeleteSubscription(subscription.ID, ctx)

	if err != nil {
		log.WithFields(ctx.GetLoggerFields()).Errorf("error updating subscription whilst " +
			"deleting old one: %s", err.Error())
		return "", fmt.Errorf("error updating subscription whilst " +
			"deleting old one: %s", err.Error())
	}

	return id, nil
}
