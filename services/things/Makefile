SERVICE := things
PROJECT := hwsh-api
REGION := europe-west1
GATEWAY_SA := hwsh-proxy@hwsh-api.iam.gserviceaccount.com

run:
	GOOGLE_CLOUD_PROJECT=${PROJECT} go run *.go

build:
	gcloud builds submit --tag=gcr.io/${PROJECT}/${SERVICE} .

update:
	gcloud run services update ${PROJECT}-${SERVICE} --platform=managed --region=${REGION} --image=gcr.io/${PROJECT}/${SERVICE} --set-env-vars="HWSH_GATEWAY_SA=${GATEWAY_SA},HWSH_ENVIRONMENT=prod,GOOGLE_CLOUD_PROJECT=${PROJECT}"

.PHONY: run build update
