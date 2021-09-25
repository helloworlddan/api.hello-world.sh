init:
	gsutil mb -f gs://hwsh-api--terraform-state/ || exit 0
	terraform init

apply:
	terraform apply

destroy:
	terraform destroy

.PHONY: init apply destroy
