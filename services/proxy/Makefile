SERVICE := proxy
PROJECT := hwsh-api
REGION := europe-west1

init:
	go mod init

run:
	go run server.go

build:
	sh configure.sh
	gcloud builds submit --tag=gcr.io/${PROJECT}/${SERVICE} .
	echo "undefined" > service.json

update:
	gcloud run services update hwsh-${SERVICE} --platform=managed --region=${REGION} --image=gcr.io/${PROJECT}/${SERVICE}

.PHONY: init run build update
