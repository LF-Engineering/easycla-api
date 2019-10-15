// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package apidocs

import (
	"net/http"

	"github.com/go-openapi/runtime"
)

// GetDocOK Success
type GetDocOK struct {
}

// NewGetDocOK creates GetDocOK with default headers values
func NewGetDocOK() *GetDocOK {
	return &GetDocOK{}
}

// WriteResponse to the client
func (o *GetDocOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {
	html := `<!DOCTYPE html>
    <html lang="en">
	  <head>
		<title>CLA API Service Doc</title>
		<!-- needed for adaptive design -->
		<meta charset="utf-8"/>
		<meta name="viewport" content="width=device-width, initial-scale=1">
		<link rel="shortcut icon" href="https://www.linuxfoundation.org/wp-content/uploads/2017/08/favicon.png">
		<link href="https://fonts.googleapis.com/css?family=Montserrat:300,400,700|Roboto:300,400,700" rel="stylesheet">

		<!--
		ReDoc doesn't change outer page styles
		-->
		<style>
		  body {
			margin: 0;
			padding: 0;
		  }
		</style>
	  </head>
	  <body>
		<redoc spec-url='/v4/swagger.json'></redoc>
		<script src="https://cdn.jsdelivr.net/npm/redoc@next/bundles/redoc.standalone.js"> </script>
	  </body>
	</html>`

	rw.Header().Set("Content-Type", "text/html")
	_, err := rw.Write([]byte(html))
	if err != nil {
		panic(err)
	}
}
