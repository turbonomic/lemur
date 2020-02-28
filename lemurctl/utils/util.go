package utils

import (
	"fmt"
	"net"
	"os"
	"os/user"
	"strings"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func GetLocalIP() string {
	host, _ := os.Hostname()
	addrs, _ := net.LookupIP(host)
	for _, addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil {
			return ipv4.String()
		}
	}
	// Return the loop-back address
	return "127.0.0.1"
}

func GetHostName(hostip string) (string, error) {
	names, err := net.LookupAddr(hostip)
	if err == nil {
		return strings.Split(names[0], ".")[0], nil
	}
	return hostip, err
}

func GetUser(username string) (*user.User, error) {
	if username == "" {
		return user.Current()
	}
	return user.Lookup(username)
}

func GetKubeClient(kubeconfig string) (*kubernetes.Clientset, error) {
	if kubeconfig == "" {
		return nil, fmt.Errorf("no kubeconfig is specified")
	}
	kubeConfig, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}
	kubeClient, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		return nil, err
	}
	return kubeClient, nil
}
