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
