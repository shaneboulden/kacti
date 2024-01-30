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
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	goyaml "gopkg.in/yaml.v3"
	k8syaml "sigs.k8s.io/yaml"

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

var runDeploy = false
var runPod = false
var runFile = false

var trialName = ""
var trialNamespace = ""
var trialImage = ""

var verbose = false

type Trial struct {
	Name      string `yaml:"name"`
	Namespace string `yaml:"namespace"`
	Image     string `yaml:"image"`
	Template  string `yaml:"template"`
}

var trialsCmd = &cobra.Command{
	Use:   "trials",
	Short: "Functionally test admission control",
	Long:  `Perform functional verification trials against Kubernetes admission controllers`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// set up kubeconfig
		var kubeconfig *string
		if os.Getenv("KUBECONFIG") != "" {
			kubeconfig = flag.String("kubeconfig", os.Getenv("KUBECONFIG"), "environment variable holding kubeconfig")
		} else if home := homedir.HomeDir(); home != "" {
			kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		} else {
			kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
		}

		flag.Parse()
		if verbose {
			fmt.Println(color.YellowString("Setting up kubeconfig from: " + *kubeconfig))
		}
		config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			panic(err)
		}
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			panic(err)
		}

		// run the trials using kubeconfig and flags provided
		var errCode = 1
		if runFile {
			errCode = runTrialsFromFile(strings.Join(args, " "), clientset)
		} else if runDeploy {
			trial := Trial{strings.Join(args, ""), trialNamespace, trialImage, ""}
			errCode = runDeployTrialStandalone(trial, clientset)
		}
		os.Exit(errCode)
	},
}

func runDeployTrialStandalone(trial Trial, clientset *kubernetes.Clientset) int {
	if verbose {
		fmt.Println(color.YellowString("Running trial: " + trial.Name + " { ns: " + trial.Namespace + " / img: " + trial.Image + " }"))
	}
	runDeploymentTrial(trial, clientset)

	// Currently this is just a sleep to give StackRox time to scale
	// the deployment replicas down. It's not waiting on any condition...
	time.Sleep(8 * time.Second)

	result, err := checkAdmissionControl(clientset, trial.Namespace, trial.Name)

	// print results and set return code
	var retCode = ExitErr
	if err != nil {
		fmt.Printf(" -> %s, %s\n\n", color.RedString("Failed"), result)
	} else {
		fmt.Printf(" -> %s, %s\n\n", color.GreenString("Success"), result)
		retCode = ExitOk

	}

	// cleanup
	deploymentsClient := clientset.AppsV1().Deployments(trial.Namespace)
	deploymentsClient.Delete(context.Background(), trial.Name, *&metav1.DeleteOptions{})

	return retCode
}

func runTrialsFromFile(file string, clientset *kubernetes.Clientset) int {
	// read the YAML file defining trials
	if verbose {
		fmt.Println(color.YellowString("Using trials from: " + file))
	}
	data, err := os.ReadFile(file)
	if err != nil {
		fmt.Println(err)
		return ExitErr
	}

	var trials []Trial

	err = goyaml.Unmarshal(data, &trials)
	if err != nil {
		fmt.Println(err)
		return ExitErr
	}

	for _, trial := range trials {
		if verbose {
			fmt.Println(color.YellowString("Running trial: " + trial.Name + " { ns: " + trial.Namespace + " / img: " + trial.Image + " }"))
		}
		runDeploymentTrial(trial, clientset)
	}

	// Currently this is just a sleep to give StackRox time to scale
	// the deployment replicas down. It's not waiting on any condition...
	time.Sleep(8 * time.Second)

	r := make(map[Trial]string)
	n := make(map[Trial]int32)
	for _, trial := range trials {
		result, err := checkAdmissionControl(clientset, trial.Namespace, trial.Name)
		r[trial] = result
		if err != nil {
			n[trial] = ResultFail
		} else {
			n[trial] = ResultPass
		}

		// cleanup
		deploymentsClient := clientset.AppsV1().Deployments(trial.Namespace)
		deploymentsClient.Delete(context.Background(), trial.Name, *&metav1.DeleteOptions{})
	}

	// print results and set return code
	fmt.Printf("Results:\n")
	var retCode = ExitOk
	for _, trial := range trials {
		fmt.Printf(trial.Name + " { ns: " + trial.Namespace + " / img:" + trial.Image + " }\n")
		if n[trial] == ResultPass {
			fmt.Printf(" -> %s, %s\n\n", color.GreenString("Success"), r[trial])
		} else {
			retCode = ExitErr
			fmt.Printf(" -> %s, %s\n\n", color.RedString("Failed"), r[trial])
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

func runDeploymentTrial(trial Trial, clientset *kubernetes.Clientset) {
	i := int32(1)
	deploymentsClient := clientset.AppsV1().Deployments(trial.Namespace)

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: trial.Name,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &i,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": trial.Name,
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": trial.Name,
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  trial.Name,
							Image: trial.Image,
						},
					},
				},
			},
		},
	}

	if trial.Template != "" {
		parseTemplate(trial.Template, trial, deployment)
	}

	Ctx := context.Background()

	// Create the deployment on the cluster
	_, err := deploymentsClient.Create(Ctx, deployment, metav1.CreateOptions{})
	if err != nil {
		fmt.Printf("Error creating Deployment: %v\n", err)
	}
}

func parseTemplate(templateRef string, trial Trial, deployment *appsv1.Deployment) {

	yamlFile, err := os.ReadFile(templateRef)
	if err != nil {
		fmt.Printf("Error reading deployment template: %v\n", err)
		os.Exit(1)
	}

	err = k8syaml.Unmarshal(yamlFile, deployment)
	if err != nil {
		fmt.Printf("Error unmarshalling deployment template: %v\n", err)
		os.Exit(1)
	}

	// set name, namespace
	deployment.ObjectMeta.Name = trial.Name
	deployment.ObjectMeta.Namespace = trial.Namespace

	// set labels
	deployment.ObjectMeta.Labels["app"] = trial.Name
	deployment.Spec.Template.ObjectMeta.Labels["app"] = trial.Name
	deployment.Spec.Template.ObjectMeta.Namespace = trial.Namespace

	// set container images
	deployment.Spec.Template.Spec.Containers[0].Image = trial.Image
}

func init() {
	// add flags for declarative and imperative trial runs
	trialsCmd.Flags().BoolVarP(&runDeploy, "deploy", "d", false, "Run a deployment trial")
	trialsCmd.Flags().BoolVarP(&runFile, "file", "f", false, "Run a set of trials from a file")

	// add flags required for `deploy` and `pod` trials
	trialsCmd.Flags().StringVarP(&trialNamespace, "namespace", "n", "", "Namespace for the trial")
	trialsCmd.Flags().StringVarP(&trialImage, "image", "i", "", "Image for the trial")

	// set flags as mutually exclusive, using the following rules:
	// --deploy -> used standalone
	// --file -> used standalone
	trialsCmd.MarkFlagsMutuallyExclusive("deploy", "file")

	// set flags required for imperative trials
	trialsCmd.MarkFlagsRequiredTogether("deploy", "namespace", "image")

	trialsCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")

	rootCmd.AddCommand(trialsCmd)
}
