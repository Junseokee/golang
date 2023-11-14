package main

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"os"
	"path/filepath"
)

func main() {
	var (
		client           *kubernetes.Clientset
		deploymentLabels map[string]string
		err              error
	)
	ctx := context.Background()
	if client, err = getClient(); err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(1)
	}
	if deploymentLabels, err = getService(ctx, client); err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(1)
	}
	fmt.Printf("deploy finished. Did a deploy with labels: %+v \n", deploymentLabels)
}
func getClient() (*kubernetes.Clientset, error) {
	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(homedir.HomeDir(), ".kube", "config"))
	if err != nil {
		return nil, err
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return clientset, nil
}

//func deploy(ctx context.Context, client *kubernetes.Clientset) (map[string]string, error) {
//	var deployment *v1.Deployment
//
//	appFile, err := ioutil.ReadFile("app.yaml")
//	if err != nil {
//		return nil, fmt.Errorf("readfile error: %s", err)
//	}
//	obj, groupVersionKind, err := scheme.Codecs.UniversalDeserializer().Decode(appFile, nil, nil)
//
//	switch obj.(type) {
//	case *v1.Deployment:
//		deployment = obj.(*v1.Deployment)
//	default:
//		return nil, fmt.Errorf("Unrecognized type: %s\n", groupVersionKind)
//	}
//
//	deploymentResponse, err := client.AppsV1().Deployments("default").Create(ctx, deployment, metav1.CreateOptions{}) //namespace = default
//	if err != nil {
//		return nil, fmt.Errorf("deployment error: %s", err)
//	}
//	return deploymentResponse.Spec.Template.Labels, nil
//
//}

func getService(ctx context.Context, client *kubernetes.Clientset) (map[string]string, error) {
	service, err := client.CoreV1().Endpoints("monitoring").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("Service not found: %s", err)
	}
	for _, s := range service.Items {
		fmt.Printf("Service Name: %s\n", s.Name)
	}
	return nil, fmt.Errorf("Service not found: %s", err)
}
