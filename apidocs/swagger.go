// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package apidocs

import (
	"net/http"

	"github.com/communitybridge/easycla-api/gen/restapi"

	"github.com/go-openapi/runtime"
)

// GetSwaggerOK Success
type GetSwaggerOK struct {
}

// NewGetSwaggerOK creates GetSwaggerOK with default headers values
func NewGetSwaggerOK() *GetSwaggerOK {
	return &GetSwaggerOK{}
}

// WriteResponse streams the swagger.json to the client
func (o *GetSwaggerOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {
	rw.Header().Set("Content-Type", "application/json")
	_, err := rw.Write([]byte(restapi.SwaggerJSON))
	if err != nil {
		panic(err)
	}
}
