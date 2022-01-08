package controller

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"regexp"

	"github.com/edreed/azure-permissions-checker/pkg/apis/permissions/v1alpha1"

	azauth "github.com/Azure/azure-sdk-for-go/services/authorization/mgmt/2015-07-01/authorization"
	azauthapi "github.com/Azure/azure-sdk-for-go/services/authorization/mgmt/2015-07-01/authorization/authorizationapi"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/retry"
	controllerruntime "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type Reconciler struct {
	k8sClient         client.Client
	permissionsClient azauthapi.PermissionsClientAPI
	clientID          string
}

const (
	resourcePathRegExString      = `(?i:^/subscriptions/(?P<subscription>[[:xdigit:]]{8}-[[:xdigit:]]{4}-[[:xdigit:]]{4}-[[:xdigit:]]{4}-[[:xdigit:]]{12})/resourcegroups/(?P<group>[-\w\._\(\)]+)/providers/(?P<provider>[\w\.]+)(?P<parentPath>/.*?)?/(?P<resourceType>[\w]+)/(?P<resourceName>[-\w_]+)$)`
	resourceGroupPathRegExString = `(?i:^/subscriptions/(?P<subscription>[[:xdigit:]]{8}-[[:xdigit:]]{4}-[[:xdigit:]]{4}-[[:xdigit:]]{4}-[[:xdigit:]]{12})/resourcegroups/(?P<group>[-\w\._\(\)]+)$)`
)

var (
	resourcePathRegEx      = regexp.MustCompile(resourcePathRegExString)
	resourceGroupPathRegEx = regexp.MustCompile(resourceGroupPathRegExString)
)

func NewReconciler(k8sclient client.Client, permissionsClient azauthapi.PermissionsClientAPI, clientID string) *Reconciler {
	r := &Reconciler{
		k8sClient:         k8sclient,
		permissionsClient: permissionsClient,
		clientID:          clientID,
	}

	return r
}

func (r *Reconciler) Reconcile(ctx context.Context, request controllerruntime.Request) (controllerruntime.Result, error) {
	log := log.FromContext(ctx)
	log.V(1).Info("reconciling permissions request")

	var permissions v1alpha1.AzPermission
	if err := r.k8sClient.Get(ctx, request.NamespacedName, &permissions); err != nil {
		log.Error(err, "failed to get permissions request")
		return controllerruntime.Result{}, nil
	}

	log = log.WithValues("resourcePath", permissions.Spec.ResourcePath)

	if permissions.Status != nil {
		log.V(1).Info("permissions request was already handled")
		return controllerruntime.Result{}, nil
	}

	err := r.updateStatus(ctx, &permissions, func(obj client.Object) error {
		permissions := obj.(*v1alpha1.AzPermission)

		if permissions.Status != nil {
			return InvalidObjectStateErr
		}

		permissions.Status = &v1alpha1.AzPermissionStatus{State: v1alpha1.Pending}

		return nil
	})
	if err != nil {
		if errors.Is(err, InvalidObjectStateErr) {
			log.V(1).Info("permissions request was already handled")
			return controllerruntime.Result{}, nil
		}

		log.Error(err, "failed to update permissions request status to Pending")
		return controllerruntime.Result{}, nil
	}

	var listFunc func() (azauth.PermissionGetResultIterator, error)

	if submatches := resourcePathRegEx.FindStringSubmatch(permissions.Spec.ResourcePath); len(submatches) == 7 {
		resourceGroupName := submatches[2]
		providerNamespace := submatches[3]
		parentPath := submatches[4]
		resourceType := submatches[5]
		resourceName := submatches[6]

		listFunc = func() (azauth.PermissionGetResultIterator, error) {
			log.V(1).Info("listing permissions for resource object", "resourceGroupName", resourceGroupName, "provider", providerNamespace, "parentPath", parentPath, "resourceType", resourceType, "resourceName", resourceName)
			return r.permissionsClient.ListForResourceComplete(ctx, resourceGroupName, providerNamespace, parentPath, resourceType, resourceName)
		}
	} else if submatches := resourceGroupPathRegEx.FindStringSubmatch(permissions.Spec.ResourcePath); len(submatches) == 3 {
		resourceGroupName := submatches[2]

		listFunc = func() (azauth.PermissionGetResultIterator, error) {
			log.V(1).Info("listing permissions for resource group", "resourceGroupName", resourceGroupName)
			return r.permissionsClient.ListForResourceGroupComplete(ctx, resourceGroupName)
		}
	} else {
		err = NewInvalidResourcePathErr(permissions.Spec.ResourcePath)
	}

	allowed := make([]string, 0)
	denied := make([]string, 0)

	if listFunc != nil {
		var iter azauth.PermissionGetResultIterator
		iter, err = listFunc()
		if err == nil {
			for iter.NotDone() {
				result := iter.Value()
				if result.Actions != nil {
					allowed = append(allowed, *result.Actions...)
				}
				if result.NotActions != nil {
					denied = append(denied, *result.NotActions...)
				}
				err = iter.NextWithContext(ctx)
				if err != nil {
					break
				}
			}
		}
	}

	updateFunc := func(obj client.Object) error {
		permissions := obj.(*v1alpha1.AzPermission)
		if permissions.Status == nil {
			permissions.Status = &v1alpha1.AzPermissionStatus{}
		}
		permissions.Status.Principal = r.clientID
		if err == nil {
			log.V(1).Info("successfully got permissions", "resourcePath", permissions.Spec.ResourcePath)
			permissions.Status.State = v1alpha1.Ready
			permissions.Status.Allowed = allowed
			permissions.Status.Denied = denied
		} else {
			log.Error(err, "failed to get permissions", "resourcePath", permissions.Spec.ResourcePath)
			permissions.Status.State = v1alpha1.Failed
			permissions.Status.StateReason = reflect.TypeOf(err).Name()
			permissions.Status.StateMessage = err.Error()
		}
		return nil
	}

	if err = r.updateStatus(ctx, &permissions, updateFunc); err != nil {
		log.Error(err, fmt.Sprintf("failed to update permissions request status to %s", string(permissions.Status.State)))
		return controllerruntime.Result{}, nil
	}

	return controllerruntime.Result{}, nil
}

func (r *Reconciler) updateStatus(ctx context.Context, obj client.Object, updateFunc func(client.Object) error) error {
	attempt := 0
	objKey := types.NamespacedName{Name: obj.GetName(), Namespace: obj.GetNamespace()}

	return retry.RetryOnConflict(retry.DefaultBackoff, func() error {
		if attempt > 0 {
			if err := r.k8sClient.Get(ctx, objKey, obj); err != nil {
				return err
			}
		}

		attempt++

		if err := updateFunc(obj); err != nil {
			return err
		}

		return r.k8sClient.Status().Update(ctx, obj)
	})
}
