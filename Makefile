
NAME ?= lawn
BACKEND_BUCKET ?=
PROJECT_ID ?=
SERVICE_ACCOUNT ?=
GITHUB_USERNAME ?=
IMAGE_TAG ?= "latest"

.PHONY: terraform/terraform.tf
terraform/terraform.tf:
	PREFIX=${NAME} BUCKET=${BACKEND_BUCKET} \
	envsubst < terraform/terraform.tf.tmpl > terraform/terraform.tf

.PHONY: terraform/variables.tfvars
terraform/variables.tfvars:
	PROJECT_ID=${PROJECT_ID} SERVICE_ACCOUNT=${SERVICE_ACCOUNT} \
	GITHUB_USERNAME=${GITHUB_USERNAME} IMAGE_TAG=${IMAGE_TAG} \
	envsubst < terraform/variables.tfvars.tmpl > terraform/variables.tfvars

.PHONY: init
init: enable_apis
init: terraform/terraform.tf
init:
	cd terraform && terraform init

plan: terraform/variables.tfvars
plan:
	cd terraform && terraform plan -var-file=variables.tfvars

apply: terraform/variables.tfvars
apply:
	cd terraform && terraform apply -var-file=variables.tfvars

.PHONY: login
login:
	gcloud auth application-default login --project ${PROJECT_ID}

.PHONY: enable_apis
enable:
	gcloud services enable containerregistry.googleapis.com --project ${PROJECT_ID}
	gcloud services enable run.googleapis.com --project ${PROJECT_ID}

.PHONY: run
run:
	go run ./cmd/lawn
