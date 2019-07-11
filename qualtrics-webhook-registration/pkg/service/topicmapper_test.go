package service

import "testing"

func TestNew(t *testing.T) {

	_, err := NewTopicmapper("../../testdata/topic-config.json")

	if err != nil {
		t.Errorf("Reading valid config failed: %s", err.Error())
	}

}

func TestMapper_MapEventTypeVersionToTopic(t *testing.T) {
	mapper, err := NewTopicmapper("../../testdata/topic-config.json")

	if err != nil {
		t.Errorf("Reading valid config failed: %s", err.Error())
	}

	//success case

	topic, err := mapper.MapEventTypeVersionToTopic("controlpanel.activateSurvey", "v1")

	if err != nil {
		t.Errorf("mapping valid event failed: %s", err.Error())
	}

	if topic != "controlpanel.activateSurvey" {
		t.Errorf("topic should be \"controlpanel.activateSurvey\" but was %q", topic)
	}

	// error case
	topic, err = mapper.MapEventTypeVersionToTopic("controlpanel.activateSurvey", "v2")

	if err == nil {
		t.Error("mapping invalid event type did not fail")
	}
}

func TestMapper_MapTopicToEventTypeVersion(t *testing.T) {
	mapper, err := NewTopicmapper("../../testdata/topic-config.json")

	if err != nil {
		t.Errorf("Reading valid config failed: %s", err.Error())
	}

	//success case

	event, eventVersion, err := mapper.MapTopicToEventTypeVersion("controlpanel.activateSurvey")

	if err != nil {
		t.Errorf("mapping valid topic failed: %s", err.Error())
	}

	if event != "controlpanel.activateSurvey" {
		t.Errorf("event should be \"controlpanel.activateSurvey\" but was %q", event)
	}

	if eventVersion != "v1" {
		t.Errorf("eventVersion should be \"v1\" but was %q", eventVersion)
	}

	//error case

	event, eventVersion, err = mapper.MapTopicToEventTypeVersion("controlpanel.activateSurvey.123")

	if err == nil {
		t.Error("mapping invalid topic did not fail")
	}

}

