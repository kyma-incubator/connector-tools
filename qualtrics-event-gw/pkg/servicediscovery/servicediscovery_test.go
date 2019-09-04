package servicediscovery

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	testclient "k8s.io/client-go/kubernetes/fake"
	"testing"
)

func TestKubernetesClient_DiscoverEventServiceURL(t *testing.T) {

	client := KubernetesClient{
		client: testclient.NewSimpleClientset(&v1.Service{
			TypeMeta: metav1.TypeMeta{
				Kind: "Service",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: "testservice",
				Namespace: "kyma-integration",
				Labels: map[string]string{
					"application": "qualtrics",
					"heritage": "Tiller-event-service",
				},
			},
		}),
	}

	url, err := client.DiscoverEventServiceURL("kyma-integration",
		"application=qualtrics, heritage=Tiller-event-service", "qualtrics")

	if err != nil {
		t.Fatalf("Event Service discovery must not fail: %s", err.Error())
	}

	targetUrl := "http://testservice.kyma-integration.svc.cluster.local:8081/qualtrics/v1/events"

	if url != targetUrl {
		t.Errorf("Returned Url should be %q, but is %q", targetUrl, url)
	}


	//Test error for non existant service
	_, err = client.DiscoverEventServiceURL("kyma-integration",
		"application=doesnotexist, heritage=Tiller-event-service", "qualtrics")

	if err == nil {
		t.Errorf("Event Service discovery must fail, but did not")
	}
}
