SERVICE := things
PROJECT := hwsh-api
REGION := europe-west1

run:
	bundle install
	GOOGLE_CLOUD_PROJECT=${PROJECT} bundle exec ruby app.rb

build:
	gcloud builds submit --tag=gcr.io/${PROJECT}/${SERVICE} .

update:
	gcloud run services update ${PROJECT}-${SERVICE} --region=${REGION} --image=gcr.io/${PROJECT}/${SERVICE}

.PHONY: run build update
