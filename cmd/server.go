// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package cmd

import (
	"net/http"
	"os"
	"runtime"

	"github.com/communitybridge/easycla-api/cla_groups"
	"github.com/communitybridge/easycla-api/config"
	"github.com/communitybridge/easycla-api/events"
	"github.com/communitybridge/easycla-api/orgs"
	"github.com/communitybridge/easycla-api/projects"

	"github.com/communitybridge/easycla-api/apidocs"
	"github.com/communitybridge/easycla-api/gen/restapi"
	"github.com/communitybridge/easycla-api/gen/restapi/operations"
	"github.com/communitybridge/easycla-api/health"
	ini "github.com/communitybridge/easycla-api/init"
	log "github.com/communitybridge/easycla-api/logging"

	"github.com/go-openapi/loads"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Version is the application version - either a git SHA or tag value
	Version string

	// Commit is the application commit hash
	Commit string

	// Branch the build branch
	Branch string

	// BuildDate is the date of the build
	BuildDate string
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "server",
	Short: "Run the backend server",
	Long:  `Run the backend server which listens for http requests over a given port.`,
	Run:   runServer,
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

// server function called by environment specific server functions
func server(localMode bool) http.Handler {

	host, err := os.Hostname()
	if err != nil {
		log.Fatalf("unable to get hostname. Error: %v", err)
	}

	stage := viper.GetString("STAGE")

	awsSession, err := ini.GetAWSSession()
	if err != nil {
		log.Panicf("Unable to load AWS session - Error: %v", err)
	}

	configFile, err := config.LoadConfig(configFile, awsSession, stage)
	if err != nil {
		log.Panicf("Unable to load config - Error: %v", err)
	}

	log.Infof("Service %s starting...", ini.ServiceName)

	// Show the version and build info
	log.Infof("Name                    : %s", ini.ServiceName)
	log.Infof("Version                 : %s", Version)
	log.Infof("Git commit hash         : %s", Commit)
	log.Infof("Branch                  : %s", Branch)
	log.Infof("Build date              : %s", BuildDate)
	log.Infof("Golang OS               : %s", runtime.GOOS)
	log.Infof("Golang Arch             : %s", runtime.GOARCH)
	log.Infof("LOCAL_MODE              : %t", localMode)
	log.Infof("STAGE                   : %s", stage)
	log.Infof("AWS_REGION              : %s", os.Getenv("AWS_REGION"))
	log.Infof("Service Host            : %s", host)
	log.Infof("Service Port (localonly): %d", *portFlag)
	log.Infof("RDS Host                : %s", configFile.RDSHost)
	log.Infof("RDS Database            : %s", configFile.RDSDatabase)
	log.Infof("RDS Username            : %s", configFile.RDSUsername)
	log.Infof("RDS Port                : %d", configFile.RDSPort)

	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Panicf("Invalid swagger file for initializing cla - Error: %v", err)
	}

	api := operations.NewClaAPI(swaggerSpec)

	// Initialize the DB connection
	db := initDB(configFile)

	healthRepo := health.NewRepository(db)
	healthService := health.New(healthRepo, Version, Commit, Branch, BuildDate)

	health.Configure(api, healthService)
	apidocs.Configure(api)

	orgsRepo := orgs.NewRepository()
	orgsService := orgs.NewService(orgsRepo)
	orgs.Configure(api, orgsService)

	projectsRepo := projects.NewRepository()
	projectsService := projects.NewService(projectsRepo)
	projects.Configure(api, projectsService)

	eventsRepo := events.NewRepository(db)
	eventsService := events.NewService(eventsRepo)
	events.Configure(api, eventsService)

	claGroupsRepo := cla_groups.NewRepository(db)
	claGroupsService := cla_groups.NewService(claGroupsRepo, eventsService)
	cla_groups.Configure(api, claGroupsService)
	return api.Serve(setupMiddlewares)
}

// setupMiddlewares The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return responseLoggingMiddleware(handler)
}

// LoggingResponseWriter is a wrapper around an http.ResponseWriter which captures the
// status code written to the response, so that it can be logged.
type LoggingResponseWriter struct {
	wrapped    http.ResponseWriter
	StatusCode int
	// Response content could also be captured here, but I was only interested in logging the response status code
}

// NewLoggingResponseWriter creates a new logging response writer
func NewLoggingResponseWriter(wrapped http.ResponseWriter) *LoggingResponseWriter {
	return &LoggingResponseWriter{wrapped: wrapped}
}

// Header returns the header
func (lrw *LoggingResponseWriter) Header() http.Header {
	return lrw.wrapped.Header()
}

// Write writes the contents
func (lrw *LoggingResponseWriter) Write(content []byte) (int, error) {
	return lrw.wrapped.Write(content)
}

// WriteHeader writes the header
func (lrw *LoggingResponseWriter) WriteHeader(statusCode int) {
	lrw.StatusCode = statusCode
	lrw.wrapped.WriteHeader(statusCode)
}

// responseLoggingMiddleware logs the responses from API endpoints
func responseLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(NewLoggingResponseWriter(w), r)
		if r.Response != nil {
			log.Debugf("%s %s, response code: %d response status: %s",
				r.Method, r.URL.String(), r.Response.StatusCode, r.Response.Status)
		} else {
			log.Debugf("%s %s", r.Method, r.URL.String())
		}
	})
}
