package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/edreed/azure-permissions-checker/pkg/apis/permissions/v1alpha1"
	azpermclient "github.com/edreed/azure-permissions-checker/pkg/client/clientset/versioned"
	"github.com/google/uuid"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	_ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	var (
		homeDir    string = homedir.HomeDir()
		kubeconfig string
		timeout    time.Duration
	)

	defaultKubeConfig := os.Getenv("KUBECONFIG")
	if len(defaultKubeConfig) == 0 && len(homeDir) != 0 {
		defaultKubeConfig = filepath.Join(homeDir, ".kube", "config")
	}

	if len(defaultKubeConfig) != 0 {
		flag.StringVar(&kubeconfig, "kubeconfig", defaultKubeConfig, "(optional) absolutepath to the kubeconfig file")
	} else {
		flag.StringVar(&kubeconfig, "kubeconfig", "", "absolute path to the kubeconfig file")
	}

	flag.DurationVar(&timeout, "timeout", 2*time.Minute, "(optional) permissions check timeout")

	flag.Parse()

	resourcePath := flag.Arg(0)
	if len(resourcePath) == 0 {
		_, _ = fmt.Fprintln(os.Stderr, "USAGE: azcheckperms [options] RESOURCE_GROUP|RESOURCE")
		flag.PrintDefaults()
		os.Exit(1)
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "ERROR: failed to get kubeconfig", err)
		os.Exit(1)
	}

	clientset, err := azpermclient.NewForConfig(config)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "ERROR: failed to create kubeclient", err)
		os.Exit(1)
	}

	azperms := &v1alpha1.AzPermission{
		ObjectMeta: v1.ObjectMeta{
			Name: fmt.Sprintf("permcheck-%s", uuid.New()),
		},
		Spec: v1alpha1.AzPermissionSpec{
			ResourcePath: resourcePath,
		},
	}

	ctx, cancel := context.WithTimeout(context.TODO(), timeout)
	defer cancel()

	azperms, err = clientset.PermissionsV1alpha1().AzPermissions().Create(ctx, azperms, v1.CreateOptions{})
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "ERROR: failed to create permissions request:", err)
		os.Exit(1)
	}
	defer func() {
		_ = clientset.PermissionsV1alpha1().AzPermissions().Delete(ctx, azperms.Name, v1.DeleteOptions{})
	}()

	err = wait.PollImmediateInfiniteWithContext(ctx, 1*time.Second, func(c context.Context) (done bool, err error) {
		azperms, err = clientset.PermissionsV1alpha1().AzPermissions().Get(c, azperms.Name, v1.GetOptions{})
		if err != nil {
			return false, err
		}

		return azperms.Status != nil && len(string(azperms.Status.State)) != 0 && azperms.Status.State != v1alpha1.Pending, nil
	})
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "ERROR: failed to get permissions:", err)
		os.Exit(1)
	}

	if azperms.Status.State == v1alpha1.Failed {
		fmt.Fprintln(os.Stderr, "ERROR: failed to get permissions:", azperms.Status.StateMessage)
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println("The principal", azperms.Status.Principal, "has the following permissions to", resourcePath, ":")
	fmt.Println()
	fmt.Println("ALLOWED:")
	for _, allowed := range azperms.Status.Allowed {
		fmt.Println("  ", allowed)
	}
	fmt.Println()
	fmt.Println("DENIED:")
	for _, denied := range azperms.Status.Denied {
		fmt.Println("  ", denied)
	}
	fmt.Println()
}
