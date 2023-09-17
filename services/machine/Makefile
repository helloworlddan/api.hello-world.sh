SERVICE := machine
PROJECT := $(shell terraform -chdir=../../infrastructure output -json | jq -r '.project.value')
REGION  := $(shell terraform -chdir=../../infrastructure output -json | jq -r '.region.value')
PREFIX  := $(shell terraform -chdir=../../infrastructure output -json | jq -r '.prefix.value')

run:
	ENVIRONMENT=dev GOOGLE_CLOUD_PROJECT=${PROJECT} TOP_SESSION=ca683f00-d51c-4f1a-af5e-5f9a25b3f4a8 TOP_OWNER=dan@hello-world.sh CGO_ENABLED=0 go run *.go

build:
	gcloud services enable cloudbuild.googleapis.com
	gcloud builds submit --tag=gcr.io/${PROJECT}/${SERVICE} .

update:
	gcloud run services update ${PREFIX}-${SERVICE} --region=${REGION} --image=gcr.io/${PROJECT}/${SERVICE}

.PHONY: run build update
