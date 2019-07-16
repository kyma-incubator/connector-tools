package service

import (
	"errors"
	"github.com/kyma-incubator/connector-tools/qualtrics-webhook-registration/pkg/apiclient"
	"github.com/kyma-incubator/connector-tools/qualtrics-webhook-registration/pkg/util"
	"net/http"
	"testing"
)

type qualtricsAPICLientMock struct {
	QualtricsSubscriptionListError bool
}

type eventServiceAPIClientMock struct {
	Fail bool
}

func (q *qualtricsAPICLientMock) DeleteSubscription(subscriptionID string, ctx *util.RequestContext) error {

	if subscriptionID != "SUB_8BXPypz7dWUm8e1" {
		return errors.New("wrong subscription deleted")
	}
	return nil
}
func (q *qualtricsAPICLientMock) GetSubscriptionList(ctx *util.RequestContext) ([]apiclient.QualtricsSubscription, error) {
	if q.QualtricsSubscriptionListError {
		return nil, errors.New("error result")
	}

	return []apiclient.QualtricsSubscription{
		{
			ID:             "SUB_8BXPypz7dWUm8e1",
			Topics:         "controlpanel.activateSurvey",
			PublicationURL: "https://kyma-project.io/qualtrics",
			SharedKey:      "",
		},
		{
			ID:             "SUB_8BXPypz7dWUm8e1",
			Topics:         "surveyengine.completedResponse.*",
			PublicationURL: "https://kyma-project.io/qualtrics",
			SharedKey:      "",
		},
		{
			ID:             "SUB_8BXPypz7dWUm8e1",
			Topics:         "surveyengine.completedResponse.*",
			PublicationURL: "https://kyma-project2.io/qualtrics",
			SharedKey:      "",
		},
		{
			ID:             "SUB_8BXPypz7dWUm8e1",
			Topics:         "controlpanel.activateSurvey",
			PublicationURL: "https://kyma-project2.io/qualtrics",
			SharedKey:      "",
		},
		{
			ID:             "SUB_8BXPypz7dWUm8e1",
			Topics:         "*",
			PublicationURL: "https://kyma-project3.io/qualtrics",
			SharedKey:      "",
		},
	}, nil
}
func (q *qualtricsAPICLientMock) CreateSubscription(subscription *apiclient.QualtricsSubscription,
	ctx *util.RequestContext) (string, error) {

	if subscription.Topics != "controlpanel.deactivateSurvey" {
		return "", errors.New("wrong topic to register")
	}
	return "SUB_8pSh6Kbqjgvwynb", nil
}
func (q *qualtricsAPICLientMock) UpdateSubscription(subscription *apiclient.QualtricsSubscription,
	ctx *util.RequestContext) (string, error) {

	//error as not ready for test
	return "", errors.New("Not implemented")
}

func (e eventServiceAPIClientMock) GetActiveSubscriptions(ctx *util.RequestContext) ([]apiclient.EventSubscription, error) {
	if e.Fail {
		return nil, errors.New("something went wrong")
	}
	return []apiclient.EventSubscription{
		{
			EventType:    "surveyengine.completedResponse",
			EventVersion: "v1",
		},
		{
			EventType:    "controlpanel.deactivateSurvey",
			EventVersion: "v1",
		},
	}, nil
}

