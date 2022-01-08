package main

import (
	"flag"
	"os"

	azauth "github.com/Azure/azure-sdk-for-go/services/authorization/mgmt/2015-07-01/authorization"
	"github.com/Azure/go-autorest/autorest"
	"github.com/edreed/azure-permissions-checker/pkg/apis/permissions/v1alpha1"
	"github.com/edreed/azure-permissions-checker/pkg/permissions/controller"
	"github.com/edreed/azure-permissions-checker/pkg/permissions/util"

	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
	"k8s.io/klog/v2/klogr"
	azproviderauth "sigs.k8s.io/cloud-provider-azure/pkg/auth"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/event"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

func init() {
	klog.InitFlags(nil)
}

func main() {
	flag.Parse()

	logf.SetLogger(klogr.New())

	log := logf.Log.WithName("azpermissions")

	kubeConfig := config.GetConfigOrDie()

	kubeClient, err := clientset.NewForConfig(kubeConfig)
	if err != nil {
		log.Error(err, "failed to create kubeclient")
		os.Exit(1)
	}

	azcloud, err := util.GetCloudProviderFromClient(kubeClient, "", "", "")
	if err != nil {
		log.Error(err, "failed to create Azure cloud provider")
		os.Exit(1)
	}

	token, err := azproviderauth.GetServicePrincipalToken(&azcloud.AzureAuthConfig, &azcloud.Environment, azcloud.Environment.ServiceManagementEndpoint)
	if err != nil {
		log.Error(err, "failed to get service identity token")
		os.Exit(1)
	}

	permissionsClient := azauth.NewPermissionsClient(azcloud.SubscriptionID)
	permissionsClient.Authorizer = autorest.NewBearerAuthorizer(token)

	mgr, err := manager.New(kubeConfig, manager.Options{Logger: log})
	if err != nil {
		log.Error(err, "failed to create controller runtime manager")
		os.Exit(1)
	}

	err = v1alpha1.AddToScheme(mgr.GetScheme())
	if err != nil {
		log.Error(err, "failed to add scheme")
		os.Exit(1)
	}

	onCreation := &predicate.Funcs{
		CreateFunc: func(ce event.CreateEvent) bool {
			return ce.Object.(*v1alpha1.AzPermission).Status == nil
		},
		UpdateFunc: func(ue event.UpdateEvent) bool {
			return false
		},
		GenericFunc: func(ge event.GenericEvent) bool {
			return false
		},
		DeleteFunc: func(de event.DeleteEvent) bool {
			return false
		},
	}

	err = builder.
		ControllerManagedBy(mgr).
		For(&v1alpha1.AzPermission{}, builder.WithPredicates(onCreation)).
		Complete(controller.NewReconciler(mgr.GetClient(), permissionsClient, azcloud.AADClientID))
	if err != nil {
		log.Error(err, "failed to create controller")
		os.Exit(1)
	}

	err = mgr.Start(signals.SetupSignalHandler())
	if err != nil {
		log.Error(err, "failed to start controller runtime manager")
		os.Exit(1)
	}
}
