package event

type TopicMapper interface {
	MapTopic(qualtricsTopicName string) (eventName string, eventVersion string, err error)
}
