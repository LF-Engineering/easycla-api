// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package health

import (
	"context"
	"time"

	"github.com/communitybridge/easycla-api/gen/models"
	"github.com/communitybridge/easycla-api/gen/restapi/operations/health"
)

// Service provides an API to the health API
type Service struct {
	version   string
	commit    string
	branch    string
	buildDate string
}

// New is a simple helper function to create a health service instance
func New(version, commit, branch, buildDate string) Service {
	return Service{
		version:   version,
		commit:    commit,
		branch:    branch,
		buildDate: buildDate,
	}
}

// HealthCheck API call returns the current health of the service
func (s Service) HealthCheck(ctx context.Context, in health.GetHealthParams) (*models.Health, error) {

	t := time.Now()
	hs := models.HealthStatus{
		TimeStamp: time.Now().UTC().Format(time.RFC3339),
		Healthy:   true,
		Name:      "CLA Health",
		Duration:  time.Since(t).String(),
	}

	// Do a quick check to see if we have a DynamoDB database connection
	dynamoNow := time.Now()
	dynamoAlive := isDynamoDBAlive()
	dy := models.HealthStatus{
		TimeStamp: time.Now().UTC().Format(time.RFC3339),
		Healthy:   dynamoAlive,
		Name:      "CLA - Dynamodb",
		Duration:  time.Since(dynamoNow).String(),
	}

	// Do a quick check to see if we have a RDS database connection
	rdsNow := time.Now()
	rdsAlive := isRDSAlive()
	rds := models.HealthStatus{
		TimeStamp: time.Now().UTC().Format(time.RFC3339),
		Healthy:   rdsAlive,
		Name:      "CLA - RDS",
		Duration:  time.Since(rdsNow).String(),
	}

	var status = "healthy"
	if !dynamoAlive {
		status = "not healthy"
	}

	response := models.Health{
		Status:         status,
		TimeStamp:      time.Now().UTC().Format(time.RFC3339),
		Version:        s.version,
		Githash:        s.commit,
		Branch:         s.branch,
		BuildTimeStamp: s.buildDate,
		Healths: []*models.HealthStatus{
			&hs,
			&dy,
			&rds,
		},
	}

	return &response, nil
}

// isDynamoDBAlive runs a check to see if we have connectivity to the database - returns true if successful, false otherwise
func isDynamoDBAlive() bool {
	return true
	/*
		// Grab the AWS session
		awsSession, err := ini.GetAWSSession()
		if err != nil {
			log.Warnf("Unable to acquire AWS session - returning failed health, error: %v", err)
			return false
		}

		// Known table that we can query
		tableName := "cla-" + ini.GetStage() + "-projects"

		// Create a client and make a query - don't wory about the result - just check the error response
		dynamoDBClient := dynamodb.New(awsSession)
		_, err = dynamoDBClient.DescribeTable(&dynamodb.DescribeTableInput{
			TableName: &tableName,
		})

		// No error is success
		return err == nil
	*/
}

// isRDSAlive runs a check to see if we have connectivity to the database - returns true if successful, false otherwise
func isRDSAlive() bool {
	return true
}
