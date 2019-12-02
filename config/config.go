// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package config

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws/session"
	log "github.com/communitybridge/easycla-api/logging"
)

// Config data model
type Config struct {
	// RDS
	RDSHost             string `json:"rds_host"`
	RDSDatabase         string `json:"rds_database"`
	RDSUsername         string `json:"rds_username"`
	RDSPassword         string `json:"rds_password"`
	RDSPort             int    `json:"rds_port"`
	GithubWebhookSecret string `json:"github_webhook_secret"`
	GithubAppPrivateKey string `json:"github_app_private_key"`
	GithubAppID         int    `json:"github_app_id"`
}

// LoadConfig loads the configuration
func LoadConfig(configFilePath string, awsSession *session.Session, awsStage string) (Config, error) {
	var config Config
	var err error

	if configFilePath != "" {
		// Read from local env.jso
		log.Infof("Loading local config from file: %s", configFilePath)
		config, err = loadLocalConfig(configFilePath)

	} else if awsSession != nil {
		// Read from SSM
		log.Info("Loading SSM config...")
		config, err = loadSSMConfig(awsSession, awsStage)
	} else {
		return Config{}, errors.New("config not found")
	}

	if err != nil {
		log.Warnf("Error fetching SSM parameters for configuration, error: %+v", err)
		return Config{}, err
	}

	return config, nil
}
