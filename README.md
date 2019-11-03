# easycla-api

[![CircleCI](https://circleci.com/gh/communitybridge/easycla-api/tree/master.svg?style=svg&circle-token=1e1a6bfbb182a6d2acd24a9886580670604aebc0)](https://circleci.com/gh/communitybridge/easycla-api/tree/master)

The Contributor License Agreement (CLA) service of the Linux Foundation lets
project contributors read, sign, and submit contributor license agreements easily.

This repository contains the backend API for supporting and
managing the application.

This platform supports both GitHub and Gerrit source code repositories.
Additional information can be found in the [Getting Started Guide](#getting-started-guide).

## Getting Started Guide

See the [Getting Started Guide](docs/getting-started.md) to get started with EasyCLA.

## Third-party Services

Besides integration with Auth0 and Salesforce, the CLA system has the following third party services:

- [Docusign](https://www.docusign.com/) for CLA agreement e-sign flow
- [Docraptor](https://docraptor.com/) for convert html CLA template as PDF file

## Architecture

See the [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) document.

## Building

To build the software, you will need the following tools installed.

- Go version >= 1.13
- GNU Make

### Tool Setup

To install the go build tools, run the setup [Makefile](Makefile) target.
This only needs to be done once or after a tool/version upgrade.

```bash
make setup
```

This will install:

- [Swagger](https://github.com/go-swagger/go-swagger) - A Swagger 2.0 implementation for Go
- [Go Imports](https://godoc.org/golang.org/x/tools/cmd/goimports) - A tool that updates your Go import lines, adding missing ones and removing unreferenced ones
- [Dep](https://github.com/golang/dep) - A Go dependency management tool
- [DBMate](https://github.com/amacneil/dbmate) - A lightweight, framework-agnostic database migration tool.
- [SafeSQL](https://github.com/stripe/safesql) - a linter for SQL

### Clean Build

To build the software from scratch, run:

```bash
make clean swagger swagger-validate deps fmt build-mac test lint
```

## Running

To run the software, you will need to setup your AWS account and export a few
environment variables:

- `STAGE` - the stage/environment for the deployment, typically one of: `dev`, `staging`, or `prod`
- `AWS_REGION` - the AWS region, typically: `us-east-1`
- `AWS_ACCESS_KEY_ID` - the AWS access key ID
- `AWS_SECRET_ACCESS_KEY` - the AWS secret access key

AWS credentials are required to load the environment configuration values
from the AWS SSM parameter store. Once these values are exported to your
environment, then run the executable generated during the build:

```bash
./bin/cla-api
```

Typical output would look like (`dev` environment example):

```code
Logging format not defined - setting value to default: 'text'
Logging configured with level: debug, format: text
INFO[2019-10-18T16:52:43Z] Running init...                              
INFO[2019-10-18T16:52:43Z] Staring the HTTP server in local mode...     
DEBU[2019-10-18T16:52:43Z] Creating a new AWS session for region: ...   
DEBU[2019-10-18T16:52:43Z] Successfully created a new AWS session for region: ... 
INFO[2019-10-18T16:52:43Z] Loading SSM config...                        
INFO[2019-10-18T16:52:45Z] Service CLA_API_SERVICE starting...          
INFO[2019-10-18T16:52:45Z] Name                    : CLA_API_SERVICE    
INFO[2019-10-18T16:52:45Z] Version                 : f9dd2ce            
INFO[2019-10-18T16:52:45Z] Git commit hash         : f9dd2ce            
INFO[2019-10-18T16:52:45Z] Branch                  : master             
INFO[2019-10-18T16:52:45Z] Build date              : 2019-10-18T09:45:35-0700 
INFO[2019-10-18T16:52:45Z] Golang OS               : darwin             
INFO[2019-10-18T16:52:45Z] Golang Arch             : amd64              
INFO[2019-10-18T16:52:45Z] LOCAL_MODE              : true               
INFO[2019-10-18T16:52:45Z] STAGE                   : dev                
INFO[2019-10-18T16:52:45Z] Service Host            : Davids-MBP.attlocal.net 
INFO[2019-10-18T16:52:45Z] Service Port            : 8080               
INFO[2019-10-18T16:52:45Z] Sender email address is : admin@dev.lfcla.com 
INFO[2019-10-18T16:52:45Z] RDS Database            : easycla-rds-database-dev 
INFO[2019-10-18T16:52:45Z] RDS Username            : <value redacted>
INFO[2019-10-18T16:52:45Z] Running http server on port: 8080 - set PORT environment variable to change port 
```

## Validating the Installation

To validate the EasyCLA API is up and running, make a request to the Health Service or the API documentation:

View the Health endpoint:

```bash
open http://localhost:8080/v4/ops/health
# via API Gateway
open https://api-gw.dev.platform.linuxfoundation.org/cla-service/v4/ops/health


# or use cURL
curl -s -XGET http://localhost:8080/v4/ops/health
# via API Gateway
curl -s -XGET https://api-gw.dev.platform.linuxfoundation.org/cla-service/v4/ops/health

# or use cURL with jq to make the output pretty
curl -s -XGET http://localhost:8080/v4/ops/health | jq
```

View the API documentation in a web browser:

```bash
open http://localhost:8080/v4/api-docs
```

## License

Copyright The Linux Foundation and each contributor to CommunityBridge.

This project’s source code is licensed under the MIT License. A copy of the
license is available in LICENSE.

The project includes source code from keycloak, which is licensed under the
Apache License, version 2.0 (Apache-2.0), a copy of which is available in
LICENSE-keycloak.

This project’s documentation is licensed under the Creative Commons Attribution
4.0 International License (CC-BY-4.0). A copy of the license is available in
LICENSE-docs.
