package controller

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	stdlog "log"
	"reflect"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/authorization/mgmt/2015-07-01/authorization"
	"github.com/edreed/azure-permissions-checker/pkg/apis/permissions"
	"github.com/edreed/azure-permissions-checker/pkg/apis/permissions/v1alpha1"
	"github.com/edreed/azure-permissions-checker/pkg/client/clientset/versioned/fake"
	"github.com/edreed/azure-permissions-checker/pkg/permissions/controller/mockauthorizationapi"
	"github.com/edreed/azure-permissions-checker/pkg/permissions/controller/mockclient"
	"github.com/go-logr/logr"
	"github.com/go-logr/stdr"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	controllerruntime "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	testSubscription  = "01234567-89ab-cdef-fedc-ba9876543210"
	testResourceGroup = "test-rg"
	testProvider      = "Microsoft.Compute"
	testResourceType  = "disks"
	testResourceName  = "test-disk-1"
)

var (
	testResourceGroupPathName   = "test-resource-group-permissions"
	testResourceGroupPath       = fmt.Sprintf("/subscriptions/%s/resourceGroups/%s", testSubscription, testResourceGroup)
	testResourcePathName        = "test-resource-permissions"
	testResourcePath            = fmt.Sprintf("%s/providers/%s/%s/%s", testResourceGroupPath, testProvider, testResourceType, testResourceName)
	testInvalidResourcePathName = "test-invalid-resource-path"
	testInvalidResourcePath     = "invalid path"

	testResourceGroupPermissions = v1alpha1.AzPermission{
		ObjectMeta: v1.ObjectMeta{
			Name: testResourceGroupPathName,
		},
		Spec: v1alpha1.AzPermissionSpec{
			ResourcePath: testResourceGroupPath,
		},
	}

	testResourcePermissions = v1alpha1.AzPermission{
		ObjectMeta: v1.ObjectMeta{
			Name: testResourcePathName,
		},
		Spec: v1alpha1.AzPermissionSpec{
			ResourcePath: testResourcePath,
		},
	}

	testResourcePendingPermissions = v1alpha1.AzPermission{
		ObjectMeta: v1.ObjectMeta{
			Name: testResourcePathName,
		},
		Spec: v1alpha1.AzPermissionSpec{
			ResourcePath: testResourcePath,
		},
		Status: &v1alpha1.AzPermissionStatus{
			State: v1alpha1.Pending,
		},
	}

	testInvalidResourcePermissions = v1alpha1.AzPermission{
		ObjectMeta: v1.ObjectMeta{
			Name: testInvalidResourcePathName,
		},
		Spec: v1alpha1.AzPermissionSpec{
			ResourcePath: testInvalidResourcePath,
		},
	}

	testResourceGroupPermissionsActions    = []string{"Microsoft.Compute/disks/create", "Microsoft.Compute/disks/read", "Microsoft.Compute/disks/write"}
	testResourceGroupPermissionsNotActions = []string{"Microsoft.Compute/disks/delete"}

	testResourceGroupPermissionsResult = []authorization.Permission{
		{
			Actions:    &testResourceGroupPermissionsActions,
			NotActions: &testResourceGroupPermissionsNotActions,
		},
	}

	testResourceGroupPermissionGetResult = authorization.PermissionGetResult{
		Value: &testResourceGroupPermissionsResult,
	}

	testResourcePermissionsActions    = []string{"Microsoft.Compute/disks/read", "Microsoft.Compute/disks/write"}
	testResourcePermissionsNotActions = []string{}

	testResourcePermissionsResult = []authorization.Permission{
		{
			Actions:    &testResourcePermissionsActions,
			NotActions: &testResourcePermissionsNotActions,
		},
	}

	testResourcePermissionGetResult = authorization.PermissionGetResult{
		Value: &testResourcePermissionsResult,
	}
)

