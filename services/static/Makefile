SERVICE := static
PROJECT := $(shell terraform -chdir=../../infrastructure output -json | jq -r '.project.value')
REGION := $(shell terraform -chdir=../../infrastructure output -json | jq -r '.region.value')
PREFIX := $(shell terraform -chdir=../../infrastructure output -json | jq -r '.prefix.value')
GATEWAY_SA := $(shell terraform -chdir=../../infrastructure output -json | jq -r '.gateway_sa.value')

run:
	go run server.go

build:
	gcloud builds submit --tag=gcr.io/${PROJECT}/${SERVICE} .

update:
	gcloud run services update ${PROJECT}-${SERVICE} --region=${REGION} --image=gcr.io/${PROJECT}/${SERVICE}

.PHONY: run build update
