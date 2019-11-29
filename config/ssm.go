// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package config

import (
	"fmt"
	"strconv"
	"strings"

	log "github.com/communitybridge/easycla-api/logging"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

// getSSMString is a generic routine to fetch the specified key value
func getSSMString(ssmClient *ssm.SSM, key string) (string, error) {
	log.Debugf("Loading SSM parameter: %s", key)
	value, err := ssmClient.GetParameter(&ssm.GetParameterInput{
		Name:           aws.String(key),
		WithDecryption: aws.Bool(false),
	})
	if err != nil {
		log.Warnf("unable to read SSM parameter %s - error: %+v", key, err)
		return "", err
	}

	return strings.TrimSpace(*value.Parameter.Value), nil
}

func loadSSMConfig(awsSession *session.Session, stage string) (Config, error) {
	config := Config{}

	ssmClient := ssm.New(awsSession)

	var err error
	config.RDSHost, err = getSSMString(ssmClient, fmt.Sprintf("cla-rds-host-%s", stage))
	if err != nil {
		return Config{}, err
	}

	strPort, err := getSSMString(ssmClient, fmt.Sprintf("cla-rds-port-%s", stage))
	if err != nil {
		return Config{}, err
	}
	config.RDSPort, err = strconv.Atoi(strPort)
	if err != nil {
		return Config{}, err
	}

	config.RDSDatabase, err = getSSMString(ssmClient, fmt.Sprintf("cla-rds-database-%s", stage))
	if err != nil {
		return Config{}, err
	}

	config.RDSUsername, err = getSSMString(ssmClient, fmt.Sprintf("cla-rds-username-%s", stage))
	if err != nil {
		return Config{}, err
	}

	config.RDSPassword, err = getSSMString(ssmClient, fmt.Sprintf("cla-rds-password-%s", stage))
	if err != nil {
		return Config{}, err
	}

	config.GithubWebhookSecret, err = getSSMString(ssmClient, fmt.Sprintf("cla-github-webhook-secret-%s", stage))
	if err != nil {
		return Config{}, err
	}

	config.GithubAppPrivateKey, err = getSSMString(ssmClient, fmt.Sprintf("cla-github-app-private-key-%s", stage))
	if err != nil {
		return Config{}, err
	}
	return config, nil
}