func newTestReconciler(t *testing.T, mockCtrl *gomock.Controller, fakeClientSet *fake.Clientset, updateStatusErr error, updateStatusSuccessCount int) *Reconciler {
	mockStatus := mockclient.NewMockStatusWriter(mockCtrl)
	mockStatus.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, obj client.Object, opts ...client.UpdateOption) error {
			targetObj, ok := obj.(*v1alpha1.AzPermission)
			if !ok {
				gvk := obj.GetObjectKind().GroupVersionKind()
				return apierrors.NewNotFound(schema.GroupResource{Group: gvk.Group, Resource: gvk.Kind}, obj.GetName())
			}

			updatedObj, err := fakeClientSet.PermissionsV1alpha1().AzPermissions().UpdateStatus(ctx, targetObj, v1.UpdateOptions{})
			if err != nil {
				return err
			}

			if updateStatusSuccessCount == 0 && updateStatusErr != nil {
				return updateStatusErr
			}

			if updateStatusSuccessCount > 0 {
				updateStatusSuccessCount--
			}

			updatedObj.DeepCopyInto(targetObj)

			return nil
		}).
		AnyTimes()

	mockClient := mockclient.NewMockClient(mockCtrl)
	mockClient.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, key types.NamespacedName, obj client.Object) error {
			targetObj, ok := obj.(*v1alpha1.AzPermission)
			if !ok {
				gvk := obj.GetObjectKind().GroupVersionKind()
				return apierrors.NewNotFound(schema.GroupResource{Group: gvk.Group, Resource: gvk.Kind}, key.String())
			}

			permissions, err := fakeClientSet.PermissionsV1alpha1().AzPermissions().Get(ctx, key.Name, v1.GetOptions{})
			if err != nil {
				return err
			}

			permissions.DeepCopyInto(targetObj)

			return nil
		}).
		AnyTimes()
	mockClient.EXPECT().Status().Return(mockStatus).AnyTimes()

	mockPermissionsClient := mockauthorizationapi.NewMockPermissionsClientAPI(mockCtrl)
	mockPermissionsClient.EXPECT().ListForResourceGroupComplete(gomock.Any(), testResourceGroup).
		DoAndReturn(func(ctx context.Context, resourceGroupName string) (authorization.PermissionGetResultIterator, error) {
			page := authorization.NewPermissionGetResultPage(
				testResourceGroupPermissionGetResult,
				func(c context.Context, pgr authorization.PermissionGetResult) (authorization.PermissionGetResult, error) {
					return authorization.PermissionGetResult{}, nil
				})

			return authorization.NewPermissionGetResultIterator(page), nil
		}).
		AnyTimes()
	mockPermissionsClient.EXPECT().ListForResourceComplete(gomock.Any(), testResourceGroup, testProvider, "", testResourceType, testResourceName).
		DoAndReturn(func(ctx context.Context, resourceGroupName, providerNamespace, providerPath, resourceType, resourceName string) (authorization.PermissionGetResultIterator, error) {
			page := authorization.NewPermissionGetResultPage(
				testResourcePermissionGetResult,
				func(c context.Context, pgr authorization.PermissionGetResult) (authorization.PermissionGetResult, error) {
					return authorization.PermissionGetResult{}, nil
				})

			return authorization.NewPermissionGetResultIterator(page), nil
		}).
		AnyTimes()

	return NewReconciler(mockClient, mockPermissionsClient, "test-client")
}

