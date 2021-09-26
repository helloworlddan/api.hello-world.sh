SERVICE := things
PROJECT := hwsh-api
REGION := europe-west1
GATEWAY_SA := hwsh-api-proxy@hwsh-api.iam.gserviceaccount.com

run:
	GOOGLE_CLOUD_PROJECT=${PROJECT} go run *.go

build:
	gcloud builds submit --tag=gcr.io/${PROJECT}/${SERVICE} .

update:
	gcloud run services update ${PROJECT}-${SERVICE} --region=${REGION} --image=gcr.io/${PROJECT}/${SERVICE}

.PHONY: run build update
