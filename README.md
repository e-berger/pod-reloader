# pod-reloader

Help you reload deployment/pod that use latest tag for your images.
Setup namespace to watch with annotations : `"pod-reloader/watch": "true"`
You can ignore workload in watched namespaces with this annotations on pod : `"pod-reloader/ignore": "true"`

## Setup

In `deploy` folder, you can find a chart helm to install this project

### Environment variables

LOGLEVEL : Set log level for the application (default: INFO), possible values: DEBUG, INFO, WARN, ERROR
REGISTRY : Set type of registry to use for images checking (required), possible values: DOCKER, ECR (for check by aws sdk instead of docker)
REGISTRY_AUTH_SECRET: Set the name of the secret to use for registry authentication (can be omitted if no authentication is required)
FREQUENCY_CHECK_SECONDS: Set the frequency of main loop in seconds (default: 30)

### Local development

You can run the application with `go run cmd/main.go` and set environment variables.

KUBECONFIG : Set the path to the kubeconfig file
AWS_REGION : Set the region to use for ECR authentication
AWS_ASSUME_ROLE : Set the role name to assume for aws authentication
AWS_PROFILE : Set the profile to use for aws authentication
