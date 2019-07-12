package topicmapper

import "testing"

func TestNew(t *testing.T) {

	_, err := New("../../testing/topic_config_valid.json")

	if err != nil {
		t.Errorf("Reading valid config failed: %s", err.Error())
	}

	_, err = New("../../testing/topic_config_invalid.json")

	if err == nil {
		t.Errorf("Reading invalid config failed")
	}

}

func TestMapTopic(t *testing.T) {

	mapper, _ := New("../../testing/topic_config_valid.json")

	eventName, eventVersion, err :=
		mapper.MapTopic("sapdevelopment.surveyengine.completedResponse.SV_3k3FeLtnAHsw0VD")

	if err != nil || eventName != "surveyengine.completedResponse" || eventVersion != "v1" {
		t.Error("Topic not correctly mapped")
	}

	eventName, eventVersion, err =
		mapper.MapTopic("sapdevelopment.controlpanel.deactivateSurvey")

	if err != nil || eventName != "controlpanel.deactivateSurvey" || eventVersion != "v1" {
		t.Error("Topic not correctly mapped")
	}

	eventName, eventVersion, err =
		mapper.MapTopic("sapdevelopment.surveyengine.completedResponse.SV_3k3FeLtnAHsw0VD.testtemp")

	if err == nil {
		t.Error("Topic not correctly mapped")
	}

	eventName, eventVersion, err =
		mapper.MapTopic("sapdevelopment.threesixty.person.statusChanged")

	if err != nil || eventName != "threesixty.person.statusChanged" || eventVersion != "v1" {
		t.Error("Topic not correctly mapped")
	}

}
