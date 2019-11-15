// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package init

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws/credentials"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	log "github.com/communitybridge/easycla-api/logging"
)

const (
	// DefaultAWSRegion is the region to use when the AWS_REGION environment variable is not set
	DefaultAWSRegion = "us-east-1"
)

var (
	// AWS
	awsRegion          string
	awsAccessKeyID     string
	awsSecretAccessKey string
	awsSessionToken    string

	awsSession           *session.Session
	awsCloudWatchService *cloudwatch.CloudWatch
)

// AWSInit initialization logic for the AWS resources
func AWSInit() {
	awsRegion = os.Getenv("AWS_REGION")
	if awsRegion == "" {
		log.Debugf("AWS_REGION environment variable is not set. Using default region: %s", DefaultAWSRegion)
		awsRegion = DefaultAWSRegion
	} else {
		log.Debugf("AWS_REGION set to: %s", awsRegion)
	}

	key := fmt.Sprintf("%s_AWS_ACCESS_KEY_ID", ServiceName)
	awsAccessKeyID = os.Getenv(key)
	if awsAccessKeyID == "" {
		log.Fatalf("Unable to load %s - value not set in environment. Exiting...", key)
	}
	log.Debugf("Loaded %s: %s...", key, awsAccessKeyID[:5])

	key = fmt.Sprintf("%s_AWS_SECRET_ACCESS_KEY", ServiceName)
	awsSecretAccessKey = os.Getenv(key)
	if awsSecretAccessKey == "" {
		log.Fatalf("Unable to load %s - value not set in environment. Exiting...", key)
	}
	log.Debugf("Loaded %s: %s...", key, awsSecretAccessKey[:5])

	key = fmt.Sprintf("%s_AWS_SESSION_TOKEN", ServiceName)
	awsSessionToken = os.Getenv(key)
	if awsSessionToken != "" {
		log.Debugf("Loaded %s: %s...", key, awsSessionToken[:5])
	}

	if err := startCloudWatchSession(); err != nil {
		log.Fatalf("Error starting the AWS CloudWatch session - Error: %s", err.Error())
	}
}

// GetAWSSession returns an AWS session based on the region and credentials
func GetAWSSession() (*session.Session, error) {
	if awsSession == nil {
		log.Debugf("Creating a new AWS session for region: %s", awsRegion)
		awsSession = session.Must(session.NewSession(&aws.Config{
			Region:      aws.String(awsRegion),
			Credentials: credentials.NewStaticCredentials(awsAccessKeyID, awsSecretAccessKey, awsSessionToken),
			MaxRetries:  aws.Int(5),
		}))

		/*
			awsSession = session.Must(session.NewSession(
				&aws.Config{
					Region:                        aws.String(awsRegion),
					CredentialsChainVerboseErrors: aws.Bool(true),
					MaxRetries:                    aws.Int(5),
				},
			))
		*/

		log.Debugf("Successfully created a new AWS session for region: %s", awsRegion)
	}

	return awsSession, nil
}

// startCloudWatchSession creates a new AWS CloudWatch service session
func startCloudWatchSession() error {
	sess, err := GetAWSSession()
	if err != nil {
		log.Fatal("Error creating a new AWS Session", err)
		return err
	}

	awsCloudWatchService = cloudwatch.New(sess)

	log.Info("CloudWatch service started")

	return nil
}

// GetAWSCloudWatchService returns the CloudWatch service client
func GetAWSCloudWatchService() *cloudwatch.CloudWatch {
	return awsCloudWatchService
}