func TestNewReconciler(t *testing.T) {

	topicConverter, err := NewTopicmapper("../../testdata/topic-config.json")

	if err != nil {
		t.Errorf("topicConverter creation must not fail, error %q", err.Error())
	}

	inst, err := NewReconciler(&qualtricsAPICLientMock{}, &eventServiceAPIClientMock{}, topicConverter,
		"dummy", "https://www.kyma-project.io")

	if err != nil {
		t.Errorf("instance must not fail on creation, error %q", err.Error())
	}

	if inst.SubscriptionURL != "https://www.kyma-project.io" {
		t.Errorf("url  should be https://www.kyma-project.io, but was %q", inst.SubscriptionURL)
	}

	if inst.QualtricsAPIClient == nil {
		t.Error("QualtricsAPIClient should not be nil, but was")
	}

	if inst.EventServiceAPIClient == nil {
		t.Error("EventServiceAPIClient should not be nil, but was")
	}

	if inst.TopicConverter == nil {
		t.Error("TopicConverter should not be nil, but was")
	}

	if inst.sharedKey != "dummy" {
		t.Error("sharedKey should be \"dummy\", but was not")
	}

	_, err = NewReconciler(nil, &eventServiceAPIClientMock{}, topicConverter, "dummy",
		"https://www.kyma-project.io")

	if err == nil {
		t.Error("instance should fail on creation, but did not")
	}

	_, err = NewReconciler(&qualtricsAPICLientMock{}, nil, topicConverter,"dummy",
		"https://www.kyma-project.io")

	if err == nil {
		t.Error("instance should fail on creation, but did not")
	}

	_, err = NewReconciler(&qualtricsAPICLientMock{}, &eventServiceAPIClientMock{}, topicConverter,"dummy",
		"")

	if err == nil {
		t.Error("instance should fail on creation, but did not")
	}

	inst, err = NewReconciler(&qualtricsAPICLientMock{}, &eventServiceAPIClientMock{}, nil, "dummy",
		"https://www.kyma-project.io")

	if err == nil {
		t.Error("instance should fail on creation, but did not")
	}
}

func TestReconciler_RefreshQualtricsState(t *testing.T) {

	topicConverter, err := NewTopicmapper("../../testdata/topic-config.json")

	if err != nil {
		t.Errorf("topicConverter creation must not fail, error %q", err.Error())
	}

	// success
	inst, err := NewReconciler(&qualtricsAPICLientMock{}, &eventServiceAPIClientMock{}, topicConverter, "dummy",
		"https://kyma-project.io/qualtrics")

	if err != nil {
		t.Errorf("instance must not fail on creation, error %q", err.Error())
	}

	err = inst.RefreshQualtricsState(&util.RequestContext{TraceHeaders: http.Header{}})

	if err != nil {
		t.Errorf("refreshing state from Qualtrics failed with error, but shouldn't: %s", err.Error())
		t.Fail()
	}

	if len(inst.qualtricsEventsToSubscriptions) != 2 {
		t.Errorf("expected to receive only 2 subscriptions, received %d", len(inst.qualtricsEventsToSubscriptions))
	}

	if _, ok := inst.qualtricsEventsToSubscriptions["surveyengine.completedResponse.v1"]; !ok {
		t.Error("expected to receive subscription for event surveyengine.completedResponse.v1")
	}

	if _, ok := inst.qualtricsEventsToSubscriptions["controlpanel.activateSurvey.v1"]; !ok {
		t.Error("expected to receive subscription for event controlpanel.activateSurvey.v1")
	}

	//error

	qualtricsMock := &qualtricsAPICLientMock{}
	inst, err = NewReconciler(qualtricsMock, &eventServiceAPIClientMock{},
		topicConverter, "dummy", "https://kyma-project.io/qualtrics")

	if err != nil {
		t.Fatalf("instance must not fail on creation, error %q", err.Error())
	}

	//make qualtrics mock return error
	qualtricsMock.QualtricsSubscriptionListError = true

	err = inst.RefreshQualtricsState(&util.RequestContext{TraceHeaders: http.Header{}})

	if err == nil {
		t.Error("expected qualtrics refresh to fail, but it didn't")
	}
}

