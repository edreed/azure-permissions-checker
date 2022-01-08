package network

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
//
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"context"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/tracing"
	"net/http"
)

// P2sVpnGatewaysClient is the network Client
type P2sVpnGatewaysClient struct {
	BaseClient
}

// NewP2sVpnGatewaysClient creates an instance of the P2sVpnGatewaysClient client.
func NewP2sVpnGatewaysClient(subscriptionID string) P2sVpnGatewaysClient {
	return NewP2sVpnGatewaysClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewP2sVpnGatewaysClientWithBaseURI creates an instance of the P2sVpnGatewaysClient client using a custom endpoint.
// Use this when interacting with an Azure cloud that uses a non-standard base URI (sovereign clouds, Azure stack).
func NewP2sVpnGatewaysClientWithBaseURI(baseURI string, subscriptionID string) P2sVpnGatewaysClient {
	return P2sVpnGatewaysClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// CreateOrUpdate creates a virtual wan p2s vpn gateway if it doesn't exist else updates the existing gateway.
// Parameters:
// resourceGroupName - the resource group name of the P2SVpnGateway.
// gatewayName - the name of the gateway.
// p2SVpnGatewayParameters - parameters supplied to create or Update a virtual wan p2s vpn gateway.
func (client P2sVpnGatewaysClient) CreateOrUpdate(ctx context.Context, resourceGroupName string, gatewayName string, p2SVpnGatewayParameters P2SVpnGateway) (result P2sVpnGatewaysCreateOrUpdateFuture, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/P2sVpnGatewaysClient.CreateOrUpdate")
		defer func() {
			sc := -1
			if result.FutureAPI != nil && result.FutureAPI.Response() != nil {
				sc = result.FutureAPI.Response().StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.CreateOrUpdatePreparer(ctx, resourceGroupName, gatewayName, p2SVpnGatewayParameters)
	if err != nil {
		err = autorest.NewErrorWithError(err, "network.P2sVpnGatewaysClient", "CreateOrUpdate", nil, "Failure preparing request")
		return
	}

	result, err = client.CreateOrUpdateSender(req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "network.P2sVpnGatewaysClient", "CreateOrUpdate", nil, "Failure sending request")
		return
	}

	return
}

// CreateOrUpdatePreparer prepares the CreateOrUpdate request.
func (client P2sVpnGatewaysClient) CreateOrUpdatePreparer(ctx context.Context, resourceGroupName string, gatewayName string, p2SVpnGatewayParameters P2SVpnGateway) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"gatewayName":       autorest.Encode("path", gatewayName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2021-02-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	p2SVpnGatewayParameters.Etag = nil
	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/p2svpnGateways/{gatewayName}", pathParameters),
		autorest.WithJSON(p2SVpnGatewayParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// CreateOrUpdateSender sends the CreateOrUpdate request. The method will close the
// http.Response Body if it receives an error.
func (client P2sVpnGatewaysClient) CreateOrUpdateSender(req *http.Request) (future P2sVpnGatewaysCreateOrUpdateFuture, err error) {
	var resp *http.Response
	resp, err = client.Send(req, azure.DoRetryWithRegistration(client.Client))
	if err != nil {
		return
	}
	var azf azure.Future
	azf, err = azure.NewFutureFromResponse(resp)
	future.FutureAPI = &azf
	future.Result = future.result
	return
}

// CreateOrUpdateResponder handles the response to the CreateOrUpdate request. The method always
// closes the http.Response Body.
func (client P2sVpnGatewaysClient) CreateOrUpdateResponder(resp *http.Response) (result P2SVpnGateway, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusCreated),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// Delete deletes a virtual wan p2s vpn gateway.
// Parameters:
// resourceGroupName - the resource group name of the P2SVpnGateway.
// gatewayName - the name of the gateway.
func (client P2sVpnGatewaysClient) Delete(ctx context.Context, resourceGroupName string, gatewayName string) (result P2sVpnGatewaysDeleteFuture, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/P2sVpnGatewaysClient.Delete")
		defer func() {
			sc := -1
			if result.FutureAPI != nil && result.FutureAPI.Response() != nil {
				sc = result.FutureAPI.Response().StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.DeletePreparer(ctx, resourceGroupName, gatewayName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "network.P2sVpnGatewaysClient", "Delete", nil, "Failure preparing request")
		return
	}

	result, err = client.DeleteSender(req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "network.P2sVpnGatewaysClient", "Delete", nil, "Failure sending request")
		return
	}

	return
}

// DeletePreparer prepares the Delete request.
func (client P2sVpnGatewaysClient) DeletePreparer(ctx context.Context, resourceGroupName string, gatewayName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"gatewayName":       autorest.Encode("path", gatewayName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2021-02-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsDelete(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/p2svpnGateways/{gatewayName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// DeleteSender sends the Delete request. The method will close the
// http.Response Body if it receives an error.
func (client P2sVpnGatewaysClient) DeleteSender(req *http.Request) (future P2sVpnGatewaysDeleteFuture, err error) {
	var resp *http.Response
	resp, err = client.Send(req, azure.DoRetryWithRegistration(client.Client))
	if err != nil {
		return
	}
	var azf azure.Future
	azf, err = azure.NewFutureFromResponse(resp)
	future.FutureAPI = &azf
	future.Result = future.result
	return
}

// DeleteResponder handles the response to the Delete request. The method always
// closes the http.Response Body.
func (client P2sVpnGatewaysClient) DeleteResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusAccepted, http.StatusNoContent),
		autorest.ByClosing())
	result.Response = resp
	return
}

// DisconnectP2sVpnConnections disconnect P2S vpn connections of the virtual wan P2SVpnGateway in the specified
// resource group.
// Parameters:
// resourceGroupName - the name of the resource group.
// p2sVpnGatewayName - the name of the P2S Vpn Gateway.
// request - the parameters are supplied to disconnect p2s vpn connections.
func (client P2sVpnGatewaysClient) DisconnectP2sVpnConnections(ctx context.Context, resourceGroupName string, p2sVpnGatewayName string, request P2SVpnConnectionRequest) (result P2sVpnGatewaysDisconnectP2sVpnConnectionsFuture, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/P2sVpnGatewaysClient.DisconnectP2sVpnConnections")
		defer func() {
			sc := -1
			if result.FutureAPI != nil && result.FutureAPI.Response() != nil {
				sc = result.FutureAPI.Response().StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.DisconnectP2sVpnConnectionsPreparer(ctx, resourceGroupName, p2sVpnGatewayName, request)
	if err != nil {
		err = autorest.NewErrorWithError(err, "network.P2sVpnGatewaysClient", "DisconnectP2sVpnConnections", nil, "Failure preparing request")
		return
	}

	result, err = client.DisconnectP2sVpnConnectionsSender(req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "network.P2sVpnGatewaysClient", "DisconnectP2sVpnConnections", nil, "Failure sending request")
		return
	}

	return
}

// DisconnectP2sVpnConnectionsPreparer prepares the DisconnectP2sVpnConnections request.
func (client P2sVpnGatewaysClient) DisconnectP2sVpnConnectionsPreparer(ctx context.Context, resourceGroupName string, p2sVpnGatewayName string, request P2SVpnConnectionRequest) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"p2sVpnGatewayName": autorest.Encode("path", p2sVpnGatewayName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2021-02-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/p2svpnGateways/{p2sVpnGatewayName}/disconnectP2sVpnConnections", pathParameters),
		autorest.WithJSON(request),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// DisconnectP2sVpnConnectionsSender sends the DisconnectP2sVpnConnections request. The method will close the
// http.Response Body if it receives an error.
func (client P2sVpnGatewaysClient) DisconnectP2sVpnConnectionsSender(req *http.Request) (future P2sVpnGatewaysDisconnectP2sVpnConnectionsFuture, err error) {
	var resp *http.Response
	resp, err = client.Send(req, azure.DoRetryWithRegistration(client.Client))
	if err != nil {
		return
	}
	var azf azure.Future
	azf, err = azure.NewFutureFromResponse(resp)
	future.FutureAPI = &azf
	future.Result = future.result
	return
}

// DisconnectP2sVpnConnectionsResponder handles the response to the DisconnectP2sVpnConnections request. The method always
// closes the http.Response Body.
func (client P2sVpnGatewaysClient) DisconnectP2sVpnConnectionsResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusAccepted),
		autorest.ByClosing())
	result.Response = resp
	return
}

// GenerateVpnProfile generates VPN profile for P2S client of the P2SVpnGateway in the specified resource group.
// Parameters:
// resourceGroupName - the name of the resource group.
// gatewayName - the name of the P2SVpnGateway.
// parameters - parameters supplied to the generate P2SVpnGateway VPN client package operation.
func (client P2sVpnGatewaysClient) GenerateVpnProfile(ctx context.Context, resourceGroupName string, gatewayName string, parameters P2SVpnProfileParameters) (result P2sVpnGatewaysGenerateVpnProfileFuture, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/P2sVpnGatewaysClient.GenerateVpnProfile")
		defer func() {
			sc := -1
			if result.FutureAPI != nil && result.FutureAPI.Response() != nil {
				sc = result.FutureAPI.Response().StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.GenerateVpnProfilePreparer(ctx, resourceGroupName, gatewayName, parameters)
	if err != nil {
		err = autorest.NewErrorWithError(err, "network.P2sVpnGatewaysClient", "GenerateVpnProfile", nil, "Failure preparing request")
		return
	}

	result, err = client.GenerateVpnProfileSender(req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "network.P2sVpnGatewaysClient", "GenerateVpnProfile", nil, "Failure sending request")
		return
	}

	return
}

// GenerateVpnProfilePreparer prepares the GenerateVpnProfile request.
func (client P2sVpnGatewaysClient) GenerateVpnProfilePreparer(ctx context.Context, resourceGroupName string, gatewayName string, parameters P2SVpnProfileParameters) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"gatewayName":       autorest.Encode("path", gatewayName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2021-02-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/p2svpnGateways/{gatewayName}/generatevpnprofile", pathParameters),
		autorest.WithJSON(parameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GenerateVpnProfileSender sends the GenerateVpnProfile request. The method will close the
// http.Response Body if it receives an error.
func (client P2sVpnGatewaysClient) GenerateVpnProfileSender(req *http.Request) (future P2sVpnGatewaysGenerateVpnProfileFuture, err error) {
	var resp *http.Response
	resp, err = client.Send(req, azure.DoRetryWithRegistration(client.Client))
	if err != nil {
		return
	}
	var azf azure.Future
	azf, err = azure.NewFutureFromResponse(resp)
	future.FutureAPI = &azf
	future.Result = future.result
	return
}

// GenerateVpnProfileResponder handles the response to the GenerateVpnProfile request. The method always
// closes the http.Response Body.
func (client P2sVpnGatewaysClient) GenerateVpnProfileResponder(resp *http.Response) (result VpnProfileResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusAccepted),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// Get retrieves the details of a virtual wan p2s vpn gateway.
// Parameters:
// resourceGroupName - the resource group name of the P2SVpnGateway.
// gatewayName - the name of the gateway.
func (client P2sVpnGatewaysClient) Get(ctx context.Context, resourceGroupName string, gatewayName string) (result P2SVpnGateway, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/P2sVpnGatewaysClient.Get")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.GetPreparer(ctx, resourceGroupName, gatewayName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "network.P2sVpnGatewaysClient", "Get", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "network.P2sVpnGatewaysClient", "Get", resp, "Failure sending request")
		return
	}

	result, err = client.GetResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "network.P2sVpnGatewaysClient", "Get", resp, "Failure responding to request")
		return
	}

	return
}

// GetPreparer prepares the Get request.
func (client P2sVpnGatewaysClient) GetPreparer(ctx context.Context, resourceGroupName string, gatewayName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"gatewayName":       autorest.Encode("path", gatewayName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2021-02-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/p2svpnGateways/{gatewayName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetSender sends the Get request. The method will close the
// http.Response Body if it receives an error.
func (client P2sVpnGatewaysClient) GetSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (client P2sVpnGatewaysClient) GetResponder(resp *http.Response) (result P2SVpnGateway, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// GetP2sVpnConnectionHealth gets the connection health of P2S clients of the virtual wan P2SVpnGateway in the
// specified resource group.
// Parameters:
// resourceGroupName - the name of the resource group.
// gatewayName - the name of the P2SVpnGateway.
func (client P2sVpnGatewaysClient) GetP2sVpnConnectionHealth(ctx context.Context, resourceGroupName string, gatewayName string) (result P2sVpnGatewaysGetP2sVpnConnectionHealthFuture, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/P2sVpnGatewaysClient.GetP2sVpnConnectionHealth")
		defer func() {
			sc := -1
			if result.FutureAPI != nil && result.FutureAPI.Response() != nil {
				sc = result.FutureAPI.Response().StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.GetP2sVpnConnectionHealthPreparer(ctx, resourceGroupName, gatewayName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "network.P2sVpnGatewaysClient", "GetP2sVpnConnectionHealth", nil, "Failure preparing request")
		return
	}

	result, err = client.GetP2sVpnConnectionHealthSender(req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "network.P2sVpnGatewaysClient", "GetP2sVpnConnectionHealth", nil, "Failure sending request")
		return
	}

	return
}

// GetP2sVpnConnectionHealthPreparer prepares the GetP2sVpnConnectionHealth request.
func (client P2sVpnGatewaysClient) GetP2sVpnConnectionHealthPreparer(ctx context.Context, resourceGroupName string, gatewayName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"gatewayName":       autorest.Encode("path", gatewayName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2021-02-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsPost(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/p2svpnGateways/{gatewayName}/getP2sVpnConnectionHealth", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetP2sVpnConnectionHealthSender sends the GetP2sVpnConnectionHealth request. The method will close the
// http.Response Body if it receives an error.
func (client P2sVpnGatewaysClient) GetP2sVpnConnectionHealthSender(req *http.Request) (future P2sVpnGatewaysGetP2sVpnConnectionHealthFuture, err error) {
	var resp *http.Response
	resp, err = client.Send(req, azure.DoRetryWithRegistration(client.Client))
	if err != nil {
		return
	}
	var azf azure.Future
	azf, err = azure.NewFutureFromResponse(resp)
	future.FutureAPI = &azf
	future.Result = future.result
	return
}

// GetP2sVpnConnectionHealthResponder handles the response to the GetP2sVpnConnectionHealth request. The method always
// closes the http.Response Body.
func (client P2sVpnGatewaysClient) GetP2sVpnConnectionHealthResponder(resp *http.Response) (result P2SVpnGateway, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusAccepted),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// GetP2sVpnConnectionHealthDetailed gets the sas url to get the connection health detail of P2S clients of the virtual
// wan P2SVpnGateway in the specified resource group.
// Parameters:
// resourceGroupName - the name of the resource group.
// gatewayName - the name of the P2SVpnGateway.
// request - request parameters supplied to get p2s vpn connections detailed health.
func (client P2sVpnGatewaysClient) GetP2sVpnConnectionHealthDetailed(ctx context.Context, resourceGroupName string, gatewayName string, request P2SVpnConnectionHealthRequest) (result P2sVpnGatewaysGetP2sVpnConnectionHealthDetailedFuture, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/P2sVpnGatewaysClient.GetP2sVpnConnectionHealthDetailed")
		defer func() {
			sc := -1
			if result.FutureAPI != nil && result.FutureAPI.Response() != nil {
				sc = result.FutureAPI.Response().StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.GetP2sVpnConnectionHealthDetailedPreparer(ctx, resourceGroupName, gatewayName, request)
	if err != nil {
		err = autorest.NewErrorWithError(err, "network.P2sVpnGatewaysClient", "GetP2sVpnConnectionHealthDetailed", nil, "Failure preparing request")
		return
	}

	result, err = client.GetP2sVpnConnectionHealthDetailedSender(req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "network.P2sVpnGatewaysClient", "GetP2sVpnConnectionHealthDetailed", nil, "Failure sending request")
		return
	}

	return
}

// GetP2sVpnConnectionHealthDetailedPreparer prepares the GetP2sVpnConnectionHealthDetailed request.
func (client P2sVpnGatewaysClient) GetP2sVpnConnectionHealthDetailedPreparer(ctx context.Context, resourceGroupName string, gatewayName string, request P2SVpnConnectionHealthRequest) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"gatewayName":       autorest.Encode("path", gatewayName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2021-02-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/p2svpnGateways/{gatewayName}/getP2sVpnConnectionHealthDetailed", pathParameters),
		autorest.WithJSON(request),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetP2sVpnConnectionHealthDetailedSender sends the GetP2sVpnConnectionHealthDetailed request. The method will close the
// http.Response Body if it receives an error.
func (client P2sVpnGatewaysClient) GetP2sVpnConnectionHealthDetailedSender(req *http.Request) (future P2sVpnGatewaysGetP2sVpnConnectionHealthDetailedFuture, err error) {
	var resp *http.Response
	resp, err = client.Send(req, azure.DoRetryWithRegistration(client.Client))
	if err != nil {
		return
	}
	var azf azure.Future
	azf, err = azure.NewFutureFromResponse(resp)
	future.FutureAPI = &azf
	future.Result = future.result
	return
}

// GetP2sVpnConnectionHealthDetailedResponder handles the response to the GetP2sVpnConnectionHealthDetailed request. The method always
// closes the http.Response Body.
func (client P2sVpnGatewaysClient) GetP2sVpnConnectionHealthDetailedResponder(resp *http.Response) (result P2SVpnConnectionHealth, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusAccepted),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// List lists all the P2SVpnGateways in a subscription.
func (client P2sVpnGatewaysClient) List(ctx context.Context) (result ListP2SVpnGatewaysResultPage, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/P2sVpnGatewaysClient.List")
		defer func() {
			sc := -1
			if result.lpvgr.Response.Response != nil {
				sc = result.lpvgr.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	result.fn = client.listNextResults
	req, err := client.ListPreparer(ctx)
	if err != nil {
		err = autorest.NewErrorWithError(err, "network.P2sVpnGatewaysClient", "List", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListSender(req)
	if err != nil {
		result.lpvgr.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "network.P2sVpnGatewaysClient", "List", resp, "Failure sending request")
		return
	}

	result.lpvgr, err = client.ListResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "network.P2sVpnGatewaysClient", "List", resp, "Failure responding to request")
		return
	}
	if result.lpvgr.hasNextLink() && result.lpvgr.IsEmpty() {
		err = result.NextWithContext(ctx)
		return
	}

	return
}

// ListPreparer prepares the List request.
func (client P2sVpnGatewaysClient) ListPreparer(ctx context.Context) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"subscriptionId": autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2021-02-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/providers/Microsoft.Network/p2svpnGateways", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListSender sends the List request. The method will close the
// http.Response Body if it receives an error.
func (client P2sVpnGatewaysClient) ListSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// ListResponder handles the response to the List request. The method always
// closes the http.Response Body.
func (client P2sVpnGatewaysClient) ListResponder(resp *http.Response) (result ListP2SVpnGatewaysResult, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// listNextResults retrieves the next set of results, if any.
func (client P2sVpnGatewaysClient) listNextResults(ctx context.Context, lastResults ListP2SVpnGatewaysResult) (result ListP2SVpnGatewaysResult, err error) {
	req, err := lastResults.listP2SVpnGatewaysResultPreparer(ctx)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "network.P2sVpnGatewaysClient", "listNextResults", nil, "Failure preparing next results request")
	}
	if req == nil {
		return
	}
	resp, err := client.ListSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "network.P2sVpnGatewaysClient", "listNextResults", resp, "Failure sending next results request")
	}
	result, err = client.ListResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "network.P2sVpnGatewaysClient", "listNextResults", resp, "Failure responding to next results request")
	}
	return
}

// ListComplete enumerates all values, automatically crossing page boundaries as required.
func (client P2sVpnGatewaysClient) ListComplete(ctx context.Context) (result ListP2SVpnGatewaysResultIterator, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/P2sVpnGatewaysClient.List")
		defer func() {
			sc := -1
			if result.Response().Response.Response != nil {
				sc = result.page.Response().Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	result.page, err = client.List(ctx)
	return
}

// ListByResourceGroup lists all the P2SVpnGateways in a resource group.
// Parameters:
// resourceGroupName - the resource group name of the P2SVpnGateway.
func (client P2sVpnGatewaysClient) ListByResourceGroup(ctx context.Context, resourceGroupName string) (result ListP2SVpnGatewaysResultPage, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/P2sVpnGatewaysClient.ListByResourceGroup")
		defer func() {
			sc := -1
			if result.lpvgr.Response.Response != nil {
				sc = result.lpvgr.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	result.fn = client.listByResourceGroupNextResults
	req, err := client.ListByResourceGroupPreparer(ctx, resourceGroupName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "network.P2sVpnGatewaysClient", "ListByResourceGroup", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListByResourceGroupSender(req)
	if err != nil {
		result.lpvgr.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "network.P2sVpnGatewaysClient", "ListByResourceGroup", resp, "Failure sending request")
		return
	}

	result.lpvgr, err = client.ListByResourceGroupResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "network.P2sVpnGatewaysClient", "ListByResourceGroup", resp, "Failure responding to request")
		return
	}
	if result.lpvgr.hasNextLink() && result.lpvgr.IsEmpty() {
		err = result.NextWithContext(ctx)
		return
	}

	return
}

// ListByResourceGroupPreparer prepares the ListByResourceGroup request.
func (client P2sVpnGatewaysClient) ListByResourceGroupPreparer(ctx context.Context, resourceGroupName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2021-02-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/p2svpnGateways", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListByResourceGroupSender sends the ListByResourceGroup request. The method will close the
// http.Response Body if it receives an error.
func (client P2sVpnGatewaysClient) ListByResourceGroupSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// ListByResourceGroupResponder handles the response to the ListByResourceGroup request. The method always
// closes the http.Response Body.
func (client P2sVpnGatewaysClient) ListByResourceGroupResponder(resp *http.Response) (result ListP2SVpnGatewaysResult, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// listByResourceGroupNextResults retrieves the next set of results, if any.
func (client P2sVpnGatewaysClient) listByResourceGroupNextResults(ctx context.Context, lastResults ListP2SVpnGatewaysResult) (result ListP2SVpnGatewaysResult, err error) {
	req, err := lastResults.listP2SVpnGatewaysResultPreparer(ctx)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "network.P2sVpnGatewaysClient", "listByResourceGroupNextResults", nil, "Failure preparing next results request")
	}
	if req == nil {
		return
	}
	resp, err := client.ListByResourceGroupSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "network.P2sVpnGatewaysClient", "listByResourceGroupNextResults", resp, "Failure sending next results request")
	}
	result, err = client.ListByResourceGroupResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "network.P2sVpnGatewaysClient", "listByResourceGroupNextResults", resp, "Failure responding to next results request")
	}
	return
}

// ListByResourceGroupComplete enumerates all values, automatically crossing page boundaries as required.
func (client P2sVpnGatewaysClient) ListByResourceGroupComplete(ctx context.Context, resourceGroupName string) (result ListP2SVpnGatewaysResultIterator, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/P2sVpnGatewaysClient.ListByResourceGroup")
		defer func() {
			sc := -1
			if result.Response().Response.Response != nil {
				sc = result.page.Response().Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	result.page, err = client.ListByResourceGroup(ctx, resourceGroupName)
	return
}

// Reset resets the primary of the p2s vpn gateway in the specified resource group.
// Parameters:
// resourceGroupName - the resource group name of the P2SVpnGateway.
// gatewayName - the name of the gateway.
func (client P2sVpnGatewaysClient) Reset(ctx context.Context, resourceGroupName string, gatewayName string) (result P2SVpnGatewaysResetFuture, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/P2sVpnGatewaysClient.Reset")
		defer func() {
			sc := -1
			if result.FutureAPI != nil && result.FutureAPI.Response() != nil {
				sc = result.FutureAPI.Response().StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.ResetPreparer(ctx, resourceGroupName, gatewayName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "network.P2sVpnGatewaysClient", "Reset", nil, "Failure preparing request")
		return
	}

	result, err = client.ResetSender(req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "network.P2sVpnGatewaysClient", "Reset", nil, "Failure sending request")
		return
	}

	return
}

// ResetPreparer prepares the Reset request.
func (client P2sVpnGatewaysClient) ResetPreparer(ctx context.Context, resourceGroupName string, gatewayName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"gatewayName":       autorest.Encode("path", gatewayName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2021-02-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsPost(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/p2svpnGateways/{gatewayName}/reset", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ResetSender sends the Reset request. The method will close the
// http.Response Body if it receives an error.
func (client P2sVpnGatewaysClient) ResetSender(req *http.Request) (future P2SVpnGatewaysResetFuture, err error) {
	var resp *http.Response
	resp, err = client.Send(req, azure.DoRetryWithRegistration(client.Client))
	if err != nil {
		return
	}
	var azf azure.Future
	azf, err = azure.NewFutureFromResponse(resp)
	future.FutureAPI = &azf
	future.Result = future.result
	return
}

// ResetResponder handles the response to the Reset request. The method always
// closes the http.Response Body.
func (client P2sVpnGatewaysClient) ResetResponder(resp *http.Response) (result P2SVpnGateway, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusAccepted),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// UpdateTags updates virtual wan p2s vpn gateway tags.
// Parameters:
// resourceGroupName - the resource group name of the P2SVpnGateway.
// gatewayName - the name of the gateway.
// p2SVpnGatewayParameters - parameters supplied to update a virtual wan p2s vpn gateway tags.
func (client P2sVpnGatewaysClient) UpdateTags(ctx context.Context, resourceGroupName string, gatewayName string, p2SVpnGatewayParameters TagsObject) (result P2sVpnGatewaysUpdateTagsFuture, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/P2sVpnGatewaysClient.UpdateTags")
		defer func() {
			sc := -1
			if result.FutureAPI != nil && result.FutureAPI.Response() != nil {
				sc = result.FutureAPI.Response().StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.UpdateTagsPreparer(ctx, resourceGroupName, gatewayName, p2SVpnGatewayParameters)
	if err != nil {
		err = autorest.NewErrorWithError(err, "network.P2sVpnGatewaysClient", "UpdateTags", nil, "Failure preparing request")
		return
	}

	result, err = client.UpdateTagsSender(req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "network.P2sVpnGatewaysClient", "UpdateTags", nil, "Failure sending request")
		return
	}

	return
}

// UpdateTagsPreparer prepares the UpdateTags request.
func (client P2sVpnGatewaysClient) UpdateTagsPreparer(ctx context.Context, resourceGroupName string, gatewayName string, p2SVpnGatewayParameters TagsObject) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"gatewayName":       autorest.Encode("path", gatewayName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2021-02-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPatch(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/p2svpnGateways/{gatewayName}", pathParameters),
		autorest.WithJSON(p2SVpnGatewayParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// UpdateTagsSender sends the UpdateTags request. The method will close the
// http.Response Body if it receives an error.
func (client P2sVpnGatewaysClient) UpdateTagsSender(req *http.Request) (future P2sVpnGatewaysUpdateTagsFuture, err error) {
	var resp *http.Response
	resp, err = client.Send(req, azure.DoRetryWithRegistration(client.Client))
	if err != nil {
		return
	}
	var azf azure.Future
	azf, err = azure.NewFutureFromResponse(resp)
	future.FutureAPI = &azf
	future.Result = future.result
	return
}

// UpdateTagsResponder handles the response to the UpdateTags request. The method always
// closes the http.Response Body.
func (client P2sVpnGatewaysClient) UpdateTagsResponder(resp *http.Response) (result P2SVpnGateway, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusAccepted),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