func TestReconcilerReconcile(t *testing.T) {
	invalidPathErr := NewInvalidResourcePathErr(testInvalidResourcePath)

	tests := []struct {
		description              string
		objects                  []runtime.Object
		request                  controllerruntime.Request
		updateStatusSuccessCount int
		updateStatusErr          error
		expectedLog              string
		skipVerifyObject         bool
		expectedState            v1alpha1.AzPermissionState
		expectedAllowed          []string
		expectedDenied           []string
		expectedStateReason      string
		expectedStateMessage     string
	}{
		{
			description:     "[Success] Resource Group Path",
			objects:         []runtime.Object{&testResourceGroupPermissions},
			request:         controllerruntime.Request{NamespacedName: types.NamespacedName{Name: testResourceGroupPathName}},
			expectedLog:     "listing permissions for resource group",
			expectedState:   v1alpha1.Ready,
			expectedAllowed: testResourceGroupPermissionsActions,
			expectedDenied:  testResourceGroupPermissionsNotActions,
		},
		{
			description:     "[Success] Resource Path",
			objects:         []runtime.Object{&testResourcePermissions},
			request:         controllerruntime.Request{NamespacedName: types.NamespacedName{Name: testResourcePathName}},
			expectedLog:     "listing permissions for resource object",
			expectedState:   v1alpha1.Ready,
			expectedAllowed: testResourcePermissionsActions,
			expectedDenied:  testResourcePermissionsNotActions,
		},
		{
			description:      "[Failure] AzPermission Not Found",
			objects:          []runtime.Object{},
			request:          controllerruntime.Request{NamespacedName: types.NamespacedName{Name: testResourcePathName}},
			expectedLog:      "failed to get permissions request",
			skipVerifyObject: true,
		},
		{
			description:      "[Failure] AzPermission Pending",
			objects:          []runtime.Object{&testResourcePendingPermissions},
			request:          controllerruntime.Request{NamespacedName: types.NamespacedName{Name: testResourcePathName}},
			expectedLog:      "permissions request was already handled",
			skipVerifyObject: true,
		},
		{
			description:      "[Failure] AzPermission UpdateStatus Conflict",
			objects:          []runtime.Object{&testResourcePermissions},
			request:          controllerruntime.Request{NamespacedName: types.NamespacedName{Name: testResourcePathName}},
			updateStatusErr:  apierrors.NewConflict(schema.GroupResource{Group: permissions.GroupName, Resource: "azpermissions"}, testResourcePathName, errors.New("conflict")),
			expectedLog:      "permissions request was already handled",
			skipVerifyObject: true,
		},
		{
			description:      "[Failure] AzPermission First UpdateStatus Failed",
			objects:          []runtime.Object{&testResourcePermissions},
			request:          controllerruntime.Request{NamespacedName: types.NamespacedName{Name: testResourcePathName}},
			updateStatusErr:  errors.New("failed"),
			expectedLog:      "failed to update permissions request status to Pending",
			skipVerifyObject: true,
		},
		{
			description:              "[Failure] AzPermission Second UpdateStatus Failed",
			objects:                  []runtime.Object{&testResourcePermissions},
			request:                  controllerruntime.Request{NamespacedName: types.NamespacedName{Name: testResourcePathName}},
			updateStatusSuccessCount: 1,
			updateStatusErr:          errors.New("failed"),
			expectedLog:              "failed to update permissions request status to Ready",
			skipVerifyObject:         true,
		},
		{
			description:          "[Failure] Invalid Resource Path",
			objects:              []runtime.Object{&testInvalidResourcePermissions},
			request:              controllerruntime.Request{NamespacedName: types.NamespacedName{Name: testInvalidResourcePathName}},
			expectedState:        v1alpha1.Failed,
			expectedStateReason:  reflect.TypeOf(invalidPathErr).Name(),
			expectedStateMessage: invalidPathErr.Error(),
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.description, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			fakeClientSet := fake.NewSimpleClientset(test.objects...)

			reconciler := newTestReconciler(t, mockCtrl, fakeClientSet, test.updateStatusErr, test.updateStatusSuccessCount)

			logBuffer := &bytes.Buffer{}
			log := stdr.NewWithOptions(stdlog.New(logBuffer, "", stdlog.LstdFlags), stdr.Options{LogCaller: stdr.All})
			oldVerbosity := stdr.SetVerbosity(1)
			defer stdr.SetVerbosity(oldVerbosity)

			ctx := logr.NewContext(context.TODO(), log)

			_, err := reconciler.Reconcile(ctx, test.request)
			require.NoError(t, err)

			if len(test.expectedLog) > 0 {
				assert.Contains(t, logBuffer.String(), test.expectedLog)
			}

			if !test.skipVerifyObject {
				permissions, err := fakeClientSet.PermissionsV1alpha1().AzPermissions().Get(context.TODO(), test.request.Name, v1.GetOptions{})
				require.NoError(t, err)

				require.NotNil(t, permissions.Status)
				require.Equal(t, test.expectedState, permissions.Status.State)

				switch permissions.Status.State {
				case v1alpha1.Ready:
					assert.Equal(t, permissions.Status.Allowed, test.expectedAllowed)
					assert.Equal(t, permissions.Status.Denied, test.expectedDenied)
				case v1alpha1.Failed:
					assert.Equal(t, test.expectedStateReason, permissions.Status.StateReason)
					assert.Equal(t, test.expectedStateMessage, permissions.Status.StateMessage)
				}
			}
		})
	}
}
