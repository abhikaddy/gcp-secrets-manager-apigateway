# gcp-secrets-manager-apigateway

We know GCP Secrets Manager already has an API to interact with out of the box. However, that requires the client to have an access_token which is created from a specially generated service account key or Application Default Credentials. However, if the client wants a plain REST API call, this might be lot of work on their side.

This project implements an API Gateway in Go which interacts with GCP secrets Manager using Workload Identity. Workload Identity configuration steps can be found here: https://cloud.google.com/kubernetes-engine/docs/how-to/workload-identity

Needless to say, this code will only work when container is deployed on Google Cloud Service eg: GKE or CLOUD RUN.

go mod init gcpsm-apigateway
go mod tidy
go mod run main.go