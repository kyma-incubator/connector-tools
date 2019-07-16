package service

import (
	"github.com/kyma-incubator/connector-tools/qualtrics-webhook-registration/pkg/util"
)

type ReconcilerType interface {
	Reconcile(ctx *util.RequestContext) error
	RefreshQualtricsState(ctx *util.RequestContext) error
	ReconcileState(topicsToRegister []string, subscriptionsToDeregister []string, ctx *util.RequestContext) error
	CompareState(ctx *util.RequestContext) (topicsToRegister []string, subscriptionsToDeregister []string, err error)

}

type TopicConverter interface {
	MapEventTypeVersionToTopic(eventType string, version string) (string, error)
	MapTopicToEventTypeVersion(topic string) (string, string, error)

}
