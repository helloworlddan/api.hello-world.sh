infrastructure:
	make -C infrastructure/ infrastructure

clean:
	rm -rf */.terraform
	rm -rf *.zip
	rm -rf proxy-container/service.json

build-proxy:
	$(eval PROJECT := $(shell sh ../infrastructure/project-id.sh | jq -r '.project'))
	cd proxy-container; sh configure.sh
	gcloud builds submit --tag="gcr.io/${PROJECT}/hwsh-proxy" proxy-container/
	gcloud beta run deploy --platform=managed --region=europe-west1 --image="gcr.io/${PROJECT}/hwsh-proxy" --allow-unauthenticated hwsh-api-proxy

admin-bucket:
	gsutil mb gs://hwsh-api-admin

.PHONY: infrastructure clean admin-bucket build-proxy