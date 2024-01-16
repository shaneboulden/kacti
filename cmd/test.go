/*
Copyright © 2023 Shane Boulden (sboulden@redhat.com)

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
package cmd

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

const (
	ResultFail = 0
	ResultPass = 1

	ExitOk  = 0
	ExitErr = 1
)

type KactiTests struct {
	Tests []Test `yaml:"kacti-tests"`
}

type Test struct {
	Name      string `yaml:"name"`
	Namespace string `yaml:"namespace"`
	Image     string `yaml:"image"`
}

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Functionally test admission control",
	Long: `Perform functional verification tests against Kubernetes admission controllers.

Tests are specififed in files that reference test names, namespaces and images. For example:

kacti-tests:
- name: pwnkit
  image: quay.io/the-worst-containers/pwnkit:v0.2
  namespace: app-deploy`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		errCode := runTest(strings.Join(args, " "))
		os.Exit(errCode)
	},
}

func runTest(file string) int {
	// set up k8s auth
	var kubeconfig *string

	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()
	fmt.Println("Setting up kubeconfig from: " + *kubeconfig)

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	// read the YAML file defining tests
	fmt.Println("Using tests from: " + file)
	data, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println(err)
		return ExitErr
	}

	var testConfig KactiTests

	err = yaml.Unmarshal(data, &testConfig)
	if err != nil {
		fmt.Println(err)
		return ExitErr
	}

	for _, test := range testConfig.Tests {
		fmt.Println(("Running test: " + test.Name + " { ns: " + test.Namespace + " / img: " + test.Image + " }"))
		// - try and deploy the workload
		i := int32(1)
		deploymentsClient := clientset.AppsV1().Deployments(test.Namespace)

		deployment := &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name: test.Name,
			},
			Spec: appsv1.DeploymentSpec{
				Replicas: &i,
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{
						"app": test.Name,
					},
				},
				Template: apiv1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							"app": test.Name,
						},
					},
					Spec: apiv1.PodSpec{
						Containers: []apiv1.Container{
							{
								Name:  "test",
								Image: test.Image,
							},
						},
					},
				},
			},
		}

		Ctx := context.Background()

		// we don't care about the message or response - we'll just test whether the files were created
		deploymentsClient.Create(Ctx, deployment, metav1.CreateOptions{})
	}

	// Currently this is just a sleep to give StackRox time to scale
	// the deployment replicas down. It's not waiting on any condition...
	time.Sleep(8 * time.Second)

	r := make(map[Test]string)
	n := make(map[Test]int32)
	for _, test := range testConfig.Tests {
		result, err := checkAdmissionControl(clientset, test.Namespace, test.Name)
		r[test] = result
		if err != nil {
			n[test] = ResultFail
		} else {
			n[test] = ResultPass
		}

		// cleanup
		deploymentsClient := clientset.AppsV1().Deployments(test.Namespace)
		deploymentsClient.Delete(context.Background(), test.Name, *&metav1.DeleteOptions{})
	}

	// print results and set return code
	fmt.Printf("Results:\n")
	var retCode = ExitOk
	for _, test := range testConfig.Tests {
		fmt.Printf(test.Name + " { ns: " + test.Namespace + " / img:" + test.Image + " }\n")
		if n[test] == ResultPass {
			fmt.Printf(" -> %s, %s\n\n", color.GreenString("Success"), r[test])
		} else {
			retCode = ExitErr
			fmt.Printf(" -> %s, %s\n\n", color.RedString("Failed"), r[test])
		}
	}
	return retCode
}

func checkAdmissionControl(clientset *kubernetes.Clientset, namespace, deploymentName string) (string, error) {
	// Get the deployment - if it doesn't exist, the admission controller has likely blocked it
	deployment, err := clientset.AppsV1().Deployments(namespace).Get(context.Background(), deploymentName, metav1.GetOptions{})
	if err != nil {
		return "Deployment creation was blocked", nil
	}

	// Check if the deployment is scaled-down (e.g., scaled to zero replicas)
	if *deployment.Spec.Replicas == 0 {
		return "Deployment scaled to zero replicas", nil
	}

	return "Deployment was created successfully and scaled up", errors.New("Admission control failed")
}

func init() {
	rootCmd.AddCommand(testCmd)
}
