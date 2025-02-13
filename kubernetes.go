package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type kubernetesCluster struct {
	kubernetes *kubernetes.Clientset
}

func getConnection() (*kubernetes.Clientset, error) {
	args := os.Args
	var kubeConfigPath string = defaultKubeConfigFilePath()
	if len(args) > 1 {
		kubeConfigPath = args[1]
	}
	// To add a minimim spinner time
	sleep := make(chan string)
	go func(c chan string) {
		time.Sleep(1000 * time.Millisecond)
		close(c)
	}(sleep)

	kubeConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		fmt.Printf("Error getting kubernetes config: %v\n", err)
		return nil, err
	}
	<-sleep
	k8s, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		return nil, err
	}
	res := k8s.RESTClient().Get().AbsPath("/healthz").Do(context.TODO())
	if res.Error() != nil {
		return nil, errors.New(res.Error().Error())
	}
	return k8s, nil
}

func defaultKubeConfigFilePath() string {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic("error getting user home dir: %v\n")
	}
	return filepath.Join(userHomeDir, ".kube", "config")
}