func TestReconciler_CompareState(t *testing.T) {
	topicConverter, err := NewTopicmapper("../../testdata/topic-config.json")

	if err != nil {
		t.Errorf("topicConverter creation must not fail, error %q", err.Error())
	}

	// success
	inst, err := NewReconciler(&qualtricsAPICLientMock{}, &eventServiceAPIClientMock{}, topicConverter,
		"dummy", "https://kyma-project.io/qualtrics")

	if err != nil {
		t.Errorf("instance must not fail on creation, error %q", err.Error())
	}

	topicsToRegister, subscriptionsToDeregister, err :=
		inst.CompareState(&util.RequestContext{TraceHeaders: http.Header{}})

	if err != nil {
		t.Errorf("comparing state between kyma & qualtrics failed with error, but shouldn't: %s", err.Error())
		t.Fail()
	}

	if len(topicsToRegister) != 1 {
		t.Fatalf("should have returned 1 topic to register, but returned %+v", topicsToRegister)
	}

	if topicsToRegister[0] != "controlpanel.deactivateSurvey" {
		t.Fatalf("should have returned topic \"controlpanel.deactivateSurvey\" to register, but returned %+v",
			topicsToRegister)
	}

	if len(subscriptionsToDeregister) != 1 {
		t.Fatalf("should have returned 1 subscription to deregister, but returned %+v",
			subscriptionsToDeregister)
	}
	if subscriptionsToDeregister[0] != "SUB_8BXPypz7dWUm8e1" {
		t.Fatalf("should have returned subscription \"SUB_8BXPypz7dWUm8e1\" to deregister, but returned %+v",
			subscriptionsToDeregister)
	}
}

func TestReconciler_ReconcileState(t *testing.T) {
	topicConverter, err := NewTopicmapper("../../testdata/topic-config.json")

	if err != nil {
		t.Errorf("topicConverter creation must not fail, error %q", err.Error())
	}

	// success
	inst, err := NewReconciler(&qualtricsAPICLientMock{}, &eventServiceAPIClientMock{}, topicConverter,
		"dummy", "https://kyma-project.io/qualtrics")

	if err != nil {
		t.Errorf("instance must not fail on creation, error %q", err.Error())
	}

	topicsToRegister := []string{"controlpanel.deactivateSurvey"}
	subscriptionsToDeregister:= []string{ "SUB_8BXPypz7dWUm8e1"}

	err = inst.ReconcileState(topicsToRegister, subscriptionsToDeregister,
		&util.RequestContext{http.Header{}})

	if err != nil {
		t.Errorf("reconciling state between kyma & qualtrics failed with error, but shouldn't: %s", err.Error())
	}

	if _, ok := inst.qualtricsEventsToSubscriptions["controlpanel.deactivateSurvey.v1"]; !ok {
		t.Error("reconciling state between kyma & qualtrics failed, event \"controlpanel." +
			"deactivateSurvey.v1\" should be registered")
	}

	if _, ok := inst.qualtricsSubscriptionsToEvents["SUB_8BXPypz7dWUm8e1"]; ok {
		t.Error("reconciling state between kyma & qualtrics failed, subscription \"SUB_8BXPypz7dWUm8e1\" " +
			"should not be registered")
	}

	//test error

	topicsToRegister = []string{"error"}

	err = inst.ReconcileState(topicsToRegister, subscriptionsToDeregister,
		&util.RequestContext{http.Header{}})

	if err == nil {
		t.Errorf("reconciling state between kyma & qualtrics did not fail with error, but should have")
	}

	subscriptionsToDeregister = []string{"error"}

	err = inst.ReconcileState(topicsToRegister, subscriptionsToDeregister,
		&util.RequestContext{http.Header{}})

	if err == nil {
		t.Errorf("reconciling state between kyma & qualtrics did not fail with error, but should have")
	}
}

func TestReconciler_Reconcile(t *testing.T) {
	topicConverter, err := NewTopicmapper("../../testdata/topic-config.json")

	if err != nil {
		t.Errorf("topicConverter creation must not fail, error %q", err.Error())
	}

	// success
	inst, err := NewReconciler(&qualtricsAPICLientMock{}, &eventServiceAPIClientMock{}, topicConverter,
		"dummy", "https://kyma-project.io/qualtrics")

	if err != nil {
		t.Errorf("instance must not fail on creation, error %q", err.Error())
	}

	err = inst.Reconcile(&util.RequestContext{http.Header{}})

	if err != nil {
		t.Errorf("overall reconciler test failed with error: %s", err.Error())
	}
}
