package servicediscovery

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)



type KubernetesClient struct {
	client kubernetes.Interface
}


// InitOutOfCluster initializes the client running inside of a Kubernetes Cluster
func InitInCluster() (*KubernetesClient, error) {

	config, err := rest.InClusterConfig()
	if err != nil {
		log.Errorf("error connecting to kubernetes API, %s", err.Error())
		return nil, fmt.Errorf("error connecting to kubernetes API: %s", err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Errorf("error connecting to kubernetes API, %s", err.Error())
		return nil, fmt.Errorf("error connecting to kubernetes API: %s", err.Error())
	}

	return &KubernetesClient{
		client: clientset,
	}, nil
}



// InitOutOfCluster initializes the client based on a local kubeconfig file
func InitOutOfCluster(kubeconfig string) (*KubernetesClient, error) {

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Errorf("error connecting to kubernetes API, %s", err.Error())
		return nil, fmt.Errorf("error connecting to kubernetes API: %s", err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Errorf("error connecting to kubernetes API, %s", err.Error())
		return nil, fmt.Errorf("error connecting to kubernetes API: %s", err.Error())
	}


	return &KubernetesClient{
		client: clientset,
	}, nil
}

// DiscoverService generates the in cluster url to discover the event gateway to determine active subscriptions
func (k *KubernetesClient) DiscoverEventServiceURL (namespace string, labelsesector string) (kymaEventGatewayBaseURL string,
	err error) {

	serviceList, err := k.client.CoreV1().Services(namespace).List(metav1.ListOptions {
		LabelSelector: labelsesector})
	if err != nil {
		log.Errorf("error reading services in namespace %q for labelselector %q: %s",
			namespace, labelsesector, err.Error())
		return "", fmt.Errorf("error reading services in namespace %q for labelselector %q: %s",
			namespace, labelsesector, err.Error())
	}


	//warn if more than one Service was discovered
	if len(serviceList.Items) > 1 {
		log.Warnf("more than one service discovered in namespace %q for labelselector %q",
			namespace, labelsesector)
	}

	//take the first service matching label selector
	for i := range serviceList.Items {
		return fmt.Sprintf("http://%s.%s.svc.cluster.local:8081", serviceList.Items[i].Name, namespace),
		nil
	}


	//this is only reached if there was no service discovered, hence error out
	log.Errorf("no service discovered in namespace %q for labelselector %q",
		namespace, labelsesector)

	return "",
		fmt.Errorf("no service discovered in namespace %q for labelselector %q",
			namespace, labelsesector)


}


