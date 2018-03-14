/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Note: the example only works with the code within the same release/branch.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	appsv1beta1 "k8s.io/api/apps/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {

	// Flag "--kubeconfig <kubeconfig-file>"
	//
	// Specify another kubeconfig file than "~/.kube/config"
	//
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}

	// Flag "-f <yaml-file>"
	//
	// Specify another YAML configuration file than "./configuration.yaml"
	//
	var yamlfilename *string
	yamlfilename = flag.String("f", "configuration.yaml", "(optional) absolute path to the YAML configuration file")

	// Parse commandline parameters
	flag.Parse()

	// Get filename
	yamlFilepath, err := filepath.Abs(*yamlfilename)

	// Get reader from file opening
	reader, err := os.Open(yamlFilepath)
	if err != nil {
		panic(err)
	}

	// Get size of file
	fileInfo, err := reader.Stat()
	if err != nil {
		panic(err)
	}
	fileSize := fileInfo.Size()


	// Split YAML into chunks or k8s resources, respectively
	yamlDecoder := yaml.NewDocumentDecoder(ioutil.NopCloser(reader))

	// Create decoding function used for YAML to JSON decoding
	decode := scheme.Codecs.UniversalDeserializer().Decode

	// Read first resource
	yamlDeployment := make([]byte, fileSize)
	_, err = yamlDecoder.Read(yamlDeployment)
	if err != nil {
		panic(err)
	}

	// Trim unnecessary trailing 0x0 signs which are not accepted
	trimmedYamlDeployment := strings.TrimRight(string(yamlDeployment), string(byte(0)))

	// Decode deployment resource from YAML to JSON
	jsonDeployment, groupVersionKind, err := decode([]byte(trimmedYamlDeployment), nil, nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("groupVersionKind.Group: %s, groupVersionKind.Kind: %s, groupVersionKind.Version: %s\n", groupVersionKind.Group, groupVersionKind.Kind, groupVersionKind.Version)

	// Marshall JSON deployment
	d, err := json.Marshal(&jsonDeployment)
	if err != nil {
		panic(err)
	}

	// Unmarshall JSON into deployment struct
	var deployment appsv1beta1.Deployment
	err = json.Unmarshal(d, &deployment)
	if err != nil {
		panic(err)
	}

	// Read second resource
	yamlService := make([]byte, fileSize)
	_, err = yamlDecoder.Read(yamlService)
	if err != nil {
		panic(err)
	}

	// Trim unnecessary trailing 0x0 signs which are not accepted
	trimmedYamlService := strings.TrimRight(string(yamlService), string(byte(0)))

	// Decode service resource from YAML to JSON
	jsonService, _, err := decode([]byte(trimmedYamlService), nil, nil)
	if err != nil {
		panic(err)
	}

	// Marshall JSON Service
	s, err := json.Marshal(&jsonService)
	if err != nil {
		panic(err)
	}

	//Unmarshall JSON into service struct
	var service corev1.Service
	err = json.Unmarshal(s, &service)
	if err != nil {
		panic(err)
	}

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}

	// Create clientset from outside of the cluster
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	// Create client for deployments
	deploymentsClient := clientset.AppsV1beta1().Deployments(corev1.NamespaceDefault)

	// Create client for services
	servicesClient := clientset.CoreV1().Services(corev1.NamespaceDefault)

	// create deployment - try update on error
	fmt.Printf("Create deployment %q\n", deployment.Name)
	createdDeployment, err := deploymentsClient.Create(&deployment)
	if err != nil {
		fmt.Printf("Info: %s\n\n", err)

		// update deployment
		fmt.Printf("Update deployment %q\n", deployment.Name)
		updatedDeployment, err := deploymentsClient.Update(&deployment)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Deployment %q updated\n\n", updatedDeployment.Name)
	} else {
		fmt.Printf("Deployment %q created\n\n", createdDeployment.Name)
	}

	// create service - try update on error
	fmt.Printf("Create service %q\n", service.Name)
	createdService, err := servicesClient.Create(&service)
	if err != nil {
		fmt.Printf("Info: %s\n\n", err)

		// update service
		fmt.Printf("Update service %q\n", service.Name)
		updatedService, err := deploymentsClient.Update(&deployment)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Service %q updated\n", updatedService.Name)
	} else {
		fmt.Printf("Service %q created\n", createdService.Name)
	}

	// Get the port number by service's name
	runningService, err := servicesClient.Get(service.Name, metav1.GetOptions{})
	if err != nil {
		panic(err)
	}

	// Get Pod "kube-addon-manager-minikube" of "kube-system" to retrieve 'minikube ip'
	//
	// Only valid for a minikube environment
	//
	pod, err := clientset.CoreV1().Pods("kube-system").Get("kube-addon-manager-minikube", metav1.GetOptions{})
	if err != nil {
		fmt.Printf("\nInfo: %s", err)
		fmt.Printf("\nPlease help yourself and view: http://<ip-address>:%v\n\n", runningService.Spec.Ports[0].NodePort)
	} else {
		fmt.Printf("\nPlease view: http://%s:%v\n\n", pod.Status.HostIP, runningService.Spec.Ports[0].NodePort)
	}
}
