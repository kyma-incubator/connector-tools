package service

import (
	"fmt"
	"github.com/kyma-incubator/connector-tools/qualtrics-webhook-registration/pkg/apiclient"
	"github.com/kyma-incubator/connector-tools/qualtrics-webhook-registration/pkg/util"
	log "github.com/sirupsen/logrus"
	"net/http"
	"sync"
)

type Reconciler struct {
	QualtricsAPIClient             apiclient.QualtricsAPIClient
	EventServiceAPIClient          apiclient.EventServiceAPIClient
	SubscriptionURL                string
	TopicConverter                 TopicConverter
	sharedKey                      string
	qualtricsEventsToSubscriptions map[string]string
	qualtricsSubscriptionsToEvents map[string]string
	mapAccess                      *sync.Mutex
}

func NewReconciler(qualtricsAPIClient apiclient.QualtricsAPIClient,
	eventServiceAPIClient apiclient.EventServiceAPIClient, topicConverter TopicConverter,
	sharedKey string, subscriptionURL string) (*Reconciler, error) {

	if qualtricsAPIClient == nil {
		return nil, fmt.Errorf("qualtricsAPIClient must not be nil")
	}

	if eventServiceAPIClient == nil {
		return nil, fmt.Errorf("eventServiceAPIClient must not be nil")
	}

	if subscriptionURL == "" {
		return nil, fmt.Errorf("subscriptionURL must not be empty")
	}

	if topicConverter == nil {
		return nil, fmt.Errorf("topicConverter must not be nil")
	}

	reconciler := &Reconciler{
		QualtricsAPIClient:             qualtricsAPIClient,
		EventServiceAPIClient:          eventServiceAPIClient,
		SubscriptionURL:                subscriptionURL,
		TopicConverter:                 topicConverter,
		sharedKey:                      sharedKey,
		qualtricsEventsToSubscriptions: make(map[string]string),
		qualtricsSubscriptionsToEvents: make(map[string]string),
		mapAccess: 						&sync.Mutex{},
	}

	err := reconciler.RefreshQualtricsState(&util.RequestContext{TraceHeaders: http.Header{}})

	if err != nil {
		log.Errorf("Error refreshing subscriptions from qualtrics after instantiation: %s", err.Error())
		return nil, err
	}

	return reconciler, nil
}

func (r *Reconciler) RefreshQualtricsState(ctx *util.RequestContext) error {
	log.WithFields(ctx.GetLoggerFields()).Debug("refreshing subscription state on qualtrics")

	subscriptions, err := r.QualtricsAPIClient.GetSubscriptionList(ctx)

	if err != nil {
		log.WithFields(ctx.GetLoggerFields()).Errorf("Error reading subscriptions from qualtrics: %s",
			err.Error())
		return fmt.Errorf("Error reading subscriptions from qualtrics: %s",
			err.Error())
	}

	relevantSubscriptions := make(map[string]string)
	relevantEvents := make(map[string]string)

	for i, _ := range subscriptions {
		//filter for target url
		if subscriptions[i].PublicationURL == r.SubscriptionURL {
			event, version, err := r.TopicConverter.MapTopicToEventTypeVersion(subscriptions[i].Topics)

			if err != nil {
				log.WithFields(ctx.GetLoggerFields()).Errorf("error converting topic %s to event type",
					subscriptions[i].PublicationURL)
			} else {
				relevantEvents[fmt.Sprintf("%s.%s", event, version)] = subscriptions[i].ID
				relevantSubscriptions[subscriptions[i].ID] = fmt.Sprintf("%s.%s", event, version)
				log.WithFields(ctx.GetLoggerFields()).Debugf("event %s.%s mapped to subscription %s",
					event, version, subscriptions[i].ID)
			}
		}
	}
	//Update Subscriptions
	r.mapAccess.Lock()
	r.qualtricsEventsToSubscriptions = relevantEvents
	r.qualtricsSubscriptionsToEvents = relevantSubscriptions
	r.mapAccess.Unlock()
	return nil
}

