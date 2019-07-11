package apiclient

import "github.com/kyma-incubator/connector-tools/qualtrics-webhook-registration/pkg/util"

type QualtricsAPIClient interface {
	DeleteSubscription(subscriptionID string, ctx *util.RequestContext) error
	GetSubscriptionList(ctx *util.RequestContext) ([]QualtricsSubscription, error)
	CreateSubscription(subscription *QualtricsSubscription, ctx *util.RequestContext) (string, error)
	UpdateSubscription(subscription *QualtricsSubscription, ctx *util.RequestContext) (string, error)
}

type EventServiceAPIClient interface {
	GetActiveSubscriptions(ctx *util.RequestContext) ([]EventSubscription, error)
}
