// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package init

import (
	"os"

	log "github.com/communitybridge/easycla-api/logging"
)

const (
	// ServiceName is the name of the service - this is used in logs and as a environment prefix
	ServiceName = "CLA_SERVICE"

	// DefaultStage is the default run-time environment, typically one of: dev, test, staging or prod
	DefaultStage = "dev"
)

var (
	stage string
)

// CommonInit initializes the common properties
func CommonInit() {
	stage = os.Getenv("STAGE")
	if stage == "" {
		log.Debugf("STAGE environment variable is not set. Using default: %s", DefaultStage)
		stage = DefaultStage
	} else {
		log.Debugf("STAGE set to: %s", DefaultStage)
	}
}

/*
// getProperty is a common routine to bind and return the specified environment variable
func getProperty(property string) string {
	err := viper.BindEnv(property)
	if err != nil {
		log.Fatalf("Unable to load property: %s - value not defined or empty", property)
	}

	value := viper.GetString(property)
	if value == "" {
		err := fmt.Errorf("%s environment variable cannot be empty", property)
		log.Fatal(err.Error())
	}

	return value
}
*/

// Init initialization logic for all the handlers
func Init() {
	CommonInit()
	AWSInit()
}

// GetStage returns the deployment stage, e.g. dev, test, stage or prod
func GetStage() string {
	return stage
}
