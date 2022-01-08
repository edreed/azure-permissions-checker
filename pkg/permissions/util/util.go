package util

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/edreed/azure-permissions-checker/pkg/permissions/constants"

	"k8s.io/apimachinery/pkg/api/errors"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
	"sigs.k8s.io/cloud-provider-azure/pkg/azureclients"
	azure "sigs.k8s.io/cloud-provider-azure/pkg/provider"
)

// GetCloudProviderFromClient get Azure Cloud Provider
func GetCloudProviderFromClient(kubeClient *clientset.Clientset, secretName string, secretNamespace string, userAgent string) (*azure.Cloud, error) {
	var config *azure.Config
	var fromSecret bool

	// Try to get the configuration from a K8s secret first...
	if kubeClient != nil {
		az := &azure.Cloud{
			InitSecretConfig: azure.InitSecretConfig{
				SecretName:      secretName,
				SecretNamespace: secretNamespace,
				CloudConfigKey:  "cloud-config",
			},
		}

		az.KubeClient = kubeClient

		var err error

		config, err = az.GetConfigFromSecret()
		if err == nil {
			fromSecret = true
		} else {
			if !errors.IsNotFound(err) {
				klog.Warningf("failed to create cloud config from secret %s/%s: %v", az.SecretNamespace, az.SecretName, err)
			}
		}
	}

	// ... and fallback to reading configuration file on disk.
	if config == nil {
		credFile, ok := os.LookupEnv(constants.DefaultAzureCredentialFileEnv)
		if ok && strings.TrimSpace(credFile) != "" {
			klog.V(2).Infof("%s env var set as %v", constants.DefaultAzureCredentialFileEnv, credFile)
		} else {
			if IsWindowsOS() {
				credFile = constants.DefaultCredFilePathWindows
			} else {
				credFile = constants.DefaultCredFilePathLinux
			}
			klog.V(2).Infof("use default %s env var: %v", constants.DefaultAzureCredentialFileEnv, credFile)
		}

		credFileConfig, err := os.Open(credFile)
		if err != nil {
			err = fmt.Errorf("failed to load cloud config from file %q: %v", credFile, err)
			klog.Errorf(err.Error())
			return nil, err
		}
		defer credFileConfig.Close()

		config, err = azure.ParseConfig(credFileConfig)
		if err != nil {
			err = fmt.Errorf("failed to parse cloud config file %q: %v", credFile, err)
			klog.Errorf(err.Error())
			return nil, err
		}
	}

	// Override configuration values
	config.DiskRateLimit = &azureclients.RateLimitConfig{
		CloudProviderRateLimit: false,
	}
	config.SnapshotRateLimit = &azureclients.RateLimitConfig{
		CloudProviderRateLimit: false,
	}
	config.UserAgent = userAgent

	// Create a new cloud provider
	az, err := azure.NewCloudWithoutFeatureGatesFromConfig(config, fromSecret, false)
	if err != nil {
		err = fmt.Errorf("failed to create cloud: %v", err)
		klog.Errorf(err.Error())
		return nil, err
	}

	// reassign kubeClient
	if kubeClient != nil && az.KubeClient == nil {
		az.KubeClient = kubeClient
	}

	return az, nil
}

// GetCloudProvider get Azure Cloud Provider
func GetCloudProvider(kubeConfig, secretName, secretNamespace, userAgent string) (*azure.Cloud, error) {
	kubeClient, err := GetKubeClient(kubeConfig)
	if err != nil {
		klog.Warningf("get kubeconfig(%s) failed with error: %v", kubeConfig, err)
		if !os.IsNotExist(err) && err != rest.ErrNotInCluster {
			return nil, fmt.Errorf("failed to get KubeClient: %v", err)
		}
	}
	return GetCloudProviderFromClient(kubeClient, secretName, secretNamespace, userAgent)
}

// GetKubeConfig gets config object from config file
func GetKubeConfig(kubeconfig string) (config *rest.Config, err error) {
	if kubeconfig != "" {
		if config, err = clientcmd.BuildConfigFromFlags("", kubeconfig); err != nil {
			return nil, err
		}
	} else {
		if config, err = rest.InClusterConfig(); err != nil {
			return nil, err
		}
	}
	return config, err
}

func GetKubeClient(kubeconfig string) (*clientset.Clientset, error) {
	config, err := GetKubeConfig(kubeconfig)
	if err != nil {
		return nil, err
	}

	return clientset.NewForConfig(config)
}

func GetKubeClientOrDie(kubeconfig string) *clientset.Clientset {
	kubeClient, err := GetKubeClient(kubeconfig)
	if err != nil {
		klog.Error("failed to create kubeclient")
		os.Exit(10)
	}

	return kubeClient
}

func IsWindowsOS() bool {
	return runtime.GOOS == "windows"
}
