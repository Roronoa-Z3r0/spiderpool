// Code generated by go-swagger; DO NOT EDIT.

// Copyright 2022 Authors of spidernet-io
// SPDX-License-Identifier: Apache-2.0

package daemonset

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
)

// NewGetWorkloadendpointParams creates a new GetWorkloadendpointParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewGetWorkloadendpointParams() *GetWorkloadendpointParams {
	return &GetWorkloadendpointParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewGetWorkloadendpointParamsWithTimeout creates a new GetWorkloadendpointParams object
// with the ability to set a timeout on a request.
func NewGetWorkloadendpointParamsWithTimeout(timeout time.Duration) *GetWorkloadendpointParams {
	return &GetWorkloadendpointParams{
		timeout: timeout,
	}
}

// NewGetWorkloadendpointParamsWithContext creates a new GetWorkloadendpointParams object
// with the ability to set a context for a request.
func NewGetWorkloadendpointParamsWithContext(ctx context.Context) *GetWorkloadendpointParams {
	return &GetWorkloadendpointParams{
		Context: ctx,
	}
}

// NewGetWorkloadendpointParamsWithHTTPClient creates a new GetWorkloadendpointParams object
// with the ability to set a custom HTTPClient for a request.
func NewGetWorkloadendpointParamsWithHTTPClient(client *http.Client) *GetWorkloadendpointParams {
	return &GetWorkloadendpointParams{
		HTTPClient: client,
	}
}

/* GetWorkloadendpointParams contains all the parameters to send to the API endpoint
   for the get workloadendpoint operation.

   Typically these are written to a http.Request.
*/
type GetWorkloadendpointParams struct {
	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the get workloadendpoint params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetWorkloadendpointParams) WithDefaults() *GetWorkloadendpointParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the get workloadendpoint params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetWorkloadendpointParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the get workloadendpoint params
func (o *GetWorkloadendpointParams) WithTimeout(timeout time.Duration) *GetWorkloadendpointParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get workloadendpoint params
func (o *GetWorkloadendpointParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get workloadendpoint params
func (o *GetWorkloadendpointParams) WithContext(ctx context.Context) *GetWorkloadendpointParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get workloadendpoint params
func (o *GetWorkloadendpointParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get workloadendpoint params
func (o *GetWorkloadendpointParams) WithHTTPClient(client *http.Client) *GetWorkloadendpointParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get workloadendpoint params
func (o *GetWorkloadendpointParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WriteToRequest writes these params to a swagger request
func (o *GetWorkloadendpointParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