func (r *Reconciler) CompareState(ctx *util.RequestContext) (topicsToRegister []string, subscriptionsToDeregister []string, err error) {
	log.WithFields(ctx.GetLoggerFields()).Debug("comparing state to kyma kymaSubscriptions")

	kymaSubscriptions, err := r.EventServiceAPIClient.GetActiveSubscriptions(ctx)

	if err != nil {
		log.WithFields(ctx.GetLoggerFields()).Errorf("Error reading kymaSubscriptions from kyma: %s",
			err.Error())
		return nil, nil, fmt.Errorf("Error reading kymaSubscriptions from kyma: %s",
			err.Error())
	}

	topicsToRegister = []string{}

	kymaSubscriptionSet := make(map[string]bool)

	r.mapAccess.Lock()

	for i, _ := range kymaSubscriptions {

		//if everything is fine
		if subscriptionId, ok := r.qualtricsEventsToSubscriptions[fmt.Sprintf("%s.%s", kymaSubscriptions[i].EventType,
			kymaSubscriptions[i].EventVersion)]; ok {
			log.WithFields(ctx.GetLoggerFields()).Debugf("Event Subscription for %s.%s already exists (%s)",
				kymaSubscriptions[i].EventType, kymaSubscriptions[i].EventVersion, subscriptionId)
		} else {
			topic, err := r.TopicConverter.MapEventTypeVersionToTopic(kymaSubscriptions[i].EventType,
				kymaSubscriptions[i].EventVersion)

			if err != nil {
				log.WithFields(ctx.GetLoggerFields()).Errorf("error converting event %s and version %s "+
					"to topic: %s", kymaSubscriptions[i].EventType, kymaSubscriptions[i].EventVersion, err.Error())
			} else {
				topicsToRegister = append(topicsToRegister, topic)
			}
		}

		//index kyma subscriptions to save time later
		kymaSubscriptionSet[fmt.Sprintf("%s.%s", kymaSubscriptions[i].EventType,
			kymaSubscriptions[i].EventVersion)] = true

	}

	subscriptionsToDeregister = []string{}

	for kymaEvent, qualtricsSubscription := range r.qualtricsEventsToSubscriptions {
		if _, ok := kymaSubscriptionSet[kymaEvent]; ok {
			log.WithFields(ctx.GetLoggerFields()).Debugf("Kyma Event %s already active through subscription %s",
				kymaEvent, qualtricsSubscription)
		} else {
			subscriptionsToDeregister = append(subscriptionsToDeregister, qualtricsSubscription)
			log.WithFields(ctx.GetLoggerFields()).Debugf("Found outdated subscription %s (bound to Kyma Event %s)",
				qualtricsSubscription, kymaEvent)
		}
	}
	r.mapAccess.Unlock()

	return topicsToRegister, subscriptionsToDeregister, nil
}

func (r *Reconciler) ReconcileState(topicsToRegister []string, subscriptionsToDeregister []string,
	ctx *util.RequestContext) error {

	log.WithFields(ctx.GetLoggerFields()).Debug("reconciling state between kyma and qualtrics")

	log.WithFields(ctx.GetLoggerFields()).Debug("creating new subscriptions")

	for i, _ := range topicsToRegister {

		qualtricsSubscription := &apiclient.QualtricsSubscription{
			Topics:         topicsToRegister[i],
			PublicationURL: r.SubscriptionURL,
			SharedKey:      r.sharedKey,
		}
		subscriptionId, err := r.QualtricsAPIClient.CreateSubscription(qualtricsSubscription, ctx)

		if err != nil {
			log.WithFields(ctx.GetLoggerFields()).Errorf("error creating subscription for topic %q: %s",
				topicsToRegister[i], err.Error())
			return fmt.Errorf("error creating subscription for topic %q: %s",
				topicsToRegister[i], err.Error())
		}

		eventType, eventVersion, err := r.TopicConverter.MapTopicToEventTypeVersion(topicsToRegister[i])
		if err != nil {
			log.WithFields(ctx.GetLoggerFields()).Errorf("topic %q can't be converted to event type: %s",
				topicsToRegister[i], err.Error())
			return fmt.Errorf("topic %q can't be converted to event type: %s",
				topicsToRegister[i], err.Error())
		}
		//clean qualtrics state
		r.mapAccess.Lock()
		r.qualtricsEventsToSubscriptions[fmt.Sprintf("%s.%s", eventType, eventVersion)] = subscriptionId
		r.qualtricsSubscriptionsToEvents[subscriptionId] = fmt.Sprintf("%s.%s", eventType, eventVersion)
		r.mapAccess.Unlock()

		log.WithFields(ctx.GetLoggerFields()).Debugf("subscription for topic %q created", topicsToRegister[i])
	}

	for i, _ := range subscriptionsToDeregister {

		err := r.QualtricsAPIClient.DeleteSubscription(subscriptionsToDeregister[i], ctx)

		if err != nil {
			log.WithFields(ctx.GetLoggerFields()).Errorf("error deleting subscription for topic %q: %s",
				subscriptionsToDeregister[i], err.Error())
			return fmt.Errorf("error deleting subscription for topic %q: %s",
				subscriptionsToDeregister[i], err.Error())
		}

		//clean qualtrics state
		if eventAndVersion, ok := r.qualtricsSubscriptionsToEvents[subscriptionsToDeregister[i]]; ok {
			r.mapAccess.Lock()
			delete(r.qualtricsSubscriptionsToEvents, subscriptionsToDeregister[i])
			delete(r.qualtricsEventsToSubscriptions, eventAndVersion)
			r.mapAccess.Unlock()
		}

		log.WithFields(ctx.GetLoggerFields()).Debugf("subscription for topic %q deleted", subscriptionsToDeregister[i])
	}

	return nil
}

func (r *Reconciler) Reconcile(ctx *util.RequestContext) error {

	topicsToRegister, subscriptionsToDeregister, err := r.CompareState(ctx)

	if err != nil {
		return err
	}

	return r.ReconcileState(topicsToRegister, subscriptionsToDeregister, ctx)

}
